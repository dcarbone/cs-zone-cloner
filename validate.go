package main

import (
	"fmt"
	"os"
	"strings"
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

	if !ok {
		os.Exit(1)
	}

	fmt.Println("Using parameters:")
	fmt.Println("  APIKey: " + apiKey)
	fmt.Println("  APISecret: " + apiSecret)
	fmt.Println("  HostScheme: " + hostScheme)
	fmt.Println("  HostAddr: " + hostAddr)
	fmt.Println("  HostPath: " + hostPath)
	if zoneID == "" {
		fmt.Println("  ZoneName: " + zoneName)
	} else {
		fmt.Println("  ZoneID: " + zoneID)
	}
}
