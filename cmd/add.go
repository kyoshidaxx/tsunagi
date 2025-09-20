/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/kyoshidaxx/tsunagi/internal/domain/cloud"
	"github.com/kyoshidaxx/tsunagi/internal/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Save connection information to file",
	Long: `Save connection information to a YML file
Use Cloud SQL Auth Proxy based on the saved connection information`,
	Run: func(cmd *cobra.Command, args []string) {
		// gcloud command check
		err := utils.CheckGcloudCmd()
		if err != nil {
			return
		}
		// gcloud auth check
		err = utils.CheckGcloudAuth()
		if err != nil {
			return
		}

		// get information
		ctx := context.Background()
		pc, err := cloud.NewProjectClient(ctx)
		if err != nil {
			return
		}
		defer pc.Close()
		var projects []cloud.Project
		projects, err = pc.GetProjectList(ctx)
		if err != nil {
			return
		}

		// select project
		prompt := promptui.Select{
			Label: "Select Project",
			Items: projects,
		}
		index, _, err := prompt.Run()
		if err != nil {
			return
		}
		projectId := projects[index].ID

		fmt.Println(projectId)
		// todo db instance の選択

		// r := datastore.NewConfigFileRepository()
		// c := config.NewConfig(r)

		// c.Add(config.ConfigParam{
		// 	Name:         args[0],
		// 	Port:         strconv.Atoi(args[1]),
		// 	ProjectName:  args[2],
		// 	Region:       args[3],
		// 	InstanceName: args[4],
		// })
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
