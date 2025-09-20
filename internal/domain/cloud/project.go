package cloud

import (
	"context"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"google.golang.org/api/iterator"
)

type Project struct {
	ID   string
	Name string
}

type ProjectClient struct {
	client *resourcemanager.ProjectsClient
}

func NewProjectClient(ctx context.Context) (*ProjectClient, error) {
	client, err := resourcemanager.NewProjectsClient(ctx)
	if err != nil {
		return nil, err
	}
	return &ProjectClient{
		client: client,
	}, nil
}

func (c *ProjectClient) Close() error {
	return c.client.Close()
}

func (c *ProjectClient) GetProjectList(ctx context.Context) ([]Project, error) {
	it := c.client.ListProjects(ctx, &resourcemanagerpb.ListProjectsRequest{})
	var projects []Project
	for {
		project, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		projects = append(projects, Project{
			ID:   project.GetProjectId(),
			Name: project.GetName(),
		})
	}
	return projects, nil
}
