package main

import (
	"log"

	"google.golang.org/api/admin/directory/v1"
	"github.com/ken5scal/gsuite_toolkit/client"
)

const (
	clientSecretFileName = "client_secret.json"
)

func main() {
	scopes := []string{admin.AdminDirectoryUserReadonlyScope, admin.AdminDirectoryUserScope}
	c := client.NewClient(clientSecretFileName, scopes)
	srv, err := admin.New(c.Client)
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}
	user := FindUser(srv.Users, "suzuki.kengo@moneyforward.co.jp")
	ChangeOrgUnitPath(srv.Users, user, "CISO室")
}

func FindUser(service *admin.UsersService, email string)  *admin.User {
	user ,err := service.Get(email).Do()
	if err != nil {
		log.Fatalf("Some error%v\n", err )
	}
	return user
}

func ChangeOrgUnitPath(service *admin.UsersService, user *admin.User, unit string) {
	user.OrgUnitPath = "/" + unit
	_, err := service.Update(user.PrimaryEmail, user).Do()
	if err != nil {
		log.Fatalf("fuga", err)
	}
}
