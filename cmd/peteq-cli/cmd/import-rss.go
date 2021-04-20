package cmd

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/mmcdole/gofeed"
	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/spf13/cobra"
)

var importRSSFlags struct {
	url     string
	project string
}

var importRSSCmd = &cobra.Command{
	Use:   "import-rss",
	Short: "Import RSS items as tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		logr := logger.New(logger.Options{})
		logr.Info("Pasing RSS", "url", importRSSFlags.url)
		rss, err := gofeed.NewParser().ParseURL(importRSSFlags.url)
		if err != nil {
			return err
		}
		cnf, auth, err := createClientConfiguration()
		if err != nil {
			return err
		}
		c := client.NewAPIClient(cnf)
		projectID := ""
		if importRSSFlags.project != "" {
			logr.Info("Requesting projects", "name", importRSSFlags.project)
			projects, _, err := c.RestAPIApi.ApiProjectGet(auth).Execute()
			if err != nil {
				return err
			}
			for _, p := range projects {
				if *p.Metadata.Name == importRSSFlags.project {
					logr.Info("Project matched to ID", "id", *p.Metadata.Id, "name", *p.Metadata.Name)
					projectID = *p.Metadata.Id
					break
				}
			}
			// project not found - create one
			if projectID == "" {
				resp, _, err := c.ProjectCommandAPIApi.CProjectCreatePost(auth).Body(client.CreateProjectRequestBody{
					Name:     importRSSFlags.project,
					Color:    utils.PtrString(colors[randomdata.Number(0, len(colors))]),
					ImageUrl: &rss.Image.URL,
				}).Execute()
				if err != nil {
					return err
				}
				projectID = *resp.Id
			}
		}
		for _, i := range rss.Items {
			_, _, err := c.TaskCommandAPIApi.CTaskCreatePost(auth).Body(client.CreateTaskResponse{
				Name:        i.Title,
				Description: utils.PtrString(fmt.Sprintf("Link: %s\n%s", i.Link, i.Description)),
				Project:     utils.PtrString(projectID),
			}).Execute()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importRSSCmd)
	importRSSCmd.Flags().StringVar(&importRSSFlags.url, "url", "", "URL of the RSS")
	importRSSCmd.Flags().StringVar(&importRSSFlags.project, "project", "", "Assign all tasks to given project (create the project if not exists)")
}
