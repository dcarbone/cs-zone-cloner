package main

import (
	"fmt"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func expandTrafficType(zd *ZoneDefinition, csttype *cloudstack.TrafficType) (TrafficType, error) {
	var err error
	fmt.Println("    Expanding Traffic Type " + csttype.TrafficType + "...")
	ttype := &TrafficType{
		TrafficType: *csttype,
		Networks:    make(map[string]cloudstack.Network),
	}
	fmt.Println("    Fetching Traffic Type " + csttype.TrafficType + " Networks...")
	params := client.Network.NewListNetworksParams()
	params.SetZoneid(zd.Zone.Id)
	params.SetTraffictype(csttype.TrafficType)
	params.SetIssystem(true)
	params.SetListall(true)
	csnetworks, err := client.Network.ListNetworks(params)
	if err != nil {
		goto done
	}
	fmt.Println("    Traffic Type " + csttype.TrafficType + " Networks fetched")
	for _, network := range csnetworks.Networks {
		ttype.Networks[network.Name] = *network
		fmt.Println("      Network: " + network.Name)
	}
done:
	return *ttype, err
}

func expandPhysicalNetwork(zd *ZoneDefinition, cspn *cloudstack.PhysicalNetwork) (PhysicalNetwork, error) {
	var err error
	fmt.Println("  Expanding Physical Network " + cspn.Name + "...")
	ps := &PhysicalNetwork{
		PhysicalNetwork: *cspn,
		TrafficTypes:    make(map[string]TrafficType),
	}

	fmt.Println("  Fetching Physical Network " + cspn.Name + " Traffic Types...")
	csttypes, err := client.Usage.ListTrafficTypes(client.Usage.NewListTrafficTypesParams(cspn.Id))
	if err != nil {
		goto done
	}
	fmt.Println("  Physical Network " + cspn.Name + " Traffic Types fetched")
	for _, csttype := range csttypes.TrafficTypes {
		if ps.TrafficTypes[csttype.TrafficType], err = expandTrafficType(zd, csttype); err != nil {
			goto done
		}
	}

done:
	return *ps, err
}
