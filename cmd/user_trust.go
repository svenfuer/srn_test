/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"srn-vault/lib"
)

// trustUserCmd represents the trust command
var trustUserCmd = &cobra.Command{
	Use:   "trust",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if listTrust {
			runShowTrustUserCmd(cmd, args)
		} else {
			runTrustUserCmd(cmd, args)
		}
	},
}

func init() {
	userCmd.AddCommand(trustUserCmd)

	trustUserCmd.Flags().BoolVarP(&trustAll, "all", "a", false, "Trust All new Users")

	trustUserCmd.Flags().BoolVarP(&listTrust, "show", "s", false, "Show Trusted Users")

	trustUserCmd.Flags().StringVarP(&userName, "user", "u", "", "own user name")

	if err := trustUserCmd.MarkFlagRequired("user"); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

}

func runShowTrustUserCmd(cmd *cobra.Command, args []string) {

}

func runTrustUserCmd(cmd *cobra.Command, args []string) {

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

	// Load Key Pair
	keyPair := lib.LoadRSAKeypairDialog(user.PrivateKeyPath)

	if trustAll {
		for i := 0; i < len(dataStore.Users); i++ {

			if dataStore.Users[i].UserName == user.UserName {
				continue // skip self
			}

			iUser := dataStore.Users[i]

			// check if signature is present
			if _, ok := user.TrustedKeysSignaturesMap[iUser.UserName]; ok {
				// Signature already in Map
				continue
			}

			// Generate Signature
			newSignature := dataStore.Users[i].CreateUserPublicKeySignature(keyPair)

			// Store Signature
			user.TrustedKeysSignaturesMap[iUser.UserName] = newSignature

			fmt.Printf("Trusted User: %s with Id: %s\n", dataStore.Users[i].UserName, dataStore.Users[i].Id)
		}
	} else {

		if len(args) < 1 {
			fmt.Println("no username found")
			_ = cmd.Help()
			os.Exit(2)
		}

		if len(args) > 1 {
			fmt.Println("too many usernames")
			_ = cmd.Help()
			os.Exit(2)
		}

		// Get User Id

		trustedUserId := dataStore.GetUserArrId(args[0])

		if trustedUserId == -1 {
			fmt.Println("username not found")
			os.Exit(2)
		}

		// Get User
		trustedUser := dataStore.Users[trustedUserId]

		if trustedUser.UserName == user.UserName {
			fmt.Println("You can not Trust yourself")
			os.Exit(0)
		}

		// check if signature is present
		if _, ok := user.TrustedKeysSignaturesMap[trustedUser.UserName]; !ok {
			// Signature not in  Map

			// Sign Key
			newSignature := trustedUser.CreateUserPublicKeySignature(keyPair)

			// Store Signature
			user.TrustedKeysSignaturesMap[trustedUser.UserName] = newSignature

		} else {
			fmt.Println("Signature is already present")
			os.Exit(2)
		}

		fmt.Printf("Trusted User: %s with Id: %s", trustedUser.UserName, trustedUser.Id)

	}

	// TODO: Private Key im RAM überschreiben

	dataStore.WriteDataStoreToFile(dataPath)

}
