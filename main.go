package main

import (
	"log"

	"github.com/ken5scal/gsuite_toolkit/client"
	"github.com/ken5scal/gsuite_toolkit/services/users"
	"google.golang.org/api/admin/directory/v1"

	"fmt"
	"github.com/ken5scal/gsuite_toolkit/services/structure/organizations"
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

	parentName := "CISO室"
	units := []string{"サービスインフラグループ", "社内インフラグループ", "情報セキュリティ管理部"}
	orgUnitService := &organizations.Service{srv.Orgunits}

	if len(units) < 1 {
		return
	}

	_, err = orgUnitService.GetOrganizationUnit(parentName)
	if err != nil {
		log.Fatalf("Could Not find Path: %v, %v", parentName, err)
	}

	for _, unitName := range units {
		_, err = orgUnitService.CreateOrganizationUnit(unitName, "/" + parentName)
		if err != nil {
			log.Fatalln("Failed creating New Org Unit.", err)
		}
	}
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
