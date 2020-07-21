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
		Short: "For a path, check if a CODEOWNER entry rule apples, excluding member from the ignore flag",
		Long: `For a given path, goes through the CODEONWERS file trying to find a rule that matches the path,
		Also, you can specify a list of members to ignore with the flag -i or --ignore. Example:
		codeowners-verifier verify folder1 --ignore @user1 --ignore @group1`,
		Run: func(cmd *cobra.Command, args []string) {
			co, err := verifier.ReadCodeownersFile(cmd.Flag(codeowners).Value.String())
			if err != nil {
				log.Fatalf("Couldn't read CODEOWNERS file: %s", err)
			}
			rule, valid := verifier.VerifyCodeowner(co, args[0], ignore)
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
