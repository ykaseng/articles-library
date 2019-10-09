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
	"github.com/spf13/cobra"
	"github.com/ykaseng/articles-library/database/migrate"
)

var reset bool

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "use go-pg migration tool",
	Long:  `migrate uses go-pg migration tool under the hood supporting the same commands and an additional reset command`,
	Run: func(cmd *cobra.Command, args []string) {
		argsMig := args[:0]
		for _, arg := range args {
			switch arg {
			case "migrate", "--reset":
			default:
				argsMig = append(argsMig, arg)
			}
		}

		if reset {
			migrate.Reset()
		}
		migrate.Migrate(argsMig)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	migrateCmd.Flags().BoolVar(&reset, "reset", false, "migrate down to version 0 then up to latest. WARNING: all data will be lost!")
}
