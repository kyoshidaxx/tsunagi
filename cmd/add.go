/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	f "github.com/kyoshidaxx/tsunagi/internal/datastore/file"
	"github.com/kyoshidaxx/tsunagi/internal/domain/config"
	"github.com/kyoshidaxx/tsunagi/internal/utils"
	"github.com/spf13/cobra"
)

var projectID string
var region string
var instanceName string
var port int
var name string

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

		if projectID == "" {
			prompt := &survey.Input{
				Message: "Enter Project ID",
			}
			err := survey.AskOne(prompt, &projectID)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		if region == "" {
			prompt := &survey.Select{
				Message: "Select Region",
				Options: utils.GetRegionList(),
			}
			err := survey.AskOne(prompt, &region)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		if instanceName == "" {
			prompt := &survey.Input{
				Message: "Enter Instance Name",
			}
			err := survey.AskOne(prompt, &instanceName)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		if port == 0 {
			prompt := &survey.Input{
				Message: "Enter Bind Port",
			}
			err := survey.AskOne(prompt, &port)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		if name == "" {
			prompt := &survey.Input{
				Message: "Enter Config Name",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		r := f.NewConfigFileRepository(os.Getenv("CONFIG_FILE_PATH"))
		c := config.NewConfig(r)

		err = c.Add(config.ConfigParam{
			Name:         name,
			Port:         port,
			ProjectName:  projectID,
			Region:       region,
			InstanceName: instanceName,
		})

		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID")
	addCmd.Flags().StringVarP(&region, "region", "r", "", "Region")
	addCmd.Flags().StringVarP(&instanceName, "instance", "i", "", "Instance name")
	addCmd.Flags().IntVarP(&port, "port", "o", 0, "Port")
	addCmd.Flags().StringVarP(&name, "name", "n", "", "Name")
}
