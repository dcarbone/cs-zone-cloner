package definition

import (
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func fetchPods(zd *ZoneDefinition) error {
	log.Println("Fetching Pods...")
	params := client.Pod.NewListPodsParams()
	params.SetZoneid(zd.Zone.Id)
	pods, err := client.Pod.ListPods(params)
	if err != nil {
		return err
	}
	log.Println("Pods fetched")
	for _, pod := range pods.Pods {
		zd.Pods[pod.Name] = *pod
		log.Println("  Pod: " + pod.Name)
	}
	return nil
}

func fetchClusters(zd *ZoneDefinition) error {
	log.Println("Fetching Clusters...")
	params := client.Cluster.NewListClustersParams()
	params.SetZoneid(zd.Zone.Id)
	clusters, err := client.Cluster.ListClusters(params)
	if err != nil {
		return err
	}
	log.Println("Clusters fetched")
	for _, cluster := range clusters.Clusters {
		zd.Clusters[cluster.Name] = *cluster
		log.Println("  Cluster: " + cluster.Name)
	}
	return nil
}

func fetchHosts(zd *ZoneDefinition) error {
	log.Println("Fetching Hosts...")
	params := client.Host.NewListHostsParams()
	params.SetZoneid(zd.Zone.Id)
	hosts, err := client.Host.ListHosts(params)
	if err != nil {
		return err
	}
	log.Println("Hosts fetched")
	for _, host := range hosts.Hosts {
		zd.Hosts[host.Name] = *host
		log.Println("  Host: " + host.Name)
	}
	return nil
}

func fetchPrimaryStoragePools(zd *ZoneDefinition) error {
	log.Println("Fetching Primary Storage Pools...")
	params := client.Pool.NewListStoragePoolsParams()
	params.SetZoneid(zd.Zone.Id)
	pools, err := client.Pool.ListStoragePools(params)
	if err != nil {
		return err
	}
	log.Println("Fetched Primary Storage Pools")
	for _, pool := range pools.StoragePools {
		zd.PrimaryStoragePools[pool.Name] = *pool
		log.Println("  Pool: " + pool.Name)
	}
	return nil
}

func fetchSecondaryStoragePools(zd *ZoneDefinition) error {
	log.Println("Fetching Secondary (Image) Storage Pools...")
	params := client.ImageStore.NewListImageStoresParams()
	params.SetZoneid(zd.Zone.Id)
	pools, err := client.ImageStore.ListImageStores(params)
	if err != nil {
		return err
	}
	log.Println("Secondary (Image) Storage Pools fetched")
	for _, pool := range pools.ImageStores {
		zd.SecondaryStoragePools[pool.Name] = *pool
		log.Println("  Pool: " + pool.Name)
	}
	return nil
}

func expandTrafficType(zd *ZoneDefinition, csttype *cloudstack.TrafficType) (TrafficType, error) {
	var err error
	log.Println("    Expanding Traffic Type " + csttype.TrafficType + "...")
	ttype := &TrafficType{
		TrafficType: *csttype,
		Networks:    make(map[string]cloudstack.Network),
	}
	log.Println("    Fetching Traffic Type " + csttype.TrafficType + " Networks...")
	params := client.Network.NewListNetworksParams()
	params.SetZoneid(zd.Zone.Id)
	params.SetTraffictype(csttype.TrafficType)
	params.SetIssystem(true)
	params.SetListall(true)
	csnetworks, err := client.Network.ListNetworks(params)
	if err != nil {
		goto done
	}
	log.Println("    Traffic Type " + csttype.TrafficType + " Networks fetched")
	for _, network := range csnetworks.Networks {
		ttype.Networks[network.Name] = *network
		log.Println("      Network: " + network.Name)
	}
done:
	return *ttype, err
}

func expandPhysicalNetwork(zd *ZoneDefinition, cspn *cloudstack.PhysicalNetwork) (PhysicalNetwork, error) {
	var err error
	log.Println("  Expanding Physical Network " + cspn.Name + "...")
	ps := &PhysicalNetwork{
		PhysicalNetwork: *cspn,
		TrafficTypes:    make(map[string]TrafficType),
	}

	log.Println("  Fetching Physical Network " + cspn.Name + " Traffic Types...")
	csttypes, err := client.Usage.ListTrafficTypes(client.Usage.NewListTrafficTypesParams(cspn.Id))
	if err != nil {
		goto done
	}
	log.Println("  Physical Network " + cspn.Name + " Traffic Types fetched")
	for _, csttype := range csttypes.TrafficTypes {
		if ps.TrafficTypes[csttype.TrafficType], err = expandTrafficType(zd, csttype); err != nil {
			goto done
		}
	}

done:
	return *ps, err
}

func fetchPhysicalNetworks(zd *ZoneDefinition) error {
	var err error
	log.Println("Fetching Physical Networks...")
	params := client.Network.NewListPhysicalNetworksParams()
	params.SetZoneid(zd.Zone.Id)
	cspns, err := client.Network.ListPhysicalNetworks(params)
	if err != nil {
		return err
	}
	log.Println("Physical Networks fetched")
	for _, cspn := range cspns.PhysicalNetworks {
		if zd.PhysicalNetworks[cspn.Name], err = expandPhysicalNetwork(zd, cspn); err != nil {
			return err
		}
	}
	return nil
}

func fetchComputeOfferings(zd *ZoneDefinition) error {
	log.Println("Fetching Compute Offerings...")
	params := client.ServiceOffering.NewListServiceOfferingsParams()
	params.SetIsrecursive(true)
	params.SetIssystem(false)
	params.SetListall(true)
	offerings, err := client.ServiceOffering.ListServiceOfferings(params)
	if err != nil {
		return err
	}
	log.Println("Compute Offerings fetched")
	for _, offering := range offerings.ServiceOfferings {
		zd.ComputeOfferings[offering.Name] = *offering
		log.Println("  Offering: " + offering.Name)
	}
	return nil
}

func fetchDiskOfferings(zd *ZoneDefinition) error {
	log.Println("Fetching Disk Offerings...")
	params := client.DiskOffering.NewListDiskOfferingsParams()
	params.SetIsrecursive(true)
	params.SetListall(true)
	offerings, err := client.DiskOffering.ListDiskOfferings(params)
	if err != nil {
		return err
	}
	log.Println("Disk Offerings fetched")
	for _, offering := range offerings.DiskOfferings {
		zd.DiskOfferings[offering.Name] = *offering
		log.Println("  Offering: " + offering.Name)
	}
	return nil
}

func fetchGlobalConfigs(zd *ZoneDefinition) error {
	log.Println("Fetching Global Configuration...")
	params := client.Configuration.NewListConfigurationsParams()
	configs, err := client.Configuration.ListConfigurations(params)
	if err != nil {
		return err
	}
	log.Println("Global Configuration fetched")
	for _, config := range configs.Configurations {
		zd.GlobalConfigs[config.Name] = *config
	}
	log.Println("Fetching Zone-specific Configuration...")
	params.SetZoneid(zd.Zone.Id)
	configs, err = client.Configuration.ListConfigurations(params)
	if err != nil {
		return err
	}
	log.Println("Zone-specific Configuration fetched")
	for _, config := range configs.Configurations {
		if c, ok := zd.GlobalConfigs[config.Name]; ok {
			log.Printf("  Overwriting Global Config %v value %v with %v\n", config.Name, c.Value, config.Value)
		} else {
			log.Printf("  Setting Config %v to %v\n", config.Name, config.Value)
		}
		zd.GlobalConfigs[config.Name] = *config
	}
	return nil
}
