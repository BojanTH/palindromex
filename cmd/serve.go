/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"palindromex/web"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long: "Used to server files from the 'web' dir",
	Run: func(cmd *cobra.Command, args []string) {
		web.Make()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serveCmd.Flags().StringVar(&web.AppPort, "app_port", ":8080", "Port on which to run the app.")
	serveCmd.Flags().StringVar(&web.SessionSecret, "app_session_secret", "fhvnvehpe8wunfe899he9rhifnhwliea", "A secret used for the session store")
	serveCmd.Flags().StringVar(&web.DbHost, "db_host", "localhost", "Database host.")
	serveCmd.Flags().StringVar(&web.DbName, "db_name", "postgres", "Database name.")
	serveCmd.Flags().StringVar(&web.DbUser, "db_user", "root", "Database user.")
	serveCmd.Flags().StringVar(&web.DbPassword, "db_pwd", "root", "Database password.")
	serveCmd.Flags().StringVar(&web.DbPort, "db_port", "5432", "Database port.")
	serveCmd.Flags().StringVar(&web.DbSslMode, "db_ssl_mode", "disable", "Database ssl mode.")
}
