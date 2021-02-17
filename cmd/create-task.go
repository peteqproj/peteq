package cmd

import (
	"fmt"

	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/spf13/cobra"
)

var createTaskFlags struct {
	project     string
	list        string
	description string
}

var createTaskCmd = &cobra.Command{
	Use:   "task ...names",
	Short: "Create task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := ""
		listID := ""
		logr := logger.New(logger.Options{})
		cnf, auth, err := createClientConfiguration()
		if err != nil {
			return err
		}
		c := client.NewAPIClient(cnf)
		if createTaskFlags.project != "" {
			logr.Info("Requesting projects", "name", createTaskFlags.project)
			projects, _, err := c.RestAPIApi.ApiProjectGet(auth).Execute()
			if err != nil {
				return err
			}
			for _, p := range projects {
				if *p.Metadata.Name == createTaskFlags.project {
					logr.Info("Project matched to ID", "id", *p.Metadata.Id, "name", *p.Metadata.Name)
					projectID = *p.Metadata.Id
				}
			}
		}
		if createTaskFlags.list != "" {
			logr.Info("Requesting list", "name", createTaskFlags.list)
			lists, _, err := c.RestAPIApi.ApiListGet(auth).Execute()
			if err != nil {
				return err
			}
			for _, l := range lists {
				if *l.Metadata.Name == createTaskFlags.list {
					logr.Info("List matched to ID", "id", *l.Metadata.Id, "name", *l.Metadata.Name)
					listID = *l.Metadata.Id
				}
			}
		}
		for _, t := range args {
			_, _, err := c.TaskCommandAPIApi.CTaskCreatePost(auth).Body(client.CreateTaskResponse{
				Name:        t,
				Project:     &projectID,
				List:        &listID,
				Description: &createTaskFlags.description,
			}).Execute()
			if err != nil {
				return err
			}
			fmt.Printf("Task %s created\n", t)
		}
		return nil
	},
}

func init() {
	createCmd.AddCommand(createTaskCmd)
	createTaskCmd.Flags().StringVar(&createTaskFlags.project, "project", "", "Assign the task to project")
	createTaskCmd.Flags().StringVar(&createTaskFlags.list, "list", "", "Assign the task to list")
	createTaskCmd.Flags().StringVar(&createTaskFlags.description, "description", "", "Set task description (if multiple tasks are create, the same description will be set for all)")
}
