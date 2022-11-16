package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const parentDir = "./uploads"

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	// Handle multiple files
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {

		if _, err := os.Stat(fmt.Sprintf("%s/%s", parentDir, fileHeader.Filename)); err == nil {
			http.Error(w, fmt.Sprintf(">>> The file '%s/%s' already exists!\n", parentDir, fileHeader.Filename), http.StatusInternalServerError)
			log.Printf("The file '%s/%s' already exists!\n", parentDir, fileHeader.Filename)
			continue
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err = os.MkdirAll(parentDir, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.Create(fmt.Sprintf("%s/%s", parentDir, fileHeader.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, fmt.Sprintf(">>> File %s uploaded successfully", fileHeader.Filename))
	}

}

func getHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	// if we want the list as sorted, shall opt for ioutil.ReadDir
	dir, err := os.Open(parentDir)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()
	files, err := dir.ReadDir(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		w.Write([]byte(fmt.Sprintf("%s\n", f.Name())))
		fmt.Println(f.Name())
	}
}

// func removeHandler(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != "DELETE" {
// 		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
// 		return
// 	}

// }

// func modifyHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "PUT" {
// 		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	file, fileHeader, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	defer file.Close()

// 	if f, err := os.Stat(fmt.Sprintf("%s/%s", parentDir, fileHeader.Filename)); err != nil {
// 		fmt.Println(f)
// 	}
// }

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/get", getHandler)

	// TODO: Configurting a DB to keep track of the files, so that we can perform delete and update efficiently.
	// mux.HandleFunc("/delete", removeHandler)
	// mux.HandleFunc("/update", modifyHandler)

	if err := http.ListenAndServe(":4500", mux); err != nil {
		log.Fatal(err)
	}
}
