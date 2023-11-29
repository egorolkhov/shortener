package storage

import (
	"bufio"
	"encoding/json"
	"os"
)

func GetStorage(data Storage, filepath string) error {
	if filepath == "" {
		return nil
	}
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		UUID += 1
		fl := FileJSON{}
		line := scanner.Bytes()
		json.Unmarshal(line, &fl)
		data.Add("-1", fl.ShortURL, fl.OriginalURL)
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
