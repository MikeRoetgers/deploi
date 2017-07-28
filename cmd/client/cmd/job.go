package cmd

import "github.com/spf13/cobra"

// jobCmd represents the job command
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Manages deployment jobs",
}

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists jobs",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(jobCmd)
	jobCmd.AddCommand(jobListCmd)

	jobListCmd.Flags().BoolP("pending", "p", true, "Show pending jobs vs. done jobs")
}
