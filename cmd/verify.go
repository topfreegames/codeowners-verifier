package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/codeowners-verifier/pkg/verifier"
)

// checkCmd represents the check command
var (
	verifyCmd = &cobra.Command{
		Use:   "verify path",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			co, err := verifier.ReadCodeownersFile(cmd.Flag(codeowners).Value.String())
			if err != nil {
				log.Fatalf("Couldn't read CODEOWNERS file: %s", err)
			}
			rule, valid := verifier.CheckCodeowner(co, args[0], ignore)
			if valid {
				log.Infof("Found matching rule on line %d: %s %s", rule.Line, rule.Path, rule.Owners)
				os.Exit(0)
			} else {
				log.Fatalf("Missing CODEOWNER entry, matched rule from line %d don't have valid owners: %s %s. Check your ignore rules.", rule.Line, rule.Path, rule.Owners)
			}
		},
	}
	ignore []string
)

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringSliceVarP(&ignore, "ignore", "i", []string{}, "Comma separated list of entries to ignore when validating a path E.g: @user1,@group1,@user2")
}
