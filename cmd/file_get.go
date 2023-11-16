/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"srn-vault/lib"

	"github.com/spf13/cobra"
)

// getFileCmd represents the get command
var getFileCmd = &cobra.Command{
	Use:   "get",
	Short: "decrypt a file",
	Long:  "get a file from encrypted Storage",
	Run: func(cmd *cobra.Command, args []string) {
		runGetFileCmd(cmd, args)
	},
}

func init() {
	fileCmd.AddCommand(getFileCmd)

	getFileCmd.Flags().StringVarP(&userName, "user", "u", "", "own user name")

	if err := getFileCmd.MarkFlagRequired("user"); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	getFileCmd.Flags().StringVarP(&decryptDestination, "destination", "d", "", "file destination")

	if err := getFileCmd.MarkFlagRequired("destination"); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

var (
	decryptDestination string
)

func runGetFileCmd(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		fmt.Println("no filepath found")
		_ = cmd.Help()
		os.Exit(2)
	}

	if len(args) > 1 {
		fmt.Println("too many files")
		_ = cmd.Help()
		os.Exit(2)
	}

	// Load Datastore
	dataStore.LoadDataStoreFromFile(dataPath)

	// Get Array Id
	arrId := dataStore.GetUserArrId(userName)

	// Check Array Id
	if arrId < 0 {
		fmt.Printf("User %s not found\n", userName)
		os.Exit(2)
	}

	// get User with Array Id
	user := &dataStore.Users[arrId]

	var fileRecord *lib.File
	fileRecord = nil

	// get File Record
	for i := 0; i < len(dataStore.Files); i++ {
		if dataStore.Files[i].Name == path.Base(args[0]) {
			fileRecord = &dataStore.Files[i]
		}
	}

	if fileRecord == nil {
		fmt.Println("File not found")
		os.Exit(2)
	}

	// bild encrypted Path
	encryptedPath := dataPath + "/" + lib.BuildEncryptedFileName(fileRecord.Name, fileRecord.Id)

	// Load Key Pair
	keyPair := lib.LoadRSAKeypairDialog(user.PrivateKeyPath)

	encryptedKey, err := base64.StdEncoding.DecodeString(fileRecord.EncryptedKey)
	if err != nil {
		fmt.Println("Failed to Decode AES Key")
		os.Exit(2)
	}

	// Decrypt AES Key
	aesKey := lib.DecryptWithRSA(keyPair, encryptedKey)

	// Decrypt file
	lib.DecryptFile(aesKey, encryptedPath, decryptDestination)

	// Clean Key from Memory
	lib.CleanMemory(aesKey)

}
