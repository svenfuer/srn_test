/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"path"
	"srn-vault/lib"
)

// addCmd represents the add command
var addFileCmd = &cobra.Command{
	Use:   "add",
	Short: "encrypt file",
	Long:  "Encrypt the File and add it to the Storage directory",
	Run: func(cmd *cobra.Command, args []string) {
		runAddFileCmd(cmd, args)
	},
}

func runAddFileCmd(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		fmt.Println("no filepath found")
		_ = cmd.Help()
		os.Exit(2)
	}

	if len(args) > 1 {
		fmt.Println("too files")
		_ = cmd.Help()
		os.Exit(2)
	}

	// Load Datastore
	dataStore.LoadDataStoreFromFile(dataPath)

	// TODO: Detect Duplicate files / file Names of user

	// Get Array Id
	arrId := dataStore.GetUserArrId(userName)

	// Check Array Id
	if arrId < 0 {
		fmt.Printf("User %s not found\n", userName)
		os.Exit(2)
	}

	// get User with Array Id
	user := &dataStore.Users[arrId]

	// Load Key Pair
	//keyPair := lib.LoadRSAKeypairDialog(user.PrivateKeyPath)

	// Load Public Key
	pubKey := lib.LoadRSAPublicKey(user.PublicKeyPath)

	// generate AES Key
	aesKey := lib.GenerateAESKeyPointer()

	fileUUID := uuid.NewString()

	encryptedPath := dataPath + "/" + lib.BuildEncryptedFileName(args[0], fileUUID)

	// Encrpt File and AES Key
	lib.EncryptFile(aesKey, args[0], encryptedPath)

	// TODO: Backdoor Encryption

	fileRecord := lib.File{
		Id:       fileUUID,
		Name:     path.Base(args[0]),
		Owner:    userName,
		FilePath: encryptedPath,
		//EncryptedKey:   "",
		//SharedUserKeys: nil,
	}

	encryptedKey := lib.EncryptWithRSA(pubKey, aesKey)
	lib.CleanMemory(aesKey)
	fileRecord.EncryptedKey = base64.StdEncoding.EncodeToString(encryptedKey)

	dataStore.Files = append(dataStore.Files, fileRecord)

	dataStore.WriteDataStoreToFile(dataPath)

}

func init() {
	fileCmd.AddCommand(addFileCmd)

	addFileCmd.Flags().StringVarP(&userName, "user", "u", "", "own user name")

	if err := addFileCmd.MarkFlagRequired("user"); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
