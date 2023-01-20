package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

func init() {
	if len(os.Args) < 2 {
		fmt.Println("No command was entered, checking if the current directory is a fir project")
		// check if the current directory is a fir project and print a help dialogue with a syntax helper if not
	} else {
		switch os.Args[1] {
		case "save":
			fmt.Println("fir save command was run")
			loadOrCreateConfig()
			hashlist, hashlistErr := getHashListForFolder(".")
			if hashlistErr != nil {
				log.Println("Error assembling file hashes: ", hashlistErr)
			}
			writeHashErr := writeLocalHashList(hashlist)
			if writeHashErr != nil {
				log.Println("Error writing hash list: ", writeHashErr)
			}

		case "history":
			fmt.Println("fir history command was run")
		case "sync":
			fmt.Println("fir sync command was run")
		default:
			fmt.Println("Invalid command")
		}
	}
}

func main() {
	log.Println("🌿 hello fir")
	checkIfFilesOrFoldersExist([]string{"~/.fir/config", "./.fir/config", "./."}, true)

}

func checkIfFilesOrFoldersExist(fileAndFolderList []string, createFileFolder bool) bool {
	for _, fileOrFolder := range fileAndFolderList {
		if _, err := os.Stat(fileOrFolder); os.IsNotExist(err) {
			if createFileFolder {
				if err := os.MkdirAll(fileOrFolder, os.ModePerm); err != nil {
					log.Println(err)
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}

// func getHashListForFolder(folderPath string) ([]string, bool) {
// 	hashList := []string{}
// 	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
// 		if info.IsDir() {
// 			return nil
// 		}
// 		file, err := os.Open(path)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		hash := sha3.New256()
// 		if _, err := io.Copy(hash, file); err != nil {
// 			return err
// 		}
// 		hashInBytes := hash.Sum(nil)[:32]
// 		hashInString := hex.EncodeToString(hashInBytes)
// 		relativePath, _ := filepath.Rel(folderPath, path)
// 		hashList = append(hashList, fmt.Sprintf("%s %s", hashInString, relativePath))
// 		return nil
// 	})
// 	if err != nil {
// 		log.Println(err)
// 		return []string{}, false
// 	}
// 	return hashList, true
// }

func getHashListForFolder(folderPath string) ([]string, error) {
	hashList := []string{}
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		hash := sha3.New256()
		if _, err := io.Copy(hash, file); err != nil {
			return err
		}
		hashInBytes := hash.Sum(nil)[:32]
		hashInString := hex.EncodeToString(hashInBytes)
		relativePath, _ := filepath.Rel(folderPath, path)
		hashList = append(hashList, fmt.Sprintf("%s %s", hashInString, relativePath))
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return hashList, nil
}

func loadOrCreateConfig() bool {
	globalConfig := "~/.fir/config"
	localConfig := "./.fir/config"
	config := Config{
		Name:   "Fir User",
		Email:  "user@example.com",
		Remote: "",
		PubKey: "",
	}

	loadGlobalConfig(globalConfig, &config)
	loadLocalConfig(localConfig, &config)

	FirName = config.Name
	FirEmail = config.Email
	FirRemote = config.Remote
	FirPubKey = config.PubKey

	return true
}

func loadGlobalConfig(globalConfig string, config *Config) bool {
	if checkIfFilesOrFoldersExist([]string{globalConfig}, true) {
		globalConfigFile, _ := os.Open(globalConfig)
		defer globalConfigFile.Close()
		byteValue, _ := ioutil.ReadAll(globalConfigFile)
		json.Unmarshal(byteValue, config)
		return true
	} else {
		configJSON, _ := json.Marshal(config)
		ioutil.WriteFile(globalConfig, configJSON, 0644)
		return true
	}
}

func loadLocalConfig(localConfig string, config *Config) bool {
	if checkIfFilesOrFoldersExist([]string{localConfig}, true) {
		localConfigFile, _ := os.Open(localConfig)
		defer localConfigFile.Close()
		byteValue, _ := ioutil.ReadAll(localConfigFile)
		json.Unmarshal(byteValue, config)
		return true
	} else {
		configJSON, _ := json.Marshal(config)
		ioutil.WriteFile(localConfig, configJSON, 0644)
		return true
	}
}

func writeLocalHashList(hashList []string) error {
	unixTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	filepath := "./.fir/checkpoints/" + unixTimestamp + ".list"
	if !checkIfFilesOrFoldersExist([]string{filepath}, true) {
		return errors.New("unable to create the local hash list file")
	}
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range hashList {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}
