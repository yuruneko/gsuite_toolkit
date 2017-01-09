package users

import (
	"google.golang.org/api/admin/directory/v1"
)

// Service that provides User related administration Task
type Service struct {
	*admin.UsersService
}

// GetEmployees retrieves employees from Gsuite organization.
func (service *Service) GetEmployees(customer, key string, maxResults int64) (*admin.Users, error) {
	return service.List().
		Customer(customer).
		MaxResults(maxResults).
		OrderBy(key).
		Do()
}

// ChangeOrgUnit changes OrgUnit of an user.
func (service *Service) ChangeOrgUnit(user *admin.User, unit string) (*admin.User, error) {
	user.OrgUnitPath = "/" + unit
	return service.Update(user.PrimaryEmail, user).Do()
}