package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xanzy/go-cloudstack/cloudstack"
	"os"
)

var (
	fs         *flag.FlagSet
	apiKey     string
	apiSecret  string
	hostScheme string
	hostAddr   string
	hostPath   string
	zoneName   string
	zoneID     string

	help bool

	client *cloudstack.CloudStackClient

	zd *ZoneDefinition
)

type (
	Fetcher func(*ZoneDefinition) error
)

func main() {
	var err error
	var zone *cloudstack.Zone
	var count int

	fs = flag.NewFlagSet("zone-cloner", flag.ContinueOnError)
	fs.StringVar(&apiKey, "key", "", "API Key")
	fs.StringVar(&apiSecret, "secret", "", "API Secret")
	fs.StringVar(&hostScheme, "scheme", "http", "HTTP Scheme to use (http or https)")
	fs.StringVar(&hostAddr, "host", "127.0.0.1:8080", "CloudStack Management host addr including port")
	fs.StringVar(&hostPath, "path", "/client/api", "API path")
	fs.StringVar(&zoneID, "zone-id", "", "ID of Zone to clone (mutually exclusive with zone-name)")
	fs.StringVar(&zoneName, "zone-name", "", "Name of Zone to clone (mutually exclusive with zone-id)")
	fs.BoolVar(&help, "help", false, "Show help")

	fmt.Println("Welcome to the CloudStack Zone Cloner")

	if err = fs.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error parsing input: %s\n", err)
		os.Exit(1)
	}

	validateArgs()

	client = cloudstack.NewAsyncClient(fmt.Sprintf("%s://%s%s", hostScheme, hostAddr, hostPath), apiKey, apiSecret, false)

	fmt.Println("Client created")

	if zoneID == "" {
		fmt.Println("Attempting to fetch zone " + zoneName)
		zone, count, err = client.Zone.GetZoneByName(zoneName)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
		if count == 0 {
			fmt.Println("Error: Zone " + zoneName + " not found")
			os.Exit(1)
		}
	} else {
		fmt.Println("Attempting to fetch zone " + zoneID)
		zone, count, err = client.Zone.GetZoneByID(zoneID)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
		if count == 0 {
			fmt.Println("Error: Zone " + zoneID + " not found")
			os.Exit(1)
		}
	}

	fmt.Println("Zone found")
	fmt.Println("  ID: " + zone.Id)
	fmt.Println("  Name: " + zone.Name)
	fmt.Println("  Description: " + zone.Description)

	fmt.Println("")
	fmt.Println("Reading configuration...")

	zd = NewZoneDefinition(*zone)

	fetchers := []Fetcher{
		fetchPods,
		fetchClusters,
		fetchHosts,
		fetchPrimaryStoragePools,
		fetchSecondaryStoragePools,
		fetchPhysicalNetworks,
		fetchComputeOfferings,
		fetchDiskOfferings,
		fetchGlobalConfigs,
	}

	for _, fetcher := range fetchers {
		if err = fetcher(zd); err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
	}

	b, err := json.MarshalIndent(zd, "", "\t")
	fmt.Println(err)
	fmt.Println(string(b))

}
