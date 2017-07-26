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
	"os"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Manages builds",
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers a new build in deploid",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		build, err := cmd.Flags().GetString("build")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		systemUrl, err := cmd.Flags().GetString("systemUrl")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		branch, err := cmd.Flags().GetString("branch")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		b := &protobuf.Build{
			ProjectName:    project,
			BuildId:        build,
			BuildURL:       url,
			BuildSystemURL: systemUrl,
			BranchName:     branch,
		}
		req := &protobuf.NewBuildRequest{
			Build: b,
		}
		res, err := DeploiClient.RegisterNewBuild(context.Background(), req)
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		if !res.Header.Success {
			cmd.Println("Failed to register new build")
			for _, er := range res.Header.Errors {
				cmd.Printf("Code: %s | Message: %s\n", er.Code, er.Message)
			}
			os.Exit(1)
		}
	},
}

var buildListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists builds",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := cmd.LocalFlags().GetString("project")
		if err != nil {
			cmd.Printf("Failed to parse project flag: %s\n", err)
			os.Exit(1)
		}
		if project == "" {
			cmd.Printf("You have to provide a project name\n")
			os.Exit(1)
		}
		res, err := DeploiClient.GetBuilds(context.Background(), &protobuf.GetBuildsRequest{
			ProjectName: project,
		})
		if err != nil {
			cmd.Printf("Failed to connect to deploid: %s\n", err)
			os.Exit(1)
		}
		if !res.Header.Success {
			cmd.Printf("Failed to list projects")
			for _, er := range res.Header.Errors {
				cmd.Printf("Code: %s | Message: %s\n", er.Code, er.Message)
			}
			os.Exit(1)
		}
		cmd.Printf("Build\tBranch\tURL\n")
		for _, b := range res.Builds {
			cmd.Printf("%s\t%s\t%s\n", b.BuildId, b.BranchName, b.BuildURL)
		}
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.AddCommand(registerCmd)
	buildCmd.AddCommand(buildListCmd)

	registerCmd.Flags().StringP("project", "p", "", "Name of the project")
	registerCmd.Flags().StringP("build", "b", "", "Version/ID of the build")
	registerCmd.Flags().StringP("url", "u", "", "URL that leads to the build, e.g. https://registry.deploi.io:5000/myproject:161")
	registerCmd.Flags().StringP("systemUrl", "s", "", "Deeplink back to the build in the build system, e.g. Jenkins")
	registerCmd.Flags().StringP("branch", "", "", "Name of the branch the build originated from")

	buildListCmd.Flags().StringP("project", "p", "", "Name of the project")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}