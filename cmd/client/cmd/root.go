package cmd

import (
	"fmt"
	"os"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/spf13/cobra"
)

var cfgFile string
var DeploiClient protobuf.DeploiServerClient

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "deploi",
	Short: "deploi interacts with deploid to manage deployments",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	/*
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			viper.AddConfigPath(home)
			viper.SetConfigName(".deploi")
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	*/
}

func handleGRPCFeedback(err error, header *protobuf.ResponseHeader) {
	if err != nil {
		fmt.Printf("RPC request failed: %s", err)
		os.Exit(1)
	}
	if !header.Success {
		fmt.Println("Request was unsuccessful")
		for _, er := range header.Errors {
			fmt.Printf("Code: %s | Message: %s\n", er.Code, er.Message)
		}
		os.Exit(1)
	}
}
