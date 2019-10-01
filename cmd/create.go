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
	MaxSize int64 = 8*1024
	Format string = "pdf"
)

var Opts papersave.CreateOptions

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create FILE",
	Short: "Create a file backup",
	Long: `Create a file paper backup.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCreate(args[0])
	},
}

func init() {
	createCmd.Flags().Int64VarP(&MaxSize, "max-size", "",
		MaxSize, "Set maximum file size.")
	createCmd.Flags().StringVarP(&Format, "format", "f",
		Format, "Set maximum file size.")
	createCmd.Flags().BoolVarP(&Opts.Encrypt, "encrypt", "e",
		Opts.Encrypt, "Encrypt file.")
	createCmd.Flags().StringVarP(&Opts.Password, "password", "p",
		Opts.Password, "Set encryption password.")
	createCmd.Flags().BoolVarP(&Opts.Keep, "keep", "k",
		Opts.Keep, "Keep build files.")
	createCmd.Flags().BoolVarP(&Opts.ShowPassword, "show-password", "",
		Opts.ShowPassword, "Show password in generated file.")
	rootCmd.AddCommand(createCmd)
}


func checkCreateArgs(filename string) {
	info, err := os.Stat(filename)
	switch {
	case os.IsNotExist(err): {
		fmt.Println(err)
		os.Exit(1)
	}
	case info.IsDir(): {
		fmt.Fprintf(os.Stderr, "%s is a directory\n", filename)
		os.Exit(1)
	}
	case info.Size() > MaxSize: {
		fmt.Fprintf(os.Stderr, "%s is too large (%d bytes)\n",
			filename, info.Size())
		os.Exit(1)
	}
	}
}


func runCreate(filename string) {
	checkCreateArgs(filename)
	p := papersave.New(filename, Opts)
	switch {
		case Format == "pdf": {
		//p.GenQRCode()
			p.WritePDF()
		}
		case Format == "txt": {
		p.WriteText()
		}
		default: {
		fmt.Fprintf(os.Stderr, "Unknown format %s\n", Format)
			os.Exit(1)
		}
	}
}

