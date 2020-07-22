package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	token      = "token"
	baseurl    = "base-url"
	codeowners = "codeowners"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codeowners-verifier [flags] action path",
	Short: "Verify the existence of a CODEONWERS to a file based on validation rules",
	Long: `Codeowners-verifier is a tool made for running on CI pipelines.
	It verifies the integrity of your CODEOWNERS file based on a predefined provider (currently only GITLAB),
	You can check if every user and group declared actually exists. You can also check if a file has an CODEOWNER
	defined, using the --ignore flag to ignore OWNERS.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Logging setup
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	cobra.OnInitialize(initConfig)
	if err := viper.BindEnv(token, "CODEOWNER_PROVIDER_TOKEN"); err != nil {
		log.Fatal("error initializing viper for env CODEOWNER_PROVIDER_TOKEN")
	}
	rootCmd.PersistentFlags().String(token, viper.GetString(token), "Token to be used to authenticate with the provider.(Defaults to CODEOWNER_PROVIDER_TOKEN env var)")
	if err := viper.BindPFlag(token, rootCmd.PersistentFlags().Lookup(token)); err != nil {
		log.Fatal("error binding viper for flag CODEOWNER_PROVIDER_TOKEN")
	}
	if err := viper.BindEnv(baseurl, "CODEOWNER_PROVIDER_URL"); err != nil {
		log.Fatal("error initializing viper for env CODEOWNER_PROVIDER_URL")
	}
	rootCmd.PersistentFlags().String(baseurl, viper.GetString(baseurl), "BaseURL to connect to the provider (Defaults to CODEOWNER_PROVIDER_URL env var)")
	if err := viper.BindPFlag(baseurl, rootCmd.PersistentFlags().Lookup(baseurl)); err != nil {
		log.Fatal("error binding viper for flag CODEOWNER_PROVIDER_URL")
	}
	if err := viper.BindEnv(codeowners, "CODEOWNER_PATH"); err != nil {
		log.Fatal("error initializing viper for env CODEOWNER_PATH")
	}
	rootCmd.PersistentFlags().String(codeowners, viper.GetString(codeowners), "Path to the CODEOWNERS file (Defaults to CODEOWNER_PATH env var)")
	if err := viper.BindPFlag(codeowners, rootCmd.PersistentFlags().Lookup(codeowners)); err != nil {
		log.Fatal("error binding viper for flag CODEOWNER_PATH")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
