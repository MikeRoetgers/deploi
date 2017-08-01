// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/spf13/cobra"
)

// automationCmd represents the automation command
var automationCmd = &cobra.Command{
	Use:   "automation",
	Short: "Automates deployments to environments",
}

var automationListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists configured deployent automations",
	Run: func(cmd *cobra.Command, args []string) {
		req := &protobuf.GetAutomationsRequest{
			Header: &protobuf.RequestHeader{},
		}
		res, err := DeploiClient.GetAutomations(context.Background(), req)
		handleGRPCFeedback(err, res.Header)
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tProject\tEnvironment\tNamespace\tType\tDetails")
		for _, a := range res.Automations {
			switch auto := a.Automation.(type) {
			case *protobuf.Automation_BranchAutomation:
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", a.Id, auto.BranchAutomation.Project, auto.BranchAutomation.Environment.Name, auto.BranchAutomation.Environment.Namespaces[0], "branch", auto.BranchAutomation.Branch)
			}
		}
		w.Flush()
	},
}

var automationRegisterCmd = &cobra.Command{
	Use:   "register [automation type]",
	Short: "Adds a new automation",
	Long: `Adds a new automation

Currently the following automations are supported:
* branchAutomation
  Whenever a new build in a project originating from a defined branch is registered, it is automatically deployed on the target environment.
  This automation requires the flags -p -b -e -n`,
	Run: func(cmd *cobra.Command, args []string) {
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			fmt.Printf("Failed to parse flag: %s", err)
			os.Exit(1)
		}
		branch, err := cmd.Flags().GetString("branch")
		if err != nil {
			fmt.Printf("Failed to parse flag: %s", err)
			os.Exit(1)
		}
		environment, err := cmd.Flags().GetString("environment")
		if err != nil {
			fmt.Printf("Failed to parse flag: %s", err)
			os.Exit(1)
		}
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			fmt.Printf("Failed to parse flag: %s", err)
			os.Exit(1)
		}
		switch args[0] {
		case "branchAutomation":
			req := &protobuf.RegisterAutomationRequest{
				Header: &protobuf.RequestHeader{},
				Automation: &protobuf.Automation{
					Automation: &protobuf.Automation_BranchAutomation{
						BranchAutomation: &protobuf.BranchAutomation{
							Project: project,
							Branch:  branch,
							Environment: &protobuf.Environment{
								Name:       environment,
								Namespaces: []string{namespace},
							},
						},
					},
				},
			}
			res, err := DeploiClient.RegisterAutomation(context.Background(), req)
			handleGRPCFeedback(err, res.Header)
		default:
			fmt.Printf("The supplied automation type '%s' is not supported.", args[0])
			os.Exit(1)
		}
	},
}

var automationDeleteCmd = &cobra.Command{
	Use:   "delete [automation id]",
	Short: "Deletes an automation",
	Run: func(cmd *cobra.Command, args []string) {
		req := &protobuf.DeleteAutomationRequest{
			Header: &protobuf.RequestHeader{},
			Id:     args[0],
		}
		res, err := DeploiClient.DeleteAutomation(context.Background(), req)
		handleGRPCFeedback(err, res.Header)
	},
}

func init() {
	RootCmd.AddCommand(automationCmd)
	automationCmd.AddCommand(automationListCmd, automationRegisterCmd, automationDeleteCmd)

	automationRegisterCmd.ValidArgs = []string{"branchAutomation"}
	automationRegisterCmd.Args = cobra.OnlyValidArgs
	automationRegisterCmd.Flags().StringP("project", "p", "", "Name of the project")
	automationRegisterCmd.Flags().StringP("branch", "b", "", "Name of the branch")
	automationRegisterCmd.Flags().StringP("environment", "e", "", "Name of the environment")
	automationRegisterCmd.Flags().StringP("namespace", "n", "", "Name of the namespace in the environment")

	automationDeleteCmd.Args = cobra.ExactArgs(1)
}
