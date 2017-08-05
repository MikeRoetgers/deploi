package cmd

import (
	"context"
	"os"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manages software projects",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all projects",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := DeploiClient.GetProjects(context.Background(), &protobuf.StandardRequest{
			Header: getRequestHeader(),
		})
		if err != nil {
			cmd.Printf("Failed to connect to deploid: %s", err)
			os.Exit(1)
		}
		if !res.Header.Success {
			cmd.Println("Failed to list projects")
			for _, er := range res.Header.Errors {
				cmd.Printf("Code: %s | Message: %s\n", er.Code, er.Message)
			}
			os.Exit(1)
		}
		cmd.Printf("PROJECTS\n")
		for _, p := range res.Projects {
			cmd.Printf("%s\n", p)
		}
	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
}
