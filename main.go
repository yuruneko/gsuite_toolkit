package main

import (
	"log"

	"github.com/ken5scal/gsuite_toolkit/client"
	"github.com/ken5scal/gsuite_toolkit/services/users"
	"google.golang.org/api/admin/directory/v1"

	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
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

	//_, err = userService.ChangeOrgUnit(user, "CISO室/セキュリティ推進グループ")
	//if err != nil {
	//	log.Fatalln("Failed Changing user's Organizaion unit.", err)
	//}

	//parentName := "CISO室"
	//unitNames := []string{"セキュリティ推進グループ", "サービスインフラグループ", "社内インフラグループ", "情報セキュリティ管理部"}
	//orgUnitService := &organizations.Service{srv.Orgunits}
	//_, err = orgUnitService.CreateOrganizationUnits(unitNames, parentName)
	//if err != nil {
	//	log.Fatalln(err)
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

	url := "https://www.googleapis.com/batch"

	payload := strings.NewReader("--batch_0123456789\nContent-Type: application/http\nContent-ID: <item1:suzuki.kengo@moneyforward.co.jp>\n\nGET https://www.googleapis.com/admin/directory/v1/users/suzuki.kengo@moneyforward.co.jp\n\n--batch_0123456789\nContent-Type: application/http\nContent-ID: <item2:suzuki.kengo@moneyforward.co.jp>\n\nGET https://www.googleapis.com/admin/directory/v1/users/ichikawa.takashi@moneyforward.co.jp\n\n--batch_0123456789--")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "multipart/mixed; boundary=batch_0123456789")
	req.Header.Add("authorization", "Bearer someToken")
	res, _ := c.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
