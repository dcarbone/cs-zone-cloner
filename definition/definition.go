package definition

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

const (
	DefaultScheme  = "http"
	DefaultAddress = "127.0.0.1:8080"
	DefaultPath    = "/client/api"

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

		Database DatabaseConfig

		// Custom can be used by whatever custom fetchers you define
		// it is recommended you define your own formatter to take advantage of these.
		Custom map[string]interface{} `json:"-"`
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

		Custom: make(map[string]interface{}),
	}
	return zd
}

type (
	Config struct {
		Key      string `json:"key"`
		Secret   string `json:"secret"`
		Scheme   string `json:"scheme"`
		Address  string `json:"address"`
		Path     string `json:"path"`
		ZoneID   string `json:"zoneID"`
		ZoneName string `json:"zoneName"`

		Database *DatabaseConfig

		Fetchers []Fetcher `json:"-"`
	}
)

func FetchDefinition(conf Config, dbConfig *DatabaseConfig) (*ZoneDefinition, error) {
	var zone *cloudstack.Zone
	var count int
	var err error

	key := conf.Key
	secret := conf.Secret
	scheme := conf.Scheme
	address := conf.Address
	path := conf.Path
	zoneName := conf.ZoneName
	zoneID := conf.ZoneID
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}
	if secret == "" {
		return nil, errors.New("secret cannot be empty")
	}
	if scheme != "" && conf.Scheme != "http" && conf.Scheme != "https" {
		return nil, errors.New("scheme must be http or https")
	} else if scheme == "" {
		scheme = DefaultScheme
	}
	if address == "" {
		address = DefaultAddress
	}
	if path == "" {
		path = DefaultPath
	}
	if zoneName == "" && zoneID == "" {
		return nil, errors.New("zone name or id must be populated")
	}

	client := cloudstack.NewAsyncClient(fmt.Sprintf("%s://%s%s", scheme, address, path), key, secret, false)
	if zoneID == "" {
		log.Println("Attempting to fetch zone " + zoneName)
		zone, count, err = client.Zone.GetZoneByName(zoneName)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, fmt.Errorf("zone " + zoneName + " not found")
		}
	} else {
		log.Println("Attempting to fetch zone " + zoneID)
		zone, count, err = client.Zone.GetZoneByID(zoneID)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, fmt.Errorf("zone " + zoneID + " not found")
		}
	}

	zd := NewZoneDefinition(*zone)

	var fetchers []Fetcher

	if len(conf.Fetchers) == 0 {
		fetchers = defaultFetchers
	} else {
		fetchers = conf.Fetchers
	}

	for _, fetcher := range fetchers {
		if err = fetcher.Fetch(client, zd); err != nil {
			return nil, err
		}
	}

	if dbConfig != nil {
		zd.Database = *dbConfig
	}

	return zd, nil
}
