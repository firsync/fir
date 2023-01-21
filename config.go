package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var (
	FirName   string
	FirEmail  string
	FirRemote string
	FirPubKey string
)

type Config struct {
	Name   string
	Email  string
	Remote string
	PubKey string
}

var (
	GlobalConfig = filepath.Join(os.Getenv("HOME"), ".fir/fir.config")
	LocalConfig  = "./.fir/fir.config"
)

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
	ff := filepath.Join(os.Getenv("HOME"), ".fir")
	if _, err := os.Stat(ff); os.IsNotExist(err) {
		if err := os.MkdirAll(ff, os.ModePerm); err != nil {
			return false, err
		}
	}

	if _, err := os.Stat(globalConfig); os.IsNotExist(err) {
		f, err := os.Create(globalConfig)
		if err != nil {
			return false, err
		}
		defer f.Close()
		configJSON, err := json.Marshal(config)
		if err != nil {
			return false, err
		}
		if _, err := f.Write(configJSON); err != nil {
			return false, err
		}
		return true, nil
	} else {
		f, err := os.Open(globalConfig)
		if err != nil {
			return false, err
		}
		defer f.Close()
		if err := json.NewDecoder(f).Decode(config); err != nil {
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
		fileInfo, _ := localConfigFile.Stat()
		fileSize := fileInfo.Size()
		byteValue := make([]byte, fileSize)
		localConfigFile.Read(byteValue)
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
