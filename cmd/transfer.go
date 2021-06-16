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
	"math/big"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sunliang711/eth/sdk"
	"golang.org/x/crypto/ssh/terminal"
)

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer eth cmd",
	Long: `transfer eth from sender to receiver:

available flags(no flags will enter interactive mode):
    --rpc rpc url
    --from sender sk
    --to receiver address
    --value value of eth to transfer(uint: eth/1e18)
    --price optional gas price
    --nonce optional nonce`,
	Run: transfer,
}

func init() {
	rootCmd.AddCommand(transferCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transferCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transferCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// from sk
	// to addr
	// value
	transferCmd.Flags().StringP("rpc", "r", "", "rpc URL")
	transferCmd.Flags().StringP("from", "f", "", "from sk")
	transferCmd.Flags().StringP("to", "t", "", "to address")
	transferCmd.Flags().StringP("value", "v", "0", "transfer value(unit: gas)")
	transferCmd.Flags().Uint64P("price", "p", 0, "gas price [optional]")
	transferCmd.Flags().Uint64P("nonce", "n", 0, "nonce [optional]")

}

const (
	transferEthGasLimit = 21000
)

func transfer(cmd *cobra.Command, args []string) {
	rpcURL := cmd.Flags().Lookup("rpc").Value.String()
	fromSK := cmd.Flags().Lookup("from").Value.String()
	toAddr := cmd.Flags().Lookup("to").Value.String()
	value := cmd.Flags().Lookup("value").Value.String()
	priceStr := cmd.Flags().Lookup("price").Value.String()
	nonceStr := cmd.Flags().Lookup("nonce").Value.String()

	if fromSK == "" {
		fmt.Fprintf(os.Stderr, "Enter fromSK: ")
		// read fromSK
		sk, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read fromSK error: %v", err)
			os.Exit(1)
		}
		fromSK = string(sk)
		fmt.Fprintln(os.Stderr)
	}

	v := big.NewInt(0)
	_, ok := v.SetString(value, 10)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: invalid value number\n")
		os.Exit(1)
	}

	txMan, err := sdk.New(rpcURL, 0, transferEthGasLimit, 0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer txMan.Close()

	var (
		gasPrice uint64
		nonce    uint64
	)
	gasPrice, err = strconv.ParseUint(priceStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gas price parse error: %v", err)
		os.Exit(1)
	}
	nonce, err = strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nonce parse error: %v", err)
		os.Exit(1)
	}

	_, _, fromAddr, err := sdk.HexToAccount(fromSK)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: parse sender address from sk error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("info: transfer eth\nfrom: %#x\nto: %#v\nvalue: %v\ngasPrice: %v\nnonce: %v\n", fromAddr, toAddr, v, gasPrice, nonce)
	hash, err := txMan.TransferEth(fromSK, toAddr, v, gasPrice, nonce)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: transfer eth error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("transaction hash: %v\n", hash)

}
