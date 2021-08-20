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
	ethSdk "github.com/sunliang711/eth/sdk"
	"golang.org/x/crypto/ssh/terminal"
)

// deriveCmd represents the derive command
var deriveCmd = &cobra.Command{
	Use:   "derive",
	Short: "derive address by private key",
	Long: `derive address by private key
`,
	Run: derive,
}

func init() {
	rootCmd.AddCommand(deriveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deriveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deriveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	deriveCmd.Flags().StringP("sk", "k", "", "private key")

}

func derive(cmd *cobra.Command, args []string) {
	sk := cmd.Flags().Lookup("sk").Value.String()
	if sk == "" {
		fmt.Fprintf(os.Stderr, "Enter fromSK: ")
		// read fromSK
		input, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read fromSK error: %v", err)
			os.Exit(1)
		}
		sk = string(input)
		fmt.Fprintln(os.Stderr)
	}
	_, _, address, err := ethSdk.HexToAccount(sk)
	if err != nil {
		fmt.Fprintf(os.Stderr, "derive address error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("address: %v", address)

}
