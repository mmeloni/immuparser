/*
Copyright 2019-2020 vChain, Inc.

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
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var Username string
var Source string
var Ledger string

var FlushSize int

var rootCmd = &cobra.Command{
	Use:   "immuparser",
	Short: "Parse a file and insert data into immudb",
	Long: `Parse a file with following format:
2020-05-21 20:02:00 73.115.58.89:50800 Wine`,

	Example: `
main --ledger myledger
main --ledger myledger --flushSize 1000
main --ledger myledger --flushSize 1000 --source /tmp/log.txt
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Parse()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.immuparser.yaml)")
	rootCmd.Flags().StringVarP(&Source, "source", "s", "/tmp/log.txt", "source file path")
	//rootCmd.Flags().StringVarP(&Username, "username", "u", "default", "default user name")
	rootCmd.Flags().StringVarP(&Ledger, "ledger", "l", "default", "ledger name")

	rootCmd.Flags().IntVarP(&FlushSize, "flushSize", "f", 1000, "flush size")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".immuparser")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
