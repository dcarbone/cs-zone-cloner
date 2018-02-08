package main

import "fmt"

func fetchPods(zd *ZoneDefinition) error {
	fmt.Println("Fetching Pods...")
	params := client.Pod.NewListPodsParams()
	params.SetZoneid(zd.Zone.Id)
	pods, err := client.Pod.ListPods(params)
	if err != nil {
		return err
	}
	fmt.Println("Pods fetched")
	for _, pod := range pods.Pods {
		zd.Pods[pod.Name] = *pod
		fmt.Println("  Pod: " + pod.Name)
	}
	return nil
}

func fetchClusters(zd *ZoneDefinition) error {
	fmt.Println("Fetching Clusters...")
	params := client.Cluster.NewListClustersParams()
	params.SetZoneid(zd.Zone.Id)
	clusters, err := client.Cluster.ListClusters(params)
	if err != nil {
		return err
	}
	fmt.Println("Clusters fetched")
	for _, cluster := range clusters.Clusters {
		zd.Clusters[cluster.Name] = *cluster
		fmt.Println("  Cluster: " + cluster.Name)
	}
	return nil
}

func fetchHosts(zd *ZoneDefinition) error {
	fmt.Println("Fetching Hosts...")
	params := client.Host.NewListHostsParams()
	params.SetZoneid(zd.Zone.Id)
	hosts, err := client.Host.ListHosts(params)
	if err != nil {
		return err
	}
	fmt.Println("Hosts fetched")
	for _, host := range hosts.Hosts {
		zd.Hosts[host.Name] = *host
		fmt.Println("  Host: " + host.Name)
	}
	return nil
}

func fetchPrimaryStoragePools(zd *ZoneDefinition) error {
	fmt.Println("Fetching Primary Storage Pools...")
	params := client.Pool.NewListStoragePoolsParams()
	params.SetZoneid(zd.Zone.Id)
	pools, err := client.Pool.ListStoragePools(params)
	if err != nil {
		return err
	}
	fmt.Println("Fetched Primary Storage Pools")
	for _, pool := range pools.StoragePools {
		zd.PrimaryStoragePools[pool.Name] = *pool
		fmt.Println("  Pool: " + pool.Name)
	}
	return nil
}

func fetchSecondaryStoragePools(zd *ZoneDefinition) error {
	fmt.Println("Fetching Secondary (Image) Storage Pools...")
	params := client.ImageStore.NewListImageStoresParams()
	params.SetZoneid(zd.Zone.Id)
	pools, err := client.ImageStore.ListImageStores(params)
	if err != nil {
		return err
	}
	fmt.Println("Secondary (Image) Storage Pools fetched")
	for _, pool := range pools.ImageStores {
		zd.SecondaryStoragePools[pool.Name] = *pool
		fmt.Println("  Pool: " + pool.Name)
	}
	return nil
}

func fetchPhysicalNetworks(zd *ZoneDefinition) error {
	var err error
	fmt.Println("Fetching Physical Networks...")
	params := client.Network.NewListPhysicalNetworksParams()
	params.SetZoneid(zd.Zone.Id)
	cspns, err := client.Network.ListPhysicalNetworks(params)
	if err != nil {
		return err
	}
	fmt.Println("Physical Networks fetched")
	for _, cspn := range cspns.PhysicalNetworks {
		if zd.PhysicalNetworks[cspn.Name], err = expandPhysicalNetwork(zd, cspn); err != nil {
			return err
		}
	}
	return nil
}

func fetchComputeOfferings(zd *ZoneDefinition) error {
	fmt.Println("Fetching Compute Offerings...")
	params := client.ServiceOffering.NewListServiceOfferingsParams()
	params.SetIsrecursive(true)
	params.SetIssystem(false)
	params.SetListall(true)
	offerings, err := client.ServiceOffering.ListServiceOfferings(params)
	if err != nil {
		return err
	}
	fmt.Println("Compute Offerings fetched")
	for _, offering := range offerings.ServiceOfferings {
		zd.ComputeOfferings[offering.Name] = *offering
		fmt.Println("  Offering: " + offering.Name)
	}
	return nil
}

func fetchDiskOfferings(zd *ZoneDefinition) error {
	fmt.Println("Fetching Disk Offerings...")
	params := client.DiskOffering.NewListDiskOfferingsParams()
	params.SetIsrecursive(true)
	params.SetListall(true)
	offerings, err := client.DiskOffering.ListDiskOfferings(params)
	if err != nil {
		return err
	}
	fmt.Println("Disk Offerings fetched")
	for _, offering := range offerings.DiskOfferings {
		zd.DiskOfferings[offering.Name] = *offering
		fmt.Println("  Offering: " + offering.Name)
	}
	return nil
}

func fetchGlobalConfigs(zd *ZoneDefinition) error {
	fmt.Println("Fetching Global Configuration...")
	params := client.Configuration.NewListConfigurationsParams()
	configs, err := client.Configuration.ListConfigurations(params)
	if err != nil {
		return err
	}
	fmt.Println("Global Configuration fetched")
	for _, config := range configs.Configurations {
		zd.GlobalConfigs[config.Name] = *config
	}
	fmt.Println("Fetching Zone-specific Configuration...")
	params.SetZoneid(zd.Zone.Id)
	configs, err = client.Configuration.ListConfigurations(params)
	if err != nil {
		return err
	}
	fmt.Println("Zone-specific Configuration fetched")
	for _, config := range configs.Configurations {
		if c, ok := zd.GlobalConfigs[config.Name]; ok {
			fmt.Printf("  Overwriting Global Config %v value %v with %v\n", config.Name, c.Value, config.Value)
		} else {
			fmt.Printf("  Setting Config %v to %v\n", config.Name, config.Value)
		}
		zd.GlobalConfigs[config.Name] = *config
	}
	return nil
}
