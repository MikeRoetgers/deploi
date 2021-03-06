package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/spf13/cobra"
)

// environmentCmd represents the environment command
var environmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "Manage environments",
}

var environmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environments",
	Run: func(cmd *cobra.Command, args []string) {
		req := &protobuf.StandardRequest{
			Header: &protobuf.RequestHeader{},
		}
		res, err := DeploiClient.GetEnvironments(context.Background(), req)
		if err != nil {
			cmd.Printf("Failed to register environment: %s", err)
			os.Exit(1)
		}
		if !res.Header.Success {
			cmd.Printf("Failed to register environment\n")
			for _, er := range res.Header.Errors {
				cmd.Printf("Code: %s | Message: %s\n", er.Code, er.Message)
			}
			os.Exit(1)
		}
		cmd.Printf("NAME\tNAMESPACES\n")
		for _, env := range res.Environments {
			cmd.Printf("%s\t%s\n", env.Name, strings.Join(env.Namespaces, ","))
		}
	},
}

var environmentRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register new environment or namespace",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		namespaces, err := cmd.Flags().GetStringArray("namespace")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		req := &protobuf.RegisterEnvironmentRequest{
			Header: &protobuf.RequestHeader{},
			Environment: &protobuf.Environment{
				Name:       name,
				Namespaces: namespaces,
			},
		}
		res, err := DeploiClient.RegisterEnvironment(context.Background(), req)
		if err != nil {
			cmd.Printf("Failed to register environment: %s", err)
			os.Exit(1)
		}
		if !res.Header.Success {
			cmd.Printf("Failed to register environment\n")
			for _, er := range res.Header.Errors {
				cmd.Printf("Code: %s | Message: %s\n", er.Code, er.Message)
			}
			os.Exit(1)
		}
	},
}

var environmentDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment or namespaces in an environment",
	Long: `When only an environment name is supplied, the whole environment is deleted.
When an environment name and additionally one or more namespaces are supplied, the namespaces in the environment are deleted. The environment itself is preserved.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		namespaces, err := cmd.Flags().GetStringArray("namespace")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		req := &protobuf.DeleteEnvironmentRequest{
			Header: &protobuf.RequestHeader{},
			Environment: &protobuf.Environment{
				Name:       name,
				Namespaces: namespaces,
			},
		}
		res, err := DeploiClient.DeleteEnvironment(context.Background(), req)
		if err != nil {
			cmd.Printf("Failed to delete environment: %s", err)
			os.Exit(1)
		}
		if !res.Header.Success {
			cmd.Printf("Failed to delete environment\n")
			for _, er := range res.Header.Errors {
				cmd.Printf("Code: %s | Message: %s\n", er.Code, er.Message)
			}
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(environmentCmd)
	environmentCmd.AddCommand(environmentListCmd, environmentRegisterCmd, environmentDeleteCmd)

	environmentRegisterCmd.Flags().StringP("name", "n", "", "Unique name of the environment")
	environmentRegisterCmd.Flags().StringArrayP("namespace", "", []string{}, "One or more namespaces inside the environment")

	environmentDeleteCmd.Flags().StringP("name", "n", "", "Name of the environment")
	environmentDeleteCmd.Flags().StringArrayP("namespace", "", []string{}, "Namespace(s) in the environment")
}
