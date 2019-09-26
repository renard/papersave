/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
	papersave "github.com/renard/papersave/internal"
)

var (
	Method string = "zxing"
	Debug bool = false
)

// decodeCmd represents the create command
var decodeCmd = &cobra.Command{
	Use:   "decode FILE ...",
	Short: "Decode QR Codes",
	Long: `Decode QR Codes.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runDecode(args)
	},
}

func init() {
	// decodeCmd.Flags().Int64VarP(&MaxSize, "max-size", "",
	// 	MaxSize, "Set maximum file size")
	decodeCmd.Flags().StringVarP(&Method, "method", "m",
		Method, "QRCode decode method.")
	decodeCmd.Flags().BoolVarP(&Debug, "debug", "d",
		Debug, "Debug QRCode decoding.")
	rootCmd.AddCommand(decodeCmd)
}



func runDecode(args []string) {
	var method int
	switch {
	case Method == "zxing": {
		method = papersave.ZXING
	}
	case Method == "zbar": {
		method = papersave.ZBAR
	}
	case Method == "grcode": {
		method = papersave.GRCODE
	}
	default: {
		fmt.Printf(
			"Invalid QR-Code method %s. Should be one of zxing, zbar or grcode.\n",
			Method)
		os.Exit(1)
		}
	}
	var recovered  string
	for _, arg := range(args) {
		data, err := papersave.DecodeQRCode(arg, method)
		if err != nil {
			panic(err)
		}
		recovered += data
	}
	if Debug {
		papersave.CheckData(recovered)
	} else {
		fmt.Printf("%s", recovered)
	}
}

