package cmd

import (
	"log"

	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/topfreegames/codeowners-verifier/pkg/providers"
	"github.com/topfreegames/codeowners-verifier/pkg/verifier"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate provider codeowner_file",
	Short: "Validate the integrity of a CODEOWNERS file",
	Long:  `Check if every entry on the CODEOWNERS file exists on the provider.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := providers.InitProvider(args[0], cmd.Flag(token).Value.String(), cmd.Flag(baseurl).Value.String())
		if err != nil {
			log.Fatalf("Could not initialize provider: %s", err)
		}
		valid, err := verifier.CheckCodeowner(client, args[1])
		if err != nil {
			log.Fatalf("Error reading CODEOWNERS file contents: %s", err)
		}
		if valid {
			fmt.Printf("Valid CODEOWNERS file")
			os.Exit(0)
		} else {
			fmt.Printf("Invalid CODEOWNERS file")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
