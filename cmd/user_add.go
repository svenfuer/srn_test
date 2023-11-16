/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		runUserAdd(cmd, args)
	},
}

func init() {
	userCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runUserAdd(cmd *cobra.Command, args []string) {

	if len(args) > 1 {
		fmt.Println("to many arguments")
		os.Exit(2)
	}
	if len(args) < 1 {
		fmt.Println("no arguments")
		os.Exit(2)
	}

	fmt.Printf("Adding User: %s\n", args[0])

	dataStore.LoadDataStoreFromFile(dataPath)

	dataStore.AddUser(args[0], dataPath)

	dataStore.WriteDataStoreToFile(dataPath)

}
