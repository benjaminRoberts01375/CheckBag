package main

import (
	"errors"
	"os"
	"path/filepath"
)

func init() {
	fileSystem.Setup()
}

var fileSystem = FileSystem{
	BasePath: "/CheckBag",
	UserData: "userdata.txt",
}

type FileSystem struct {
	BasePath string `json:"base_path"`
	UserData string `json:"user_data"`
}

func (fs *FileSystem) Setup() {
	if err := os.MkdirAll(fs.BasePath, os.ModePerm); err != nil {
		panic(err)
	}
}

func (fs *FileSystem) Write(path string, data string) error {
	newFile, err := os.Create(path)
	if err != nil {
		return errors.New("Could not create file: " + err.Error())
	}
	defer newFile.Close()
	_, err = newFile.WriteString(data)
	if err != nil {
		return errors.New("Could not write to file: " + err.Error())
	}
	return nil
}

func (fs *FileSystem) GetUserData() (string, error) {
	data, err := os.ReadFile(fs.GetUserDataPath())
	data, err := os.ReadFile(filepath.Join(fs.BasePath, fs.UserData))
	return string(data), err
}

func (fs *FileSystem) SetUserData(data string) error {
	return fs.Write(filepath.Join(fs.BasePath, fs.UserData), data)
}
	if err != nil {
		return errors.New("Could not write to file: " + err.Error())
	}
	return nil
}
