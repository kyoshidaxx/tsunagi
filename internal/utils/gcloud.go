package utils

import (
	"fmt"
	"os/exec"
)

func CheckGcloudCmd() error {
	cmd := exec.Command("gcloud", "--version")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("gcloud command not found")
		return err
	}
	return nil
}

func CheckGcloudAuth() error {
	cmd := exec.Command("gcloud", "auth", "application-default", "print-access-token")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Your gcloud credentials are invalid or have expired.")
		fmt.Println("Please reauthenticate by executing the following command.")
		fmt.Println("")
		fmt.Println("  gcloud auth application-default login")
		fmt.Println("")
		return err
	}
	return nil
}

func GetRegionList() []string {
	return []string{
		"asia-east1",              // Changhua County, Taiwan
		"asia-east2",              // Hong Kong
		"asia-northeast1",         // Tokyo, Japan
		"asia-northeast2",         // Osaka, Japan
		"asia-northeast3",         // Seoul, South Korea
		"asia-south1",             // Mumbai, India
		"asia-south2",             // Delhi, India
		"asia-southeast1",         // Jurong West, Singapore
		"asia-southeast2",         // Jakarta, Indonesia
		"australia-southeast1",    // Sydney, Australia
		"australia-southeast2",    // Melbourne, Australia
		"europe-central2",         // Warsaw, Poland
		"europe-north1",           // Hamina, Finland
		"europe-southwest1",       // Madrid, Spain
		"europe-west1",            // St. Ghislain, Belgium
		"europe-west2",            // London, UK
		"europe-west3",            // Frankfurt, Germany
		"europe-west4",            // Eemshaven, Netherlands
		"europe-west6",            // Zurich, Switzerland
		"northamerica-northeast1", // Montreal, Canada
		"northamerica-northeast2", // Toronto, Canada
		"southamerica-east1",      // São Paulo, Brazil
		"southamerica-west1",      // Santiago, Chile
		"us-central1",             // Council Bluffs, Iowa, USA
		"us-east1",                // Moncks Corner, South Carolina, USA
		"us-east4",                // Ashburn, Virginia, USA
		"us-south1",               // Dallas, Texas, USA
		"us-west1",                // The Dalles, Oregon, USA
		"us-west2",                // Los Angeles, California, USA
		"us-west3",                // Salt Lake City, Utah, USA
		"us-west4",                // Las Vegas, Nevada, USA
	}
}
