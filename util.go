package main

import "os"

func fileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func folderExists(folder string) (bool, error) {
	fileInfo, err := os.Stat(folder)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func createFile(file string) error {
	_, err := os.Create(file)
	if err != nil {
		return err
	}
	return nil
}

func createFolder(folder string) (bool, error) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return false, err
	}
	return true, nil
}

func writeFile(file, contents string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(contents)
	if err != nil {
		return err
	}

	return nil
}
