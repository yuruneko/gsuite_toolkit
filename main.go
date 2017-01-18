package main

import (
	"log"

	"github.com/ken5scal/gsuite_toolkit/client"
	"github.com/ken5scal/gsuite_toolkit/services/users"
	"google.golang.org/api/admin/directory/v1"
	"encoding/csv"
	"os"
	"io"
	"strings"
	"net/http"
	"io/ioutil"
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

	userService := &users.Service{srv.Users}
	//r, err := userService.GetEmployees("my_customer", "email", 500)

	user, err := userService.GetUser("suzuki.kengo@moneyforward.co.jp")
	_, err = userService.ChangeOrgUnit(user, "社員・委託社員・派遣社員・アルバイト")
	if err != nil {
		log.Fatalln("Failed Changing user's Organizaion unit.", err)
	}

	//orgUnitService := &organizations.Service{srv.Orgunits}
	//_, err = orgUnitService.CreateOrganizationUnits("CISO室", []string{"セキュリティ推進グループ", "サービスインフラグループ", "社内インフラグループ", "情報セキュリティ管理部"})
	//_, err = orgUnitService.GetOrganizationUnit("CISO室/セキュリティ推進グループ")
	//_, err = orgUnitService.UpdateOrganizationUnit(r, "CISO室")

	payload := constructPayload("/Users/suzuki/Desktop/org_structure.csv")
	fmt.Println(payload)
	url := "https://www.googleapis.com/batch"
	//
	//payload := strings.NewReader("--batch_0123456789\nContent-Type: application/http\nContent-ID: <item1:suzuki.kengo@moneyforward.co.jp>\n\nGET https://www.googleapis.com/admin/directory/v1/users/suzuki.kengo@moneyforward.co.jp\n\n--batch_0123456789\nContent-Type: application/http\nContent-ID: <item2:suzuki.kengo@moneyforward.co.jp>\n\nGET https://www.googleapis.com/admin/directory/v1/users/ichikawa.takashi@moneyforward.co.jp\n\n--batch_0123456789--")
	//
	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
	//
	req.Header.Add("content-type", "multipart/mixed; boundary=batch_0123456789")
	req.Header.Add("authorization", "Bearer someToken")
	res, _ := c.Do(req)
	//
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

func constructPayload(filePath string) string {
	var reader *csv.Reader
	var row []string
	var payload string
	boundary := "batch_0123456789"
	header := "--" + boundary + "\nContent-Type: application/http\n\n"

	csv_file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer csv_file.Close()
	reader = csv.NewReader(csv_file)

	for {
		row, err = reader.Read()
		if err == io.EOF {
			return payload + "--batch_0123456789--"
		}

		if strings.Contains(row[5], "@moneyforward.co.jp") && !strings.Contains(payload, row[5]) {
			payload = payload + header + RequestLine("PUT", row[5]) + "\n\n"
		}
	}
}

func RequestLine(method string, email string) string {
	//return "GET https://www.googleapis.com/admin/directory/v1/users/" +  email
	return method + "https://www.googleapis.com/admin/directory/v1/users/" +  email
}

func Body() string {

}