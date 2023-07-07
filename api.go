package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"time"
)

func fetchAllOrgs(ctx context.Context, client *tfe.Client) ([]string, error) {
	var orgs []string
	currentPage := 1

	for {
		options := &tfe.OrganizationListOptions{
			ListOptions: tfe.ListOptions{
				PageSize:   25,
				PageNumber: currentPage,
			},
		}

		result, err := client.Organizations.List(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("could not fetch organizations: %w", err)
		}

		for _, org := range result.Items {
			orgs = append(orgs, org.Name)
		}

		if result.NextPage == 0 {
			break
		}

		currentPage++
	}

	return orgs, nil
}

func fetchAllWorkspaces(ctx context.Context, client *tfe.Client, org string) ([]*tfe.Workspace, error) {
	var workspaces []*tfe.Workspace
	currentPage := 1

	for {
		options := &tfe.WorkspaceListOptions{
			ListOptions: tfe.ListOptions{
				PageSize:   25,
				PageNumber: currentPage,
			},
		}

		result, err := client.Workspaces.List(ctx, org, options)
		if err != nil {
			return nil, fmt.Errorf("could not fetch workspaces: %w", err)
		}

		for _, workspace := range result.Items {
			workspaces = append(workspaces, workspace)
		}

		if result.NextPage == 0 {
			break
		}

		currentPage++
	}

	return workspaces, nil
}

func fetchAllResources(ctx context.Context, client *tfe.Client, workspaceId string) ([]*tfe.StateVersionResources, error) {
	state, err := client.StateVersions.ReadCurrent(ctx, workspaceId)
	if err != nil {
		return nil, fmt.Errorf("could not fetch current state: %w", err)
	}

	for {
		if state.ResourcesProcessed {
			break
		}
		time.Sleep(time.Second)
	}

	return state.Resources, nil
}
