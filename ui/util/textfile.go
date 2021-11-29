package util

import (
	"bytes"
	"io"
	"os"
)

func OpenFileAsText(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
