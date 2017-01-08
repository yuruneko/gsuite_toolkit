package organizationunits

import "google.golang.org/api/admin/directory/v1"

type UsersService struct {
	service *admin.UsersService
}

func (s *UsersService) GetUsers(customer, key string, maxResults int64) (*admin.Users, error) {
	return  s.service.List().
		Customer(customer).
		MaxResults(maxResults).
		OrderBy(key).
		Do()
}

func (s *UsersService) ChangeOrgUnitPath(user *admin.User, unit string) (*admin.User, error) {
	user.OrgUnitPath = "/" + unit
	return s.service.
		Update(user.PrimaryEmail, user).
		Do()
}