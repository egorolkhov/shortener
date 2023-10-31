package storage

import (
	"bufio"
	"encoding/json"
	"os"
)

func GetStorage(data *Data, filepath string) error {
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
		Uuid += 1
		fl := FileJson{}
		line := scanner.Bytes()
		json.Unmarshal(line, &fl)
		data.Add(fl.Short_url, fl.Original_url)
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
