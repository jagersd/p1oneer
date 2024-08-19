package pparser

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

const (
	ErrNoReqFiles   = "Err: No p1oneer configuration files found "
	ErrReqNoRead    = "Err: The config file can't be opened "
	ErrPrioConflict = "Err: Priority conflict, exiting "
)

type StartRequest struct {
	Title    string
	Priority uint8    `json:"priority"`
	ReqType  string   `json:"type"`
	Command  string   `json:"command"`
	Args     []string `json:"arguments"`
}

func ParseConfigFiles() map[uint8]StartRequest {
	configFiles := collectRequestFiles()

	var startRequests = make(map[uint8]StartRequest)
	for _, c := range configFiles {
		startReq := parseReqfile(c)
		if _, ok := startRequests[startReq.Priority]; ok {
			log.Fatal(ErrPrioConflict)
		} else {
			startRequests[startReq.Priority] = startReq
		}
	}

	return startRequests
}

func getConfigDir() string {
	d := os.Getenv("P1ONEER_CONFIG_DIR")
	if d == "" {
		log.Fatal("Please set P1ONEER_CONFIG_DIR environment variable to indicate the location of your configuration files")
	}
	return d
}

func collectRequestFiles() []string {
	dir := getConfigDir()
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(ErrNoReqFiles, err)
	}
	var reqFiles []string
	for _, f := range files {
		fileName := f.Name()
		if len(fileName) < 5 {
			continue
		}
		if fileName[len(fileName)-5:] == ".json" {
			reqFiles = append(reqFiles, dir+"/"+fileName)
		}
	}
	if len(reqFiles) == 0 {
		log.Fatal(ErrNoReqFiles)
	}
	return reqFiles
}

func parseReqfile(reqFileName string) StartRequest {
	b, err := os.ReadFile(reqFileName)
	if err != nil {
		log.Fatal(ErrReqNoRead, err)
	}
	var request StartRequest
	request.Title = strings.Replace(reqFileName, ".json", "", 1)
	if err := json.Unmarshal(b, &request); err != nil {
		log.Fatal(ErrReqNoRead, err)
	}
	return request
}
