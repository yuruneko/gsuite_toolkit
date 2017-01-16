package organizations

import "google.golang.org/api/admin/directory/v1"

// Service provides Organization Units related functionality
// Details are available in a following link
// https://developers.google.com/admin-sdk/directory/v1/guides/manage-org-units#create_ou
type Service struct {
	*admin.OrgunitsService
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
