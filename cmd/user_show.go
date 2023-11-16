/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var showUsersCmd = &cobra.Command{
	Use:   "show",
	Short: "show User info",
	Long:  "Show the Information of a User",
	Run: func(cmd *cobra.Command, args []string) {
		runShowUsersCmd(cmd, args)
	},
}

var (
	showAll bool
)

func init() {
	userCmd.AddCommand(showUsersCmd)
	showUsersCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show All Users")
}

func runShowUsersCmd(cmd *cobra.Command, args []string) {

	dataStore.LoadDataStoreFromFile(dataPath)

	if showAll {
		for i := 0; i < len(dataStore.Users); i++ {
			fmt.Println("########################################")
			dataStore.Users[i].PrintUser()
		}
		fmt.Println("########################################")
		return
	}

	if len(args) < 1 {
		fmt.Println("no username found")
		_ = cmd.Help()
		os.Exit(2)
	}

	for i := 0; i < len(args); i++ {
		id := dataStore.GetUserArrId(args[i])
		if id == -1 {
			fmt.Println("username not found")
			os.Exit(2)
		}
		fmt.Println("########################################")
		dataStore.Users[id].PrintUser()
	}
	fmt.Println("########################################")

}
