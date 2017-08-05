package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MikeRoetgers/deploi/cmd/client/config"
	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/miquella/ask"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the deploi client",
}

var setupConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Create a fresh client config for deploi",
	Run: func(cmd *cobra.Command, args []string) {
		location, err := cmd.Flags().GetString("location")
		if err != nil {
			log.Fatal(err)
		}
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Fatal(err)
		}
		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			log.Fatal(err)
		}
		if !overwrite {
			if _, err := os.Stat(location); !os.IsNotExist(err) {
				log.Fatal("Config file already exists and overwrite was set to false")
			}
		}
		conf := config.NewConfig()
		conf.Host = host
		conf.Location = location

		if err := config.WriteConfig(conf); err != nil {
			log.Fatalf("Failed to write config: %s", err)
		}

		if location != config.GetDefaultConfLocation() {
			fmt.Printf("Please remember to export the environment variable DEPLOI_CONFIG=\"%s\", so deploi will find your config in the custom location.\n", location)
		}
	},
}

var setupUserCmd = &cobra.Command{
	Use:   "user [email]",
	Short: "Register a new user",
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		pw, err := ask.HiddenAsk("Password: ")
		if err != nil {
			log.Fatal(err)
		}
		req := &protobuf.CreateUserRequest{
			Email:    email,
			Password: pw,
		}
		res, err := DeploiClient.CreateUser(context.Background(), req)
		handleGRPCFeedback(err, res.Header)
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
	setupCmd.AddCommand(setupConfigCmd, setupUserCmd)

	setupConfigCmd.Flags().StringP("location", "l", config.GetDefaultConfLocation(), "Location of the config file on your computer")
	setupConfigCmd.Flags().String("host", "localhost:8000", "On which [IP]:[PORT] does deploid listen?")
	setupConfigCmd.Flags().BoolP("overwrite", "o", false, "Overwrite config file if it already exists in given location")

	setupUserCmd.Args = cobra.ExactArgs(1)
}
