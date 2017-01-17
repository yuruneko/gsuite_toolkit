package main

import (
	"log"

	"github.com/ken5scal/gsuite_toolkit/client"
	"github.com/ken5scal/gsuite_toolkit/services/users"
	"google.golang.org/api/admin/directory/v1"

	"fmt"
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

	fmt.Println(user.PrimaryEmail)

	_, err = userService.ChangeOrgUnit(user, "CISO室/セキュリティ推進グループ")
	if err != nil {
		log.Fatalln("Failed Changing user's Organizaion unit.", err)
	}

	//orgUnitService := &organizations.Service{srv.Orgunits}
	//r, err := orgUnitService.CreateOrganizationUnit("セキュリティ推進グループ", "/dept_ciso")
	//if err != nil {
	//	log.Fatalln("Failed creating New Org Unit.", err)
	//}
	//
	//fmt.Println(r)
	//r, err := orgUnitService.GetOrganizationUnit("dept_ciso/セキュリティ推進グループ")
	//r, err := orgUnitService.GetOrganizationUnit("dept_ciso")
	//if err != nil {
	//	log.Fatalln("Failed creating New Org Unit.", err)
	//}

	//r.Name = "CISO室"
	//r, err = orgUnitService.UpdateOrganizationUnit(r, "dept_ciso")
	//if err != nil {
	//	log.Fatalln("Failed Changing Org Unit.", err)
	//}
}
