package main

import (
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

// Checkpoint struct to hold the checkpoint data
type Checkpoint struct {
	Timestamp   int64
	HashList    []string
	Diff        string
	Message     string
	FileSummary string
}

var (
	GlobalConfig = filepath.Join(os.Getenv("HOME"), ".fir/fir.config")
	LocalConfig  = "./.fir/fir.config"
)
