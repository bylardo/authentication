package helpers

import (
	"io/ioutil"
)

func ReadDatabase(path string) (bool, []byte) {

	jsonData, err := ioutil.ReadFile(path)

	if err != nil {
		return false, nil
	}
	return true, jsonData
}
