package storage

import (
	"encoding/json"
	"os"
)

var UUID int

type FileJSON struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func FileWrite(k, v string, filepath string) error {
	UUID += 1
	if filepath == "" {
		return nil
	}
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 644)
	if err != nil {
		return err
	}
	defer file.Close()

	fileJSON := FileJSON{UUID, k, v}
	line, err := json.Marshal(fileJSON)
	if err != nil {
		return err
	}

	line = append(line, '\n')

	file.Write(line)
	return nil
}
