package definition

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

type (
	Config struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
		Scheme string `json:"scheme"`
		Host   string `json:"host"`
		Path   string `json:"path"`

		ZoneID   string `json:"zoneID,omitempty"`
		ZoneName string `json:"zoneName,omitempty"`
		AllZones bool   `json:"allZones,omitempty"`

		DomainID   string `json:"domainID,omitempty"`
		DomainName string `json:"domainName,omitempty"`
		AllDomains bool   `json:"allDomains,omitempty"`

		Database *DatabaseConfig `json:"database,omitempty"`

		Fetchers []Fetcher `json:"-"`
	}
)

func FetchDefinition(conf Config, dbConfig *DatabaseConfig) (*RegionDefinition, error) {
	var zone *cloudstack.Zone
	var count int
	var err error

	key := conf.Key
	secret := conf.Secret
	scheme := conf.Scheme
	host := conf.Host
	path := conf.Path
	zoneName := conf.ZoneName
	zoneID := conf.ZoneID
	domainName := conf.DomainName
	domainID := conf.DomainID
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
	if host == "" {
		host = DefaultHost
	}
	if path == "" {
		path = DefaultPath
	}
	if zoneName != "" && zoneID != "" {
		return nil, errors.New("cannot define zone id and zone name at once")
	}
	if domainName != "" && domainID != "" {
		return nil, errors.New("cannot define domain id and domain name at once")
	}

	client := cloudstack.NewAsyncClient(fmt.Sprintf("%s://%s%s", scheme, host, path), key, secret, false)
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
