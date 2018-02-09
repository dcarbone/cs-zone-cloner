package definition

import (
	"github.com/xanzy/go-cloudstack/cloudstack"
	"sync"
)

var (
	registeredFetchers   map[string]Fetcher
	registeredFetchersMu sync.Mutex

	defaultFetchers []Fetcher
)

func init() {
	defaultFetchers = []Fetcher{
		new(FetchPods),
		new(FetchClusters),
		new(FetchHosts),
		new(FetchPrimaryStoragePools),
		new(FetchSecondaryStoragePools),
		new(FetchPhysicalNetworks),
		new(FetchComputeOfferings),
		new(FetchDiskOfferings),
		new(FetchTemplates),
		new(FetchZoneConfigurations),
		new(FetchGlobalConfigurations),
	}
	registeredFetchers = make(map[string]Fetcher, len(defaultFetchers))
	for _, df := range defaultFetchers {
		registeredFetchers[df.Name()] = df
	}
}

func DefaultFetchers() []string {
	fetchers := make([]string, len(defaultFetchers))
	for i, fetcher := range defaultFetchers {
		fetchers[i] = fetcher.Name()
	}
	return fetchers
}

func RegisterFetcher(f Fetcher) {
	registeredFetchersMu.Lock()
	registeredFetchers[f.Name()] = f
	registeredFetchersMu.Unlock()
}

func GetFetcher(name string) (Fetcher, bool) {
	registeredFetchersMu.Lock()
	f, ok := registeredFetchers[name]
	registeredFetchersMu.Unlock()
	return f, ok
}

type Fetcher interface {
	Name() string
	Fetch(*cloudstack.CloudStackClient, *ZoneDefinition) error
}

type FetchPods struct{}

func (*FetchPods) Name() string {
	return "pods"
}

func (*FetchPods) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
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

type FetchClusters struct{}

func (*FetchClusters) Name() string {
	return "clusters"
}

func (*FetchClusters) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
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

type FetchHosts struct{}

func (*FetchHosts) Name() string {
	return "hosts"
}

func (*FetchHosts) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
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

type FetchPrimaryStoragePools struct{}

func (*FetchPrimaryStoragePools) Name() string {
	return "primaryStoragePools"
}

func (*FetchPrimaryStoragePools) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
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

type FetchSecondaryStoragePools struct{}

func (*FetchSecondaryStoragePools) Name() string {
	return "secondaryStoragePools"
}

func (*FetchSecondaryStoragePools) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
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

type FetchPhysicalNetworks struct{}

func (*FetchPhysicalNetworks) Name() string {
	return "physicalNetworks"
}

func (*FetchPhysicalNetworks) expandTrafficType(client *cloudstack.CloudStackClient, zd *ZoneDefinition, csttype *cloudstack.TrafficType) (TrafficType, error) {
	var err error
	var key string
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
		if network.Name == "" {
			key = network.Id
		} else {
			key = network.Name
		}
		ttype.Networks[key] = *network
		log.Println("      Network: " + key)
	}
done:
	return *ttype, err
}

func (fpn *FetchPhysicalNetworks) expandPhysicalNetwork(client *cloudstack.CloudStackClient, zd *ZoneDefinition, cspn *cloudstack.PhysicalNetwork) (PhysicalNetwork, error) {
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
		if ps.TrafficTypes[csttype.TrafficType], err = fpn.expandTrafficType(client, zd, csttype); err != nil {
			goto done
		}
	}

done:
	return *ps, err
}

func (fpn *FetchPhysicalNetworks) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
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
		if zd.PhysicalNetworks[cspn.Name], err = fpn.expandPhysicalNetwork(client, zd, cspn); err != nil {
			return err
		}
	}
	return nil
}

type FetchComputeOfferings struct{}

func (*FetchComputeOfferings) Name() string {
	return "computeOfferings"
}

func (*FetchComputeOfferings) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
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

type FetchDiskOfferings struct{}

func (*FetchDiskOfferings) Name() string {
	return "diskOfferings"
}

func (*FetchDiskOfferings) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
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

type FetchTemplates struct{}

func (*FetchTemplates) Name() string {
	return "templates"
}

func (*FetchTemplates) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
	log.Println("Fetching Templates...")
	params := client.Template.NewListTemplatesParams("all")
	params.SetZoneid(zd.Zone.Id)
	params.SetIsrecursive(true)
	params.SetListall(true)
	templates, err := client.Template.ListTemplates(params)
	if err != nil {
		return err
	}
	log.Println("Templates fetched")
	for _, template := range templates.Templates {
		zd.Templates[template.Name] = *template
		log.Println("  Template: " + template.Name)
	}
	return nil
}

type FetchGlobalConfigurations struct{}

func (*FetchGlobalConfigurations) Name() string {
	return "globalConfigs"
}

func (*FetchGlobalConfigurations) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
	log.Println("Fetching Global Configurations...")
	params := client.Configuration.NewListConfigurationsParams()
	configs, err := client.Configuration.ListConfigurations(params)
	if err != nil {
		return err
	}
	log.Println("Global Configuration fetched")
	for _, config := range configs.Configurations {
		zd.GlobalConfiguration[config.Name] = *config
		log.Printf("  %s: %v", config.Name, config.Value)
	}
	return nil
}

type FetchZoneConfigurations struct{}

func (*FetchZoneConfigurations) Name() string {
	return "zoneConfigs"
}

func (*FetchZoneConfigurations) Fetch(client *cloudstack.CloudStackClient, zd *ZoneDefinition) error {
	log.Println("Fetching Zone Configurations...")
	log.Println("Fetching Zone-specific Configuration...")
	params := client.Configuration.NewListConfigurationsParams()
	params.SetZoneid(zd.Zone.Id)
	configs, err := client.Configuration.ListConfigurations(params)
	if err != nil {
		return err
	}
	log.Println("Zone-specific Configuration fetched")
	for _, config := range configs.Configurations {
		zd.ZoneConfiguration[config.Name] = *config
		log.Printf("  %s: %v", config.Name, config.Value)
	}
	return nil
}
