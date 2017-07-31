package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/spf13/cobra"
)

// jobCmd represents the job command
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Manages deployment jobs",
}

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists jobs",
	Run: func(cmd *cobra.Command, args []string) {
		pending, err := cmd.Flags().GetBool("pending")
		if err != nil {
			cmd.Printf("Failed to parse flag: %s", err)
			os.Exit(1)
		}
		showID, err := cmd.Flags().GetBool("showId")
		if err != nil {
			cmd.Printf("Failed to parse flag: %s", err)
			os.Exit(1)
		}
		req := &protobuf.GetJobsRequest{
			Header:  &protobuf.RequestHeader{},
			Pending: pending,
		}
		res, err := DeploiClient.GetJobs(context.Background(), req)
		if err != nil {
			cmd.Printf("Failed to connect to deploid: %s", err)
			os.Exit(1)
		}
		if !res.Header.Success {
			cmd.Printf("Failed to list jobs")
			for _, er := range res.Header.Errors {
				cmd.Printf("Code: %s | Message: %s\n", er.Code, er.Message)
			}
			os.Exit(1)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		if showID {
			fmt.Fprintln(w, "ID\tProject\tBuild\tEnvironment\tNamespace\tStart")
		} else {
			fmt.Fprintln(w, "Project\tBuild\tEnvironment\tNamespace\tStart")
		}

		for _, j := range res.Jobs {
			d := time.Unix(j.CreatedAt, 0).Format(time.RFC3339)
			if showID {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", j.Id, j.Build.ProjectName, j.Build.BuildId, j.Environment.Name, j.Environment.Namespaces[0], d)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", j.Build.ProjectName, j.Build.BuildId, j.Environment.Name, j.Environment.Namespaces[0], d)
			}
		}
		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(jobCmd)
	jobCmd.AddCommand(jobListCmd)

	jobListCmd.Flags().BoolP("pending", "p", true, "Show pending jobs vs. done jobs")
	jobListCmd.Flags().BoolP("showId", "", false, "Show a column with the job id")
}
