package users

import (
	"google.golang.org/api/admin/directory/v1"
)

// Service that provides User related administration Task
type Service struct {
	*admin.UsersService
}

// GetEmployees retrieves employees from Gsuite organization.
func (service *Service) GetEmployees(customer, key string, max int64) (*admin.Users, error) {
	return service.List().
		Customer(customer).
		OrderBy(key).
		MaxResults(max).
		Do()
}

// GetAllUsersInDomain retrieves all users in domain.
func (service *Service) GetAllUsersInDomain(domain, key string, max int64) (*admin.Users, error) {
	return service.List().
		Domain(domain).
		OrderBy(key).
		MaxResults(max).
		Do()
}

// ChangeOrgUnit changes OrgUnit of an user.
func (service *Service) ChangeOrgUnit(user *admin.User, unit string) (*admin.User, error) {
	user.OrgUnitPath = "/" + unit
	return service.Update(user.PrimaryEmail, user).Do()
}