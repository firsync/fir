package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
)

func init() {

}

func main() {
	fmt.Println("🌿 hello fir")
	if len(os.Args) < 2 {
		fmt.Println("No command was entered, checking if the current directory is a fir project")
	} else {
		switch os.Args[1] {
		case "save":
			firSaveCase()
			fmt.Println("💾 Snapshot saved")
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

func firSaveCase() {
	unixTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	loadOrCreateConfig()
	hashlist, hashlistErr := getHashListForFolder(".")
	if hashlistErr != nil {
		log.Println("Error assembling file hashes: ", hashlistErr)
	}

	// prompt the user for a save message
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("📝 Enter a save message: [optional] ")
	saveMessage, _ := reader.ReadString('\n')
	saveMessage = strings.TrimRight(saveMessage, "\r\n")

	// prepend the save message to the first line of the file
	filepath := "./.fir/checkpoints/" + unixTimestamp + ".list"
	thisFilePathExists, thisFilePathExistsErr := fileExists(filepath)
	if !thisFilePathExists {
		createFile(filepath)
	}
	if thisFilePathExistsErr != nil {
		log.Println(thisFilePathExistsErr)
	}
	file, err := os.OpenFile(filepath, os.O_RDWR, 0600)
	if err != nil {
		log.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	now := time.Now()
	dateString := now.Format("Mon Jan 2 15:04:05 MST 2006")

	firstLine := "message: " + saveMessage + " " + dateString + "\n"
	temp := firstLine + strings.Join(hashlist, "\n")
	file.Truncate(0)
	file.Seek(0, 0)
	_, _ = file.WriteString(temp)
}

func getHashListForFolder(folderPath string) ([]string, error) {

	theseFiles, theseFilesErr := filePathWalkDir(folderPath)
	if theseFilesErr != nil {
		log.Println("error walking directory: ", theseFilesErr)
	}

	ignoreList, err := readIgnoreList()
	if err != nil {
		return nil, err
	}

	var cleanList []string

	for _, s := range theseFiles {
		if strings.HasPrefix(s, ".git/") || strings.HasPrefix(s, ".fir/") {
			continue
		}
		skip := false
		for _, ss := range ignoreList {
			if strings.HasPrefix(s, ss) {
				skip = true
				break
			}
			if strings.HasPrefix(ss, s) {
				skip = true
				break
			}
			if s == ss {
				skip = true
				break
			}
		}
		if !skip {
			cleanList = append(cleanList, s)
		}
	}

	var hashList []string
	for _, file := range cleanList {
		filePath := filepath.Join(folderPath, file)
		relativePath, err := filepath.Rel(folderPath, filePath)
		if err != nil {
			return nil, err
		}

		// check if file or directory is in ignore list
		if contains(ignoreList, relativePath) {
			log.Println("File ignored: ", relativePath)
			continue
		}

		fileHash, err := calculateFileHash(filePath)
		if err != nil {
			return nil, err
		}
		hashList = append(hashList, fmt.Sprintf("%s %s", fileHash, relativePath))
	}
	return hashList, nil
}

func readIgnoreList() ([]string, error) {
	var ignoreList []string
	if _, err := os.Stat(".ignore"); !os.IsNotExist(err) {
		file, err := os.Open(".ignore")
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ignoreList = append(ignoreList, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}
	var newList []string
	for _, s := range ignoreList {
		s = strings.TrimPrefix(s, "./")
		s = strings.TrimPrefix(s, "/")
		s = strings.TrimSuffix(s, "/")
		s = strings.TrimSuffix(s, "/*")
		if s == "" {
			continue
		}
		newList = append(newList, s)
	}
	thisOtherNewListWhyDoIDoThis := removeDuplicateItems(newList)
	return thisOtherNewListWhyDoIDoThis, nil
}

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha3.New256()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:32]
	hashInString := hex.EncodeToString(hashInBytes)

	return hashInString, nil
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func filePathWalkDir(thisPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(thisPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
