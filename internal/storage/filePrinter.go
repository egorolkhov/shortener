package storage

import (
	"encoding/json"
	"os"
)

var Uuid int

type FileJson struct {
	Uuid         int    `json:"uuid"`
	Short_url    string `json:"short_url"`
	Original_url string `json:"original_url"`
}

func FileWrite(k, v string, filepath string) error {
	Uuid += 1
	if filepath == "" {
		return nil
	}
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 644)
	if err != nil {
		return err
	}
	defer file.Close()

	fileJson := FileJson{Uuid, k, v}
	line, err := json.Marshal(fileJson)
	if err != nil {
		return err
	}

	line = append(line, '\n')

	file.Write(line)
	return nil
}
