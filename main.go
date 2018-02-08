package main

import (
	"flag"
	"fmt"
	"github.com/dcarbone/cs-zone-cloner/definition"
	"os"
	"strings"
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

	format string
	output string
)

func validateArgs() {
	ok := true

	if apiKey == "" {
		fmt.Println("key cannot be empty")
		ok = false
	}
	if apiSecret == "" {
		fmt.Println("secret cannot be empty")
		ok = false
	}
	hostScheme = strings.ToLower(hostScheme)
	if hostScheme != "http" && hostScheme != "https" {
		fmt.Println("scheme must be \"http\" or \"https\"")
		ok = false
	}
	if hostAddr == "" {
		fmt.Println("host cannot be empty")
		ok = false
	}
	if hostPath == "" {
		fmt.Println("path cannot be empty")
		ok = false
	}
	if zoneName == "" && zoneID == "" {
		fmt.Println("zone-id or zone-name must be set")
		ok = false
	}
	format = strings.ToLower(format)
	if format != "json" && format != "yaml" {
		fmt.Println("format must be json or yaml")
		ok = false
	}

	if !ok {
		os.Exit(1)
	}

	log.Println("Using parameters:")
	log.Println("  APIKey: " + apiKey)
	log.Println("  APISecret: " + apiSecret)
	log.Println("  HostScheme: " + hostScheme)
	log.Println("  HostAddr: " + hostAddr)
	log.Println("  HostPath: " + hostPath)
	if zoneID == "" {
		log.Println("  ZoneName: " + zoneName)
	} else {
		log.Println("  ZoneID: " + zoneID)
	}
	log.Println("  Format: " + format)
	if output != "" {
		log.Println("  Output: " + output)
	}
}

func main() {
	var err error

	fs = flag.NewFlagSet("zone-cloner", flag.ContinueOnError)
	fs.StringVar(&apiKey, "key", "", "API Key")
	fs.StringVar(&apiSecret, "secret", "", "API Secret")
	fs.StringVar(&hostScheme, "scheme", "http", "HTTP Scheme to use (http or https)")
	fs.StringVar(&hostAddr, "host", "127.0.0.1:8080", "CloudStack Management host addr including port")
	fs.StringVar(&hostPath, "path", "/client/api", "API path")
	fs.StringVar(&zoneID, "zone-id", "", "ID of Zone to clone (mutually exclusive with zone-name)")
	fs.StringVar(&zoneName, "zone-name", "", "Name of Zone to clone (mutually exclusive with zone-id)")
	fs.StringVar(&output, "output", "", "File to write to")

	if err = fs.Parse(os.Args[1:]); err != nil {
		fmt.Println("Error parsing input: " + err.Error())
		os.Exit(1)
	}

	validateArgs()

	log.Println("Fetching definition...")

	zd, err := definition.FetchDefinition(definition.Config{
		Key:      apiKey,
		Secret:   apiSecret,
		Scheme:   hostScheme,
		Address:  hostAddr,
		Path:     hostPath,
		ZoneName: zoneName,
		ZoneID:   zoneID,
	})

	log.Println("")

	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	log.Println("Definition built")

	if format == "json" {
		b, err := definition.FormatJSONIndent(zd)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
		if output == "" {
			fmt.Println(string(b))
		} else {
			f, err := os.Create(output)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				os.Exit(1)
			}
			f.Write(b)
			f.Close()
		}
	}

	os.Exit(0)
}
