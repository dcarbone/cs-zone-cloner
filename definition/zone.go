package definition

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

const (
	DefaultScheme = "http"
	DefaultHost   = "127.0.0.1:8080"
	DefaultPath   = "/client/api"

	DefaultDBHost = "localhost"
	DefaultDBPort = 3306
)

type (
	TrafficType struct {
		cloudstack.TrafficType
		Networks map[string]cloudstack.Network
	}
	PhysicalNetwork struct {
		cloudstack.PhysicalNetwork
		TrafficTypes map[string]TrafficType
	}

	DatabaseConfig struct {
		Server   string
		Port     int
		Schema   string
		User     string
		Password string
	}

	ZoneDefinition struct {
		Zone                  cloudstack.Zone
		Pods                  map[string]cloudstack.Pod
		Clusters              map[string]cloudstack.Cluster
		Hosts                 map[string]cloudstack.Host
		PrimaryStoragePools   map[string]cloudstack.StoragePool
		SecondaryStoragePools map[string]cloudstack.ImageStore
		PhysicalNetworks      map[string]PhysicalNetwork
		ComputeOfferings      map[string]cloudstack.ServiceOffering
		DiskOfferings         map[string]cloudstack.DiskOffering
		Templates             map[string]cloudstack.Template
		GlobalConfiguration   map[string]cloudstack.Configuration
		ZoneConfiguration     map[string]cloudstack.Configuration
	}
)

func NewZoneDefinition(zone cloudstack.Zone) *ZoneDefinition {
	zd := &ZoneDefinition{
		Zone:                  zone,
		Pods:                  make(map[string]cloudstack.Pod),
		Clusters:              make(map[string]cloudstack.Cluster),
		Hosts:                 make(map[string]cloudstack.Host),
		PrimaryStoragePools:   make(map[string]cloudstack.StoragePool),
		SecondaryStoragePools: make(map[string]cloudstack.ImageStore),
		PhysicalNetworks:      make(map[string]PhysicalNetwork),
		ComputeOfferings:      make(map[string]cloudstack.ServiceOffering),
		DiskOfferings:         make(map[string]cloudstack.DiskOffering),
		Templates:             make(map[string]cloudstack.Template),
		ZoneConfiguration:     make(map[string]cloudstack.Configuration),
		GlobalConfiguration:   make(map[string]cloudstack.Configuration),
	}
	return zd
}
