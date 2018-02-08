package main

import "github.com/xanzy/go-cloudstack/cloudstack"

type (
	TrafficType struct {
		cloudstack.TrafficType
		Networks map[string]cloudstack.Network
	}
	PhysicalNetwork struct {
		cloudstack.PhysicalNetwork
		TrafficTypes map[string]TrafficType
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
		GlobalConfigs         map[string]cloudstack.Configuration
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
		GlobalConfigs:         make(map[string]cloudstack.Configuration),
	}
	return zd
}
