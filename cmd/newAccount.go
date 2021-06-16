/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunliang711/eth/sdk"
)

// newAccountCmd represents the newAccount command
var newAccountCmd = &cobra.Command{
	Use:   "new-account",
	Short: "create new account",
	Long: `create new account`,
	Run: newAccount,
}

func init() {
	rootCmd.AddCommand(newAccountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newAccountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newAccountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newAccountCmd.Flags().StringP("output", "o", "", "account output file")
}

func writeFile(w io.Writer, address, sk, pk []byte) (err error) {
	_, err = fmt.Fprintf(w, `{"address": "%#x","sk": "%x","pk": "%x"}`, address, sk, pk)
	return
}

func newAccount(cmd *cobra.Command, args []string) {
	sk, pk, address, err := sdk.GenAccount()
	// _, _, _, err := sdk.GenAccount()
	if err != nil {
		fmt.Fprintf(os.Stderr, "create account error: %v", err)
		os.Exit(1)
	}

	outputFile := cmd.Flags().Lookup("output").Value.String()
	if outputFile == "" {
		// output to stdout
		writeFile(os.Stdout, address, sk, pk)
	} else {
		// write to file
		outputFileHandler, err := os.OpenFile(outputFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: write output file error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "info: write account to file: %v\n", outputFile)
		writeFile(outputFileHandler, address, sk, pk)
	}

}
