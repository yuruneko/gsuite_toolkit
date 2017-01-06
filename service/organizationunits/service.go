package organizationunits

import "github.com/ken5scal/gsuite_toolkit/client"

type OrganizationUnits struct {
	*client.Client
}

var initClient func(*client.Client)

const (
	ServiceName ="organizationunits"
)