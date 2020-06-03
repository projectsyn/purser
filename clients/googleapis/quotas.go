package googleapis

import "google.golang.org/api/compute/v1"

// GetQuotas fetches the list of compute quotas.
func (c *Client) GetQuotas() ([]*compute.Quota, error) {
	client, _ := compute.NewService(c.ctx)
	project, err := client.Projects.Get(c.projectId).Do()
	if err != nil {
		return nil, err
	}

	return project.Quotas, nil
}
