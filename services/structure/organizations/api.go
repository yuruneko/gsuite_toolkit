package organizations

import (
	"errors"
	"fmt"
	"google.golang.org/api/admin/directory/v1"
)

// Service provides Organization Units related functionality
// Details are available in a following link
// https://developers.google.com/admin-sdk/directory/v1/guides/manage-org-units#create_ou
type Service struct {
	*admin.OrgunitsService
}

// GetOrganizationUnit retrieves specific organization unit
// EX: GET https://www.googleapis.com/admin/directory/v1/customer/my_customer/orgunits/corp/sales/frontline+sales
func (service *Service) GetOrganizationUnit(paths ...string) (*admin.OrgUnit, error) {
	var completePath []string
	for _, path := range paths {
		completePath = append(completePath, path)
	}
	return service.Get("my_customer", completePath).Do()
}

// GetAllOrganizationUnits fetch all sub-organization units
// GET https://www.googleapis.com/admin/directory/v1/customer/my_customer/orgunits?orgUnitPath=full org unit path&type=all or children
func (service *Service) GetAllOrganizationUnits() (*admin.OrgUnits, error) {
	return service.List("my_customer").Type("all").Do()
}

// CreateOrganizationUnit creates an organization unit
// POST https://www.googleapis.com/admin/directory/v1/customer/my_customer/orgunits
func (service *Service) CreateOrganizationUnit(name, parentOrgUnitPath string) (*admin.OrgUnit, error) {
	newOrgUnit := &admin.OrgUnit{
		Name:              name,
		ParentOrgUnitPath: parentOrgUnitPath,
	}
	return service.Insert("my_customer", newOrgUnit).Do()
}

// CreateOrganizationUnits creates multiple organization units under same parent Org Unit
func (service *Service) CreateOrganizationUnits(names []string, parentOrgUnitPath string) ([]*admin.OrgUnit, error) {
	if len(names) < 1 {
		return nil, errors.New("No Names are defined")
	}

	_, err := service.GetOrganizationUnit(parentOrgUnitPath)
	if err != nil {
		return nil, err
	}

	var createdOrgUnits []*admin.OrgUnit
	e := &OrgUnitCreateError{}

	for _, unitName := range names {
		r, err := service.CreateOrganizationUnit(unitName, "/"+parentOrgUnitPath)
		if err != nil {
			e.ConcatenateMessage(unitName, err)
		} else {
			createdOrgUnits = append(createdOrgUnits, r)
		}
	}

	return createdOrgUnits, e
}

// UpdateOrganizationUnit updates an org unit specified in the path.
// PUT https://www.googleapis.com/admin/directory/v1/customer/my_customer/orgunits/corp/support/sales_support
//{
//  "description": "The BEST sales support team"
//}
func (service *Service) UpdateOrganizationUnit(NewOrgUnit *admin.OrgUnit, paths ...string) (*admin.OrgUnit, error) {
	var path []string
	for _, p := range paths {
		path = append(path, p)
	}
	return service.Patch("my_customer", path, NewOrgUnit).Do()
}

// OrgUnitCreateError implements Error interface and used when creating org unit fails
type OrgUnitCreateError struct {
	messages map[string]string
}

func (err *OrgUnitCreateError) Error() string {
	errorMessage := ""

	for unit, message := range err.messages {
		errorMessage = errorMessage + unit + " -> " + message + "\n"
	}

	return fmt.Sprintf("Failed creating following orgUnit:\n%s", errorMessage)
}

// ConcatenateMessage takes organizationUnit that failed to be created.
func (err *OrgUnitCreateError) ConcatenateMessage(failedOrgUnit string, e error) {
	if err.messages == nil {
		err.messages = make(map[string]string)
	}

	err.messages[failedOrgUnit] = e.Error()
}
