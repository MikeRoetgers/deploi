package cmd

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/MikeRoetgers/deploi"
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
		manifest, err := cmd.Flags().GetString("manifest")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		manifestContent, err := ioutil.ReadFile(manifest)
		if err != nil {
			cmd.Printf("Failed to read manifest file: %s", err)
			os.Exit(1)
		}
		b := &protobuf.Build{
			ProjectName:    project,
			BuildId:        build,
			BuildURL:       url,
			BuildSystemURL: systemUrl,
			BranchName:     branch,
			Files: map[string]string{
				deploi.ManifestFile: string(manifestContent),
			},
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

var buildDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a build onto an environment",
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
		environment, err := cmd.Flags().GetString("environment")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		manifest, err := cmd.Flags().GetString("manifest")
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		req := &protobuf.DeployRequest{
			Project:     project,
			BuildId:     build,
			Environment: environment,
			Namespace:   namespace,
		}
		if manifest != "" {
			manifestContent, err := ioutil.ReadFile(manifest)
			if err != nil {
				cmd.Printf("Failed to read manifest file: %s", err)
				os.Exit(1)
			}
			req.Files = map[string]string{
				deploi.ManifestFile: string(manifestContent),
			}
		}

		res, err := DeploiClient.DeployBuild(context.Background(), req)
		if err != nil {
			cmd.Printf("Failed to connect to deploid: %s\n", err)
			os.Exit(1)
		}
		if !res.Header.Success {
			cmd.Printf("Failed to start a deployment")
			for _, er := range res.Header.Errors {
				cmd.Printf("Code: %s | Message: %s\n", er.Code, er.Message)
			}
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.AddCommand(registerCmd, buildListCmd, buildDeployCmd)

	registerCmd.Flags().StringP("project", "p", "", "Name of the project")
	registerCmd.Flags().StringP("build", "b", "", "Version/ID of the build")
	registerCmd.Flags().StringP("url", "u", "", "URL that leads to the build, e.g. https://registry.deploi.io:5000/myproject:161")
	registerCmd.Flags().StringP("systemUrl", "s", "", "Deeplink back to the build in the build system, e.g. Jenkins")
	registerCmd.Flags().StringP("branch", "", "", "Name of the branch the build originated from")
	registerCmd.Flags().StringP("manifest", "m", "", "Path to the manifest file needed for deployment")

	buildListCmd.Flags().StringP("project", "p", "", "Name of the project")

	buildDeployCmd.Flags().StringP("project", "p", "", "Name of the project")
	buildDeployCmd.Flags().StringP("build", "b", "", "Version/ID of the build")
	buildDeployCmd.Flags().StringP("environment", "e", "", "Name of the environment")
	buildDeployCmd.Flags().StringP("namespace", "n", "", "Namespace in the environment")
	buildDeployCmd.Flags().StringP("manifest", "m", "", "Overwrites manifest file included in the build")
}
