package gcp

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type GCPHelper interface {
	ListAllPreferredZones(context.Context, string) ([]Zone, error)
}

type Zone struct {
	Name, Region string
}

type GCPHelperClient struct {
	compute *compute.Service
}

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func NewGCPHelperClient(ctx context.Context) (*GCPHelperClient, error) {
	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	computeService, err := compute.New(c)
	if err != nil {
		return nil, err
	}

	return &GCPHelperClient{
		compute: computeService,
	}, nil
}

func (c *GCPHelperClient) PreferredRegions() []string {
	return []string{
		"europe-north1",
		"europe-west1",
		"europe-west2",
		"europe-west3",
		"us-east1",
		"us-east4",
		"us-central1",
		"us-west1",
		"us-west2",
		"northamerica-northeast1",
	}
}

func (c *GCPHelperClient) ListAllPreferredZones(ctx context.Context, project string, region *string, avoid ...string) ([]Zone, error) {
	results := []Zone{}
	avoidZones := sets.NewString(avoid...)
	preferredRegions := sets.NewString()

	if region == nil {
		preferredRegions.Insert(c.PreferredRegions()...)
	} else {
		preferredRegions.Insert(*region)
	}

	err := c.compute.Zones.List(project).Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			if zone.Status == "UP" && preferredRegions.Has(zone.Region) && !avoidZones.Has(zone.Name) {
				results = append(results, Zone{Name: zone.Name, Region: zone.Region})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *GCPHelperClient) SelectRandomRegion(avoid ...string) string {
	regions := sets.NewString(c.PreferredRegions()...)
	regions.Delete(avoid...)
	return regions.UnsortedList()[rng.Intn(len(regions))]
}

func (c *GCPHelperClient) SelectRandomZone(ctx context.Context, project string, region *string, avoid ...string) (*Zone, error) {
	zones, err := c.ListAllPreferredZones(ctx, project, region, avoid...)
	if err != nil {
		return nil, err
	}
	return &zones[rng.Intn(len(zones))], nil
}
