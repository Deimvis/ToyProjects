package db

import (
	"encoding/json"
	"time"
)

type Task struct {
	Key        uint64 `json:"key"`
	Text       string `json:"text"`
	Completed  bool   `json:"completed"`
	Removed    bool   `json:"removed"`
	CreateTs   int64  `json:"create_ts"`
	CompleteTs int64  `json:"complete_ts"`
	RemoveTs   int64  `json:"remove_ts"`
}

func NewTask(key uint64, text string) *Task {
	return &Task{key, text, false, false, time.Now().Unix(), 0, 0}
}

func (t *Task) FromDB(value []byte) (ret Task) {
	err := json.Unmarshal(value, &ret)
	if err != nil {
		panic(err)
	}
	return
}

func (t *Task) ToDB() []byte {
	tJson, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return tJson
}
