package restore

import (
	"fmt"
	"github.com/dcarbone/cs-zone-cloner/command"
	"github.com/dcarbone/cs-zone-cloner/definition"
)

type config struct {
	apiKey    string
	apiSecret string

	hostScheme string
	hostAddr   string
	hostPath   string

	zone     *definition.ZoneDefinition
	zoneID   string
	zoneName string
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
	return "Restore a Zone based on a previous backup"
}

func (c Command) Help() string {
	return fmt.Sprintf(`Usage: %s restore [options]
    Restore all or parts of a Zone based on an existing backup

Required:
    -key            API key
    -secret         API secret

Optional:
    -zone-id        ID of Zone to restore values into if different than in Definition.  Mutually exclusive with "zone-name"
    -zone-name      Name of Zone to restore values into if different than in Definition.  Mutually exclusive with "zone-id"
    -scheme         "http" or "https" (default: %s) 
    -host           Managment Server hostname with port (default: %s)
    -path           Managment Server api path (default: %s)
    
`,
		c.self,
		definition.DefaultScheme,
		definition.DefaultHost,
		definition.DefaultPath)
}

func (c *Command) Run(args []string) int {

	return 0
}
