package googleapis

import (
	"context"
	"fmt"

	"google.golang.org/api/dns/v1"
	"google.golang.org/api/serviceusage/v1beta1"
)

const (
	ManagedZoneVisibilityUnknown = iota
	ManagedZoneVisibilityPublic  = iota
	ManagedZoneVisibilityPrivate = iota
)

// Client represents a client to the Google APIs.
type Client struct {
	ctx       context.Context
	projectId string
}

// New creates a new Client instance.
func New(ctx context.Context, projectId string) *Client {
	return &Client{
		ctx:       ctx,
		projectId: projectId,
	}
}

// GetServiceStates gets a list of services along with their status.
func (c *Client) GetServiceStates() (map[string]bool, error) {
	svc, err := serviceusage.NewService(c.ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := svc.Services.List(fmt.Sprintf("projects/%s", c.projectId)).
		Filter("state:ENABLED").
		Do()
	if err != nil {
		return nil, err
	}

	states := map[string]bool{}
	for _, svc := range rsp.Services {
		states[svc.Config.Name] = svc.State == "ENABLED"
	}

	return states, nil
}

// GetManagedZones fetches the list of managed zones along with their visibility.
func (c *Client) GetManagedZones() (map[string]int, error) {
	client, err := dns.NewService(c.ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := client.ManagedZones.List(c.projectId).Do()
	if err != nil {
		return nil, err
	}

	zones := make(map[string]int, len(rsp.ManagedZones))
	for _, zone := range rsp.ManagedZones {
		visibility := ManagedZoneVisibilityUnknown
		switch zone.Visibility {
		case "public":
			visibility = ManagedZoneVisibilityPublic
		case "private":
			visibility = ManagedZoneVisibilityPrivate
		}

		zones[zone.DnsName] = visibility
	}

	return zones, nil
}
