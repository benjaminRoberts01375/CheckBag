package main

import (
	"os"
	"path/filepath"
)

var fileSystem = FileSystem{
	BasePath: "/CheckBag",
	UserData: "userdata.txt",
}

type FileSystem struct {
	BasePath string `json:"base_path"`
	UserData string `json:"user_data"`
}

func (fs *FileSystem) GetUserDataPath() string {
	return filepath.Join(fs.BasePath, fs.UserData)
}

func (fs *FileSystem) GetUserData() (string, error) {
	data, err := os.ReadFile(fs.GetUserDataPath())
	return string(data), err
}

func (fs *FileSystem) SetUserData(data string) error {
	newFile, err := os.Create(fs.GetUserDataPath())
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = newFile.WriteString(data)
	return err
}
