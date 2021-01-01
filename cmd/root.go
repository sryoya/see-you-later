/*
Copyright © 2020 Ryoya Sekino <ryoyasekino1993@gmail.com>

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
	"github.com/sryoya/see-you-later/internal/syl"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "see-you-later",
	Args:  cobra.MinimumNArgs(2),
	Short: "A tool that keeps a website you want to see and opens that later",
	Long: `see-you-later(syl) is a CLI tool to keep a website you want to see and see that later.
	This can be used when you find the site you want to view but you are busy and cannot view soon.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		syl.Run(args[0], args[1], nil)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
