package pparser

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

const (
	ErrNoReqFiles = "Err: No p1oneer configuration files found "
	ErrReqNoRead  = "Err: The config file can't be opened "
)

type StartRequest struct {
	Title   string
	ReqType string `json:"type"`
	Command string `json:"command"`
}

func ParseConfigFiles() []StartRequest {
	configFiles := collectRequestFiles()

	var sr []StartRequest
	for _, c := range configFiles {
		startReq := parseReqfile(c)
		sr = append(sr, startReq)
	}

	return sr
}

func collectRequestFiles() []string {
	dir := "./examples"
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
