package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/codeowners-verifier/pkg/providers"
	"github.com/topfreegames/codeowners-verifier/pkg/verifier"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate provider",
	Short: "Validate the integrity of a CODEOWNERS file",
	Long: fmt.Sprintf(`Check if every entry on the CODEOWNERS file exists on the provider.
Valid providers: %v`, providers.ListProviders()),
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := providers.InitProvider(args[0], cmd.Flag(token).Value.String(), cmd.Flag(baseurl).Value.String())
		if err != nil {
			log.Fatalf("Could not initialize provider: %s", err)
		}
		valid, err := verifier.ValidateCodeownerFile(client, cmd.Flag(codeowners).Value.String())
		if err != nil {
			log.Fatalf("Error reading CODEOWNERS file contents: %s", err)
		}
		if valid {
			log.Info("Valid CODEOWNERS file")
			os.Exit(0)
		} else {
			log.Fatal("Invalid CODEOWNERS file")
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
