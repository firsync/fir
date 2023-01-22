package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
)

func init() {

}

func main() {
	fmt.Println("🌿 Hello Fir")
	if len(os.Args) < 2 {
		firNoArgsCase()
	} else {
		switch os.Args[1] {
		case "save":
			firSaveCase()
			fmt.Println("💾 Snapshot saved")
		case "history":
			firHistoryCase()
		case "sync":
			firSyncCase()
			fmt.Println("fir sync command was run")
		default:
			fmt.Println("💥 I don't know that command yet.\n↪️  Try just typing `fir`")
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
func firSyncCase() {
	keys := initKeys()
	jsonData := map[string]string{
		"public_key": keys.publicKey,
		"signature":  keys.signedKey,
	}
	jsonValue, _ := json.Marshal(jsonData)

	req, err := http.NewRequest("POST", "https://firsync.com/register", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(body))
	}
	fmt.Println("Public key registered with the server")

	// cmd := exec.Command("rsync", "-avz", "--delete", "./", "firsync.com:"+keys.publicKey)
	cmd := exec.Command("rsync", "-avz", "--delete", "-e", "ssh -i "+privKeyFilePath+" -l "+keys.publicKey, "./", "firsync.com:"+keys.publicKey)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Code synced with the server")

}

func firNoArgsCase() {
	WeReGood, weReNotGood := folderExists(".fir")
	if weReNotGood != nil {
		fmt.Println("🤔 This folder isn't a fir project yet...\n...but if you meant to create one, you can type `fir save` at any time to get started :)")
	}
	if WeReGood {
		fmt.Println("👍 This folder is a fir project :)\n- Type `fir save` to save a snapshot of your progress\n- Type `fir history` to view your saves so far")
	}
}

func firHistoryCase() {

	now := time.Now()
	dateString := now.Format("Mon Jan 2 15:04:05 MST 2006")

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("🗒️ ", "Project history for \""+mydir+"\"")

	fmt.Println(dateString + " 👀 you are here")
	checkpointDir := "./.fir/checkpoints"
	files, err := os.ReadDir(checkpointDir)
	if err != nil {
		log.Fatal(err)
	}

	// sort files by unix timestamp in filename
	sort.Slice(files, func(i, j int) bool {
		iName := files[i].Name()
		jName := files[j].Name()
		iTimestamp := iName[:strings.Index(iName, ".")]
		jTimestamp := jName[:strings.Index(jName, ".")]
		return iTimestamp > jTimestamp
	})

	// print first line of each file
	for _, file := range files {
		filePath := filepath.Join(checkpointDir, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Scan()
		fmt.Println(scanner.Text())
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

	firstLine := dateString + " " + saveMessage + "\n"
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
