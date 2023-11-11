package db

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var taskArchiveBucket = []byte("task_archive")
var db *bolt.DB

func InitDB() error {
	dbPath, err := GetDBPath()
	if err != nil {
		return err
	}
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(taskArchiveBucket)
		return err
	})
	return err
}

func StopDB() {
	if db != nil {
		db.Close()
	}
}

func CreateTask(text string) (uint64, error) {
	var uid64 uint64
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		uid64, _ = b.NextSequence()
		return b.Put(itob(uid64), NewTask(uid64, text).ToDB())
	})
	if err != nil {
		return 0, err
	}
	return uid64, nil
}

func ListTasks(filters ...ListFilter) ([]Task, error) {
	return listTasks(taskBucket, filters...)
}

func ListArchivedTasks(filters ...ListFilter) ([]Task, error) {
	return listTasks(taskArchiveBucket, filters...)
}

func listTasks(bucketName []byte, filters ...ListFilter) ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			t := (&Task{}).FromDB(v)
			match := true
			for _, f := range filters {
				match = match && f(t)
			}
			if match {
				tasks = append(tasks, t)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func CompleteTask(key uint64) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		v := b.Get(itob(key))
		if v == nil {
			return fmt.Errorf("task not found (key = %d)", key)
		}
		task := (&Task{}).FromDB(v)
		err := b.Delete(itob(key))
		if err != nil {
			return fmt.Errorf("failed to delete task (key = %d):\n%w", key, err)
		}
		task.Completed = true
		task.CompleteTs = time.Now().Unix()
		b = tx.Bucket(taskArchiveBucket)
		return b.Put(itob(key), task.ToDB())
	})
	return err
}

func RemoveTask(key uint64) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		v := b.Get(itob(key))
		if v == nil {
			return fmt.Errorf("task not found (key = %d)", key)
		}
		task := (&Task{}).FromDB(v)
		err := b.Delete(itob(key))
		if err != nil {
			return fmt.Errorf("failed to delete task (key = %d):\n%w", key, err)
		}
		task.Removed = true
		task.RemoveTs = time.Now().Unix()
		b = tx.Bucket(taskArchiveBucket)
		return b.Put(itob(key), task.ToDB())
	})
	return err
}

type ListFilter func(task Task) bool

var CompletedFilter ListFilter = func(task Task) bool {
	return task.Completed
}

var NotCompletedFilter ListFilter = func(task Task) bool {
	return !task.Completed
}

var RemovedFilter ListFilter = func(task Task) bool {
	return task.Removed
}

var NotRemovedFilter ListFilter = func(task Task) bool {
	return !task.Removed
}
