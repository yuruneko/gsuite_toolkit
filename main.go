package main

import (
	"fmt"
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

	// r, err := srv.Users.List().Customer("my_customer").MaxResults(10). OrderBy("email").Do()
	r, err := srv.Users.List().Customer("my_customer").MaxResults(500).OrderBy("email").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve users in domain.", err)
	}

	if len(r.Users) == 0 {
		fmt.Print("No users found.\n")
	} else {
		for _, user := range r.Users {
			if user.PrimaryEmail == "suzuki.kengo@moneyforward.co.jp" {
				ChangeOrgUnitPath(srv.Users, user, "dep_ciso")
				//fmt.Printf("%s (%s)\n", user.PrimaryEmail, user.Name.FullName)
				//fmt.Println(user.OrgUnitPath)
				//fmt.Println(user.Organizations)
				//user.OrgUnitPath = "/dep_ciso"
				//_, err := srv.Users.Update("suzuki.kengo@moneyforward.co.jp", user).Do()
				//if err != nil {
				//	log.Fatalf("fuga", err)
				//}
			}
		}
	}
}

func ChangeOrgUnitPath(service *admin.UsersService, user *admin.User, unit string) {
	user.OrgUnitPath = "/" + unit
	_, err := service.Update(user.PrimaryEmail, user).Do()
	if err != nil {
		log.Fatalf("fuga", err)
	}
}
