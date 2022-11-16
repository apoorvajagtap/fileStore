/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func fileUploadRequest(uri, key string, args []string) []*http.Request {

	result := []*http.Request{}
	for _, f := range args {
		// valid file or not
		if _, err := os.Stat(f); err != nil {
			log.Fatal(err)
			return result
		}

		// open file
		file, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
			return result
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

		req, err := http.NewRequest("POST", uri, body)
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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file to server",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if empty args, error out
		if len(args) == 0 {
			log.Fatal("Pass a valid file name.")
		}

		// validation check that all args are filepaths or not, and if any of them is not; then error out.
		// create http client.
		// create a multipart request which contains all the files (from args)
		// using http client, execute the request to /upload endpoint.
		request := fileUploadRequest("http://localhost:4500/upload", "file", args)

		client := &http.Client{}

		for _, req := range request {
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			} else {
				body := &bytes.Buffer{}
				_, err := body.ReadFrom(resp.Body)
				if err != nil {
					log.Fatal(err)
				}

				resp.Body.Close()
				fmt.Println(resp.StatusCode)
				// fmt.Println(resp.Header)

				fmt.Println(body)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
