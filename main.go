package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
}

func downloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ss := strings.Split(url, "/")
	if len(ss) == 0 {
		return fmt.Errorf("no file to download")
	}

	fileName := ss[len(ss)-1]
	fmt.Println("file name", fileName)

	fullPath := path.Join(filePath, fileName)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
