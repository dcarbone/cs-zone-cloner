package backup

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dcarbone/cs-zone-cloner/command"
	"github.com/dcarbone/cs-zone-cloner/definition"
	"os"
	"strconv"
	"strings"
)

type config struct {
	apiKey    string
	apiSecret string

	hostScheme string
	hostAddr   string
	hostPath   string

	zoneID   string
	zoneName string
	allZones bool

	domainID   string
	domainName string
	allDomains bool

	format string
	output string

	dbHost     string
	dbPort     uint
	dbSchema   string
	dbUser     string
	dbPassword string

	fetch    string
	fetchers []definition.Fetcher
}

type Command struct {
	self string
	log  command.Logger
	conf *config
}

func New(self string, log command.Logger) *Command {
	c := &Command{
		self: self,
		log:  log,
		conf: new(config),
	}
	return c
}

func (Command) Synopsis() string {
	return "Back up existing Zone configuration"
}

func (c Command) Help() string {
	return fmt.Sprintf(`Usage: %s backup [options]

    Perform a backup of an existing Zone's configuration

Required:
    -key            API key
    -secret         API secret

Optional:
	-zone-name      Name of specific Zone to back up.  Mutually exclusive with "zone-id" and "all-zones".
    -zone-id        ID of specific Zone to back up.  Mutually exclusive with "zone-name" and "all-zones"
	-all-zones		Back up all Zones.  Mutually exclusive with "zone-id" and "zone-name"

	-domain-name	Name of specific Domain to back up. Mutually exclusive with "domain-id" and "all-domains"
	-domain-id		ID of specific Domain to back up.  Mutually exclusive with "domain-name" and "all-domains"
	-all-domains	Back up all Domains.  Mutually exclusive with "domain-id" and "domain-name"

    -scheme         "http" or "https" (default: %s) 
    -host           Managment Server hostname with port (default: %s)
    -path           Managment Server api path (default: %s)

    -format         Backup format (currently only "json" is supported)
    -output         File to write backup to (default: echo to stdout)

	-db-host        Database host to add to output (default: %s)
    -db-port        Database port to add to output (default: %d)
    -db-schema      Database schema to add to output
    -db-user        Database user to add to output
    -db-password    Database password to add to output
    
	-fetch          Comma-separated list of fetchers to execute (default: %s)

`,
		c.self,
		definition.DefaultScheme,
		definition.DefaultHost,
		definition.DefaultPath,
		definition.DefaultDBHost,
		definition.DefaultDBPort,
		strings.Join(definition.DefaultFetchers(), ","))
}

func (c Command) Run(args []string) int {
	var err error

	if err = c.parseFlags(args); err != nil {
		c.log.Printf("[error] Setup failed: %s", err)
		return 1
	}

	defConf := definition.Config{
		Key:      c.conf.apiKey,
		Secret:   c.conf.apiSecret,
		Scheme:   c.conf.hostScheme,
		Host:     c.conf.hostAddr,
		Path:     c.conf.hostPath,
		ZoneName: c.conf.zoneName,
		ZoneID:   c.conf.zoneID,
		Fetchers: c.conf.fetchers,
	}
	dbConf := &definition.DatabaseConfig{
		Server:   c.conf.dbHost,
		Port:     int(c.conf.dbPort),
		Schema:   c.conf.dbSchema,
		User:     c.conf.dbUser,
		Password: c.conf.dbPassword,
	}

	definition.SetPackageLogger(c.log)

	zd, err := definition.FetchDefinition(defConf, dbConf)
	if err != nil {
		if ml, ok := c.log.(command.MutableLogger); ok {
			ml.UnMute()
		}
		c.log.Printf("[error] Execution failed: %s", err)
		return 1
	}

	c.log.Println("[info] Definition built")

	if c.conf.format == "json" {
		b, err := definition.FormatJSONIndent(zd)
		if err != nil {
			if ml, ok := c.log.(command.MutableLogger); ok {
				ml.UnMute()
			}
			c.log.Printf("[error] Error formatting: %s", err)
			return 1
		}
		if c.conf.output == "" {
			fmt.Println(string(b))
		} else {
			if f, err := os.Create(c.conf.output); err != nil {
				c.log.Printf("[error] Error opening \"%s\": %s", c.conf.output, err)
				return 1
			} else if _, err = f.Write(b); err != nil {
				c.log.Printf("[error] Error writing to \"%s\": %s", c.conf.output, err)
				return 1
			} else {
				c.log.Printf("[info] Definition written to file \"%s\"", c.conf.output)
				f.Close()
			}
		}
	}

	return 0
}

func (c Command) parseFlags(args []string) error {
	var err error

	if c.conf == nil {
		return errors.New("command improperly constructed")
	}

	fs := flag.NewFlagSet("backup", flag.ContinueOnError)

	fs.StringVar(&c.conf.apiKey, "key", "", "API Key")
	fs.StringVar(&c.conf.apiSecret, "secret", "", "API Secret")

	fs.StringVar(&c.conf.zoneID, "zone-id", "", "ID of specific Zone to clone (mutually exclusive with zone-name)")
	fs.StringVar(&c.conf.zoneName, "zone-name", "", "Name of Zone to clone (mutually exclusive with zone-id)")

	fs.StringVar(&c.conf.domainID, "domain-id", "", "ID of specific Domain to clone (mutually exclusive with domain-name")
	fs.StringVar(&c.conf.domainName, "domain-name", "", "Name of specific Domain to clone (mutually exclusive with domain-id")

	fs.StringVar(&c.conf.hostScheme, "scheme", definition.DefaultScheme, "HTTP Scheme to use (http or https)")
	fs.StringVar(&c.conf.hostAddr, "host", definition.DefaultHost, "CloudStack Management host addr including port")
	fs.StringVar(&c.conf.hostPath, "path", definition.DefaultPath, "API path")

	fs.StringVar(&c.conf.format, "format", "json", "Currently supports JSON")
	fs.StringVar(&c.conf.output, "output", "", "File to write to")

	fs.StringVar(&c.conf.dbHost, "db-server", definition.DefaultDBHost, "Database host")
	fs.UintVar(&c.conf.dbPort, "db-port", definition.DefaultDBPort, "Database port")
	fs.StringVar(&c.conf.dbSchema, "db-schema", "", "Database schema")
	fs.StringVar(&c.conf.dbUser, "db-user", "", "Database user")
	fs.StringVar(&c.conf.dbPassword, "db-pass", "", "Database password")

	fs.StringVar(&c.conf.fetch, "fetch", strings.Join(definition.DefaultFetchers(), ","), "Comma-separated list of fetchers to execute")

	if err = fs.Parse(args); err != nil {
		return err
	}

	configOK := true

	if c.conf.apiKey == "" {
		c.log.Println("[error] key cannot be empty")
		configOK = false
	}
	if c.conf.apiSecret == "" {
		c.log.Println("[error] secret cannot be empty")
		configOK = false
	}
	c.conf.hostScheme = strings.ToLower(c.conf.hostScheme)
	if c.conf.hostScheme != "http" && c.conf.hostScheme != "https" {
		c.log.Println("[error] scheme must be \"http\" or \"https\"")
		configOK = false
	}
	if c.conf.hostAddr == "" {
		c.log.Println("[error] host cannot be empty")
		configOK = false
	}
	if c.conf.hostPath == "" {
		c.log.Println("[error] path cannot be empty")
		configOK = false
	}
	if c.conf.zoneName != "" && c.conf.zoneID != "" {
		c.log.Println("[error] zone-id and zone-name cannot be set at once")
		configOK = false
	}
	if c.conf.domainName != "" && c.conf.domainID != "" {
		c.log.Println("[error] domain-id and domain-name cannot be set at once")
	}
	c.conf.format = strings.ToLower(c.conf.format)
	if c.conf.format != "json" {
		c.log.Println("[error] format must be json")
		configOK = false
	}

	var fetchers []string

	if cf := strings.Split(c.conf.fetch, ","); len(cf) > 0 {
		fetchers = cf
	} else {
		fetchers = definition.DefaultFetchers()
	}

	c.conf.fetchers = make([]definition.Fetcher, 0)
	for _, name := range fetchers {
		if fn, ok := definition.GetFetcher(name); !ok {
			configOK = false
			c.log.Printf("[error] no fetcher \"%s\" defined", name)
		} else {
			c.conf.fetchers = append(c.conf.fetchers, fn)
		}
	}

	if !configOK {
		return errors.New("error parsing flags, see log")
	}

	if c.conf.output == "" {
		if ml, ok := c.log.(command.MutableLogger); ok {
			ml.Mute()
		}
	}

	c.log.Println("[info] Using parameters:")
	c.log.Println("[info]   APIKey: " + c.conf.apiKey)
	c.log.Println("[info]   APISecret: " + c.conf.apiSecret)
	c.log.Println("[info]   HostScheme: " + c.conf.hostScheme)
	c.log.Println("[info]   HostAddr: " + c.conf.hostAddr)
	c.log.Println("[info]   HostPath: " + c.conf.hostPath)
	if c.conf.zoneID != "" {
		c.log.Println("[info]   ZoneID: " + c.conf.zoneID)
	} else if c.conf.zoneName != "" {
		c.log.Println("[info]   ZoneName: " + c.conf.zoneName)
	} else {
		c.log.Println("[info] 	All Zones in region")
	}
	if c.conf.domainID != "" {
		c.log.Println("[info] 	DomainID: " + c.conf.domainID)
	} else if c.conf.domainName != "" {
		c.log.Println("[info] 	DomainName: " + c.conf.domainName)
	} else {
		c.log.Println("[info]  ")
	}
	c.log.Println("[info]   Format: " + c.conf.format)
	if c.conf.output != "" {
		c.log.Println("[info]   Output: " + c.conf.output)
	}
	if c.conf.dbHost != "" {
		c.log.Println("[info]   DB Server: " + c.conf.dbHost)
	}
	if c.conf.dbPort != 0 {
		c.log.Println("[info]   DB Port: " + strconv.FormatUint(uint64(c.conf.dbPort), 10))
	}
	if c.conf.dbSchema != "" {
		c.log.Println("[info]   DB Schema: " + c.conf.dbSchema)
	}
	if c.conf.dbUser != "" {
		c.log.Println("[info]   DB User: " + c.conf.dbUser)
	}
	if c.conf.dbPassword != "" {
		c.log.Println("[info]   DB Password: " + c.conf.dbPassword)
	}

	return nil
}
