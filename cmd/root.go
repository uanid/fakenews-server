package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fakenews-server",
	Short: "",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var ddbTableName string
var sqsUrl string
var awsProfile string
var awsRegion string

func init() {
	rootCmd.PersistentFlags().StringVar(&ddbTableName, "ddb", "fnc1-db", "Dynamodb Table Name")
	rootCmd.PersistentFlags().StringVar(&sqsUrl, "sqs", "https://sqs.ap-northeast-2.amazonaws.com/031804216199/fnc1-queue.fifo", "Full SQS Queue Url")
	rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "fnc", "AWS Shared Credential Profile Name")
	rootCmd.PersistentFlags().StringVar(&awsRegion, "region", "ap-northeast-2", "AWS Region Name")
}
