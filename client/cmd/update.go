/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update content of the specified file, or create one if does not exist",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Pass a valid file name.")
		}

		// request := FileUploadRequest("http://localhost:4500/update", "file", "PUT", args)
		request := FileUploadRequest(fmt.Sprintf("http://%v/update", url), "file", "PUT", args)

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
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
