package users

import (
	"google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

type Service struct {
	*admin.UsersService
}

func (service *Service) GetUsers(customer, key string, maxResults int64) (*admin.Users, error) {
	return service.List().
		Customer(customer).
		MaxResults(maxResults).
		OrderBy(key).
		Do()
}

func (service *Service) ChangeOrgUnitPath(user *admin.User, unit string) (*admin.User, error) {
	user.OrgUnitPath = "/" + unit
	return service.Update(user.PrimaryEmail, user).Do()
}