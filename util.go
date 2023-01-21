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

func removeDuplicateItems(elements []string) []string {
	encountered := map[string]bool{}
	for _, element := range elements {
		encountered[element] = true
	}
	var result []string
	for key := range encountered {
		result = append(result, key)
	}
	return result
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

func createFolder(folder string) error {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func writeFile(filepath string, contents string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(contents)
	if err != nil {
		return err
	}
	return nil
}
