package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"strings"
)

const (
	freeTierThreshold  int     = 500
	hourlyResourceCost float64 = 0.00014
)

var unmanagedResourcePrefixes = []string{
	"data.",
	"null_resource.",
}

func main() {
	token, err := getApiToken()
	if err != nil {
		log.Fatal(err)
	}

	cfg := &tfe.Config{
		Token: token,
	}

	client, err := tfe.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	orgs, err := fetchAllOrgs(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}

	org, err := chooseOrganization(orgs)
	if err != nil {
		log.Fatal(err)
	}

	workspaces, err := fetchAllWorkspaces(context.Background(), client, org)
	if err != nil {
		log.Fatal(err)
	}

	resourceCount := 0
	managedResourceCount := 0

	for _, workspace := range workspaces {
		resources, err := fetchAllResources(context.Background(), client, workspace.ID)
		if err != nil {
			log.Fatal(err)
		}

		for _, resource := range resources {
			if isManagedResource(resource.Type) {
				managedResourceCount += resource.Count
			}

			resourceCount += resource.Count
		}
	}

	managedResourceCount -= freeTierThreshold
	if managedResourceCount < 0 {
		managedResourceCount = 0
	}

	hourlyCost := hourlyResourceCost * float64(managedResourceCount)
	dailyCost := hourlyCost * float64(24)
	monthlyCost := dailyCost * float64(31)
	yearlyCost := monthlyCost * float64(12)

	fmt.Printf("%s has %d workspaces\n", org, len(workspaces))
	fmt.Printf("Total Resources: %d\n", resourceCount)
	fmt.Printf("Total Managed Resources: %d\n", managedResourceCount)
	fmt.Printf("Total Hourly Cost: $%f\n", hourlyCost)
	fmt.Printf("Total Daily Cost: $%.2f\n", dailyCost)
	fmt.Printf("Total Total Monthly Cost: $%.2f\n", monthlyCost)
	fmt.Printf("Total Yearly Cost: $%.2f\n", yearlyCost)
}

func getApiToken() (string, error) {
	if token := os.Getenv("TFE_TOKEN"); token != "" {
		return token, nil
	}

	prompt := promptui.Prompt{
		Label: "Please enter your Terraform Cloud API Token",
		Mask:  '*',
	}

	return prompt.Run()
}

func chooseOrganization(orgs []string) (string, error) {
	if len(orgs) == 1 {
		return orgs[0], nil
	}

	prompt := promptui.Select{
		Label: "Please choose an Organization",
		Items: orgs,
	}

	_, org, err := prompt.Run()
	return org, err
}

func isManagedResource(resourceType string) bool {
	for _, prefix := range unmanagedResourcePrefixes {
		if strings.HasPrefix(resourceType, prefix) {
			return false
		}
	}

	return true
}
