package main

import (
	"encoding/hex"
	"encoding/json"
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

}

func main() {
	fmt.Println("🌿 hello fir")
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
	_, gcExistsErr := fileExists(GlobalConfig)
	if gcExistsErr != nil {
		log.Println("gcExistsErr: ", gcExistsErr)
	}
	_, fExistsErr := fileExists(LocalConfig)
	if fExistsErr != nil {
		log.Println("fExistsErr: ", fExistsErr)
	}
}
func getHashListForFolder(folderPath string) ([]string, error) {
	hashList := []string{}
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return []string{}, err
	}

	for _, file := range files {
		if file.IsDir() {
			if file.Name() == ".fir" || file.Name() == ".git" {
				continue
			}
		}
		filePath := filepath.Join(folderPath, file.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return []string{}, err
		}
		defer file.Close()

		hash := sha3.New256()
		if _, err := io.Copy(hash, file); err != nil {
			return []string{}, err
		}
		hashInBytes := hash.Sum(nil)[:32]
		hashInString := hex.EncodeToString(hashInBytes)
		relativePath, err := filepath.Rel(folderPath, filePath)
		if err != nil {
			return []string{}, err
		}
		hashList = append(hashList, fmt.Sprintf("%s %s", hashInString, relativePath))
	}
	return hashList, nil
}

func loadOrCreateConfig() bool {
	globalConfig := filepath.Join(os.Getenv("HOME"), ".fir/fir.config")
	localConfig := "./.fir/fir.config"
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

func loadGlobalConfig(globalConfig string, config *Config) (bool, error) {

	ffexists, ffexistsErr := folderExists(filepath.Join(os.Getenv("HOME"), ".fir"))
	if ffexistsErr != nil {
		return false, ffexistsErr
	}
	if !ffexists {
		createFolder(filepath.Join(os.Getenv("HOME"), ".fir"))
	}
	fexists, fexistsErr := fileExists(globalConfig)
	if fexistsErr != nil {
		return false, fexistsErr
	}
	if fexists {
		globalConfigFile, err := os.Open(globalConfig)
		if err != nil {
			return false, err
		}
		defer globalConfigFile.Close()
		byteValue, err := ioutil.ReadAll(globalConfigFile)
		if err != nil {
			return false, err
		}
		json.Unmarshal(byteValue, config)
		return true, nil
	} else {
		configJSON, err := json.Marshal(config)
		if err != nil {
			return false, err
		}
		err = createFile(globalConfig)
		if err != nil {
			return false, err
		}
		err = writeFile(globalConfig, string(configJSON))
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

func loadLocalConfig(localConfig string, config *Config) (bool, error) {

	ffexists, ffexistsErr := folderExists("./.fir")
	if ffexistsErr != nil {
		return false, ffexistsErr
	}
	if !ffexists {
		createFolder("./.fir")
	}
	fexists, fexistsErr := fileExists(localConfig)
	if fexistsErr != nil {
		return false, fexistsErr
	}
	if fexists {
		localConfigFile, err := os.Open(localConfig)
		if err != nil {
			return false, err
		}
		defer localConfigFile.Close()
		byteValue, err := ioutil.ReadAll(localConfigFile)
		if err != nil {
			return false, err
		}
		json.Unmarshal(byteValue, config)
		return true, nil
	} else {
		configJSON, err := json.Marshal(config)
		if err != nil {
			return false, err
		}
		err = createFile(localConfig)
		if err != nil {
			return false, err
		}
		err = writeFile(localConfig, string(configJSON))
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

func writeLocalHashList(hashList []string) error {
	unixTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	filepath := "./.fir/checkpoints/" + unixTimestamp + ".list"

	folderExists, folderExistsErr := folderExists("./.fir/checkpoints/")
	if folderExistsErr != nil {
		return folderExistsErr
	}
	if !folderExists {
		createFolderErr := createFolder("./.fir/checkpoints/")
		if createFolderErr != nil {
			return createFolderErr
		}
	}
	fileExists, fileExistsErr := fileExists(filepath)
	if fileExistsErr != nil {
		return fileExistsErr
	}
	if !fileExists {
		createFileErr := createFile(filepath)
		if createFileErr != nil {
			return createFileErr
		}
	}
	for _, line := range hashList {
		writeFileErr := writeFile(filepath, line+"\n")
		if writeFileErr != nil {
			return writeFileErr
		}
	}
	return nil
}
