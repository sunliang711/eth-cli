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
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sunliang711/eth/sdk"
	"golang.org/x/crypto/ssh/terminal"
)

// exportAccountCmd represents the exportAccount command
var exportAccountCmd = &cobra.Command{
	Use:   "export-account",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: exportAccount,
}

func init() {
	rootCmd.AddCommand(exportAccountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportAccountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportAccountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	exportAccountCmd.Flags().StringP("input", "i", "", "input utc file")
	exportAccountCmd.Flags().StringP("password", "p", "", "password to decrypt utcFile")
	exportAccountCmd.Flags().StringP("output", "o", "", "output file")
}

func exportAccount(cmd *cobra.Command, args []string) {
	var (
		inputFile  string
		outputFile string
		password   string
	)
	inputFile = cmd.Flags().Lookup("input").Value.String()
	outputFile = cmd.Flags().Lookup("output").Value.String()
	password = cmd.Flags().Lookup("password").Value.String()

	fmt.Printf("input file: %v\n", inputFile)
	// sdk.ExportAccount(input,password)
	if password == "" {
		fmt.Fprintf(os.Stderr, "Enter password: ")
		// read password
		pass, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read password error: %v", err)
			os.Exit(1)
		}
		password = string(pass)
		fmt.Fprintln(os.Stderr)
	}
	fmt.Printf("debug: password: %v\n", string(password))

	bs, err := sdk.ExportAccount(inputFile, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: export account error: %v", err)
		os.Exit(1)
	}

	if outputFile != "" {
		outputFileHandler, err := os.OpenFile(outputFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: open file: %v error: %v\n", outputFile, err)
			os.Exit(1)
		}
		_, err = outputFileHandler.Write(bs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: write output file: %v error: %v\n", outputFile, err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("%v\n", string(bs))
	}

}