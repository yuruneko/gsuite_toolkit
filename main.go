package main

import (
	"fmt"
	"log"

	"github.com/ken5scal/gsuite_toolkit/client"
	"github.com/ken5scal/gsuite_toolkit/services/structure/organizations"
	"github.com/ken5scal/gsuite_toolkit/services/users"
	"google.golang.org/api/admin/directory/v1"

)

const (
	clientSecretFileName = "client_secret.json"
)

func main() {
	scopes := []string{admin.AdminDirectoryOrgunitScope, admin.AdminDirectoryUserScope}
	c := client.NewClient(clientSecretFileName, scopes)
	srv, err := admin.New(c.Client)
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}

	// r, err := srv.Users.List().Customer("my_customer").MaxResults(10). OrderBy("email").Do()
	userService := &users.Service{srv.Users}
	//r, err := userService.GetEmployees("my_customer", "email", 500)
	user, err := userService.GetUser("suzuki.kengo@moneyforward.co.jp")
	if err != nil {
		log.Fatalln("Unable to retrieve users in domain.", err)
	}

	fmt.Println(user)

	orgUnitService := &organizations.Service{srv.Orgunits}
	r, err := orgUnitService.GetAllOrganizationUnits()
	//orgUnit, err := orgUnitService.CreateOrganizationUnit("セキュリティ推進グループ", "/moneyforward.co.jp/CISO室")
	if err != nil {
		log.Fatalln("Failed creating New Org Unit.", err)
	}

	if len(r.OrganizationUnits) == 0 {
		fmt.Println("No organization units found")
		return
	}

	for _, ou := range r.OrganizationUnits {
		fmt.Println(ou)
	}
}
