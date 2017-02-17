package users

import (
	"google.golang.org/api/admin/directory/v1"
	"net/http"
	"time"
)

// Service provides User related administration Task
// Details are available in a following link
// https://developers.google.com/admin-sdk/directory/v1/guides/manage-users
type Service struct {
	*admin.UsersService
	*http.Client
}

// NewService creates instance of User related Services
func NewService(client *http.Client) (*Service, error) {
	srv, err := admin.New(client)
	if err != nil {
		return nil, err
	}
	return &Service{srv.Users, client}, nil
}

// GetEmployees retrieves employees from Gsuite organization.
// By Default customer key should be "my_customer"
// max shoudl be integer lower than 500
func (service *Service) GetEmployees(customer, key string, max int64) (*admin.Users, error) {
	return service.UsersService.
		List().
		Customer(customer).
		OrderBy(key).
		MaxResults(max).
		Do()
}

// GetAllUsersInDomain retrieves all users in domain.
// GET https://www.googleapis.com/admin/directory/v1/users?domain=example.com&maxResults=2
// Example: GetAllUsersInDomain("hoge.co.jp", "[email, familyname, givenname]", 500)
func (service *Service) GetAllUsersInDomain(domain string, max int64) (*admin.Users, error) {
	return service.UsersService.
		List().
		Domain(domain).
		OrderBy("email").
		MaxResults(max).
		Do()
}

// GetUser retrieves a user based on either email or userID
// GET https://www.googleapis.com/admin/directory/v1/users/userKey
// Example: GetUser("abc@abc.co.jp")
func (service *Service) GetUser(key string) (*admin.User, error) {
	return service.UsersService.Get(key).ViewType("domain_public").Do()
}

// ChangeOrgUnit changes user's OrgUnit.
// PUT https://www.googleapis.com/admin/directory/v1/users/{email/userID}
// Example: ChangeOrgUnit(user, "社員・委託社員・派遣社員・アルバイト")
func (service *Service) ChangeOrgUnit(user *admin.User, unit string) (*admin.User, error) {
	user.OrgUnitPath = "/" + unit
	return service.UsersService.Update(user.PrimaryEmail, user).Do()
}

func GetUsersWhoHasNotLoggedInFor30Days(c *http.Client, domain string) ([]*admin.User, error) {
	u, err := NewService(c)
	if err != nil {
		return nil, err
	}

	users, err := u.GetAllUsersInDomain(domain, 500)
	if err != nil {
		return nil, err
	}

	time30DaysAgo := time.Now().Add(-time.Duration(30) * time.Hour * 24)

	var goneUsers []*admin.User
	for _, user := range users.Users {
		lastLogin, err := time.Parse("2006-01-02T15:04:05.000Z", user.LastLoginTime)
		if err != nil {
			return nil, err
		}
		if time30DaysAgo.After(lastLogin) {
			goneUsers = append(goneUsers, user)
		}
	}

	return goneUsers, nil
}