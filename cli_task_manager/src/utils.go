package db

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func GetDBPath() (string, error) {
	appDir, err := GetAppDir()
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(appDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, DB_NAME), nil
}

// Expected to work similar to python's click: https://github.com/pallets/click/blob/ca5e1c3d75e95cbc70fa6ed51ef263592e9ac0d0/src/click/utils.py#L449
func GetAppDir() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	var folder string
	switch {
	case strings.HasPrefix(runtime.GOOS, "win"):
		folder = os.Getenv("APPDATA")
		if folder == "" {
			folder = os.Getenv("LOCALAPPDATA")
		}
		if folder == "" {
			folder = homeDir
		}
	case strings.HasPrefix(runtime.GOOS, "darwin"):
		folder = filepath.Join(homeDir, "Library/Application Support")
	default: // Linux is expected
		folder = os.Getenv("XDG_CONFIG_HOME")
		if folder == "" {
			folder = filepath.Join(homeDir, ".config")
		}
	}
	return filepath.Join(folder, APP_NAME), nil
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
