/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"log"
	"os"
	"srn-vault/lib"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Vault",
	Long:  "This Command Initializes the Vault",
	Run: func(cmd *cobra.Command, args []string) {

		// Logic goes here

		log.Println("starting Init")

		runInit()

		log.Println("finished Init")
	},
}

var (
	backdoorEnabled bool
	backdoorKeyPath string
)

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	initCmd.Flags().BoolVarP(&backdoorEnabled, "backdoor", "b", false, "Path to the Config File")
	initCmd.Flags().StringVarP(&backdoorKeyPath, "backdoor-key-path", "", "", "Path to the Backdoor Public Key")

	initCmd.MarkFlagsRequiredTogether("backdoor", "backdoor-key-path")

}

func runInit() {

	if err := os.MkdirAll(dataPath, 0644); err != nil {
		fmt.Println("Failed to Create Data directory")
		panic(err)
	}

	if err := os.MkdirAll(dataPath+"/keys", 0644); err != nil {
		fmt.Println("Failed to Create Keys Directory")
		panic(err)
	}

	if backdoorEnabled {

		// Test if Key can be loaded
		backdoorKey := lib.LoadRSAPublicKey(backdoorKeyPath)

		_bKeyPath := dataPath + "/keys/backdoor_" + uuid.NameSpaceURL.String() + ".pem"

		lib.StoreRSAPublicKey(_bKeyPath, backdoorKey)

		// Set Variables
		dataStore.UseBackdoorKey = true
		dataStore.BackdoorPubKeyPath = _bKeyPath

	}

	dataStore.WriteDataStoreToFile(dataPath)

}
