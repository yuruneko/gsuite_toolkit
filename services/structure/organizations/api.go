package organizations

import "google.golang.org/api/admin/directory/v1"

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
		completePath = append(completePath, "/"+path)
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

// UpdateOrganizationUnit
// PUT https://www.googleapis.com/admin/directory/v1/customer/my_customer/orgunits/corp/support/sales_support
//{
//  "description": "The BEST sales support team"
//}
func (service *Service) UpdateOrganizationUnit(path []string, NewOrgUnit *admin.OrgUnit) (*admin.OrgUnit, error) {
	return service.Patch("my_customer", path, NewOrgUnit).Do()
}