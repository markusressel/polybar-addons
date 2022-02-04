package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func WriteTextToFile(text string, path string) (err error) {
	return os.WriteFile(path, []byte(text), 0666)
}

func ReadTextFromFile(path string) (value string, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ReadIntFromFile(path string) (value int, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return -1, err
	}
	text := string(data)
	text = strings.TrimSpace(text)
	if len(text) <= 0 {
		return 0, errors.New(fmt.Sprintf("File is empty: %s", path))
	}
	value, err = strconv.Atoi(text)
	return value, err
}
