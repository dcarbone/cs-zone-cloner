package definition

import "github.com/xanzy/go-cloudstack/cloudstack"

type (
	DomainDefinition struct {
		Domain   cloudstack.Domain
		Accounts map[string]cloudstack.Account
		Users    map[string]cloudstack.User

		Custom map[string]interface{} `json:"-"`
	}
)

func NewDomainDefinition(domain cloudstack.Domain) *DomainDefinition {
	dd := &DomainDefinition{
		Domain:   domain,
		Accounts: make(map[string]cloudstack.Account),
		Users:    make(map[string]cloudstack.User),

		Custom: make(map[string]interface{}),
	}
	return dd
}
