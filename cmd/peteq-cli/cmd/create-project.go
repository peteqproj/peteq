package cmd

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/spf13/cobra"
)

var createProjectFlags struct {
	description string
	imageURL    string
}

var createProjectCmd = &cobra.Command{
	Use:   "project ...names",
	Short: "Create project",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cnf, auth, err := createClientConfiguration()
		if err != nil {
			return err
		}
		c := client.NewAPIClient(cnf)
		for _, t := range args {
			c.ProjectCommandAPIApi.CProjectCreatePost(auth).Body(client.CreateProjectRequestBody{
				Name:        t,
				Color:       utils.PtrString(colors[randomdata.Number(0, len(colors))]),
				Description: &createProjectFlags.description,
				ImageUrl:    &createProjectFlags.imageURL,
			}).Execute()
			if err != nil {
				return err
			}
			fmt.Printf("Project %s created\n", t)
		}
		return nil
	},
}

func init() {
	createCmd.AddCommand(createProjectCmd)
	createProjectCmd.Flags().StringVar(&createProjectFlags.description, "description", "", "Project description")
	createProjectCmd.Flags().StringVar(&createProjectFlags.imageURL, "image", "", "Image URL to set")
}
