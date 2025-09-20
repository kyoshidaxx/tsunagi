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
