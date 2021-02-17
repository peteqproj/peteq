package cmd

import (
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/spf13/cobra"
)

var importRSSFlags struct {
	url string
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

		for _, i := range rss.Items {
			fmt.Println(i.Title)
			_, _, err := c.TaskCommandAPIApi.CTaskCreatePost(auth).Body(client.CreateTaskResponse{
				Name:        i.Title,
				Description: utils.PtrString(fmt.Sprintf("Link: %s\n%s", i.Link, i.Description)),
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
}
