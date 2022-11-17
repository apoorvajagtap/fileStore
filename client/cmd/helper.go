package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ValidFile(file string) bool {
	// valid file or not
	if _, err := os.Stat(file); err != nil {
		fmt.Printf("!!! %s !!!\n\n", err)
		return false
	}

	// check if filetype is text or not
	out, err := exec.Command("file", file).Output()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if !strings.Contains(string(out), "text") {
		fmt.Printf("The file '%s' is not a text file!\n\n", file)
		return false
	}
	return true
}

func FileUploadRequest(uri, key, method string, args []string) []*http.Request {

	result := []*http.Request{}
	for _, f := range args {

		if !ValidFile(f) {
			continue
		}

		// open file
		file, err := os.Open(f)
		if err != nil {
			fmt.Printf("!!! %s !!!\n\n", err)
			continue
		}
		defer file.Close()

		// create a buffer
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile(key, filepath.Base(fmt.Sprintf("%v", f)))
		if err != nil {
			log.Fatal(err)
			return result
		}

		// copy the file content to buffer
		_, err = io.Copy(part, file)
		err = writer.Close()
		if err != nil {
			log.Fatal(err)
			return result
		}

		req, err := http.NewRequest(method, uri, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		if err != nil {
			log.Fatal(err)
			return result
		} else {
			result = append(result, req)
		}
	}

	return result
}
