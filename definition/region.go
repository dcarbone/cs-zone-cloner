package definition

import "github.com/xanzy/go-cloudstack/cloudstack"

type (
	RegionDefinition struct {
		Zones   map[string]*ZoneDefinition
		Domains map[string]*DomainDefinition

		Database DatabaseConfig

		// Custom can be used by whatever custom fetchers you define
		// it is recommended you define your own formatter to take advantage of these.
		Custom map[string]interface{} `json:"-"`
	}
)

func NewRegionDefinition(region cloudstack.Region) *RegionDefinition {
	rd := &RegionDefinition{
		Zones:   make(map[string]*ZoneDefinition),
		Domains: make(map[string]*DomainDefinition),

		Custom: make(map[string]interface{}),
	}
	return rd
}
