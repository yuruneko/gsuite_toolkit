package main

import (
	"log"

	"github.com/ken5scal/gsuite_toolkit/client"
	"github.com/ken5scal/gsuite_toolkit/services/users"

	"encoding/csv"
	"os"
	"io"
	"strings"
	"fmt"
	admin "google.golang.org/api/admin/directory/v1"
	report  "google.golang.org/api/admin/reports/v1"

)

const (
	clientSecretFileName = "client_secret.json"
)

func main() {
	scopes := []string{admin.AdminDirectoryOrgunitScope, admin.AdminDirectoryUserScope, report.AdminReportsAuditReadonlyScope, report.AdminReportsUsageReadonlyScope}
	c := client.NewClient(clientSecretFileName, scopes)
	srv, err := admin.New(c.Client)
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}

	userService := &users.Service{srv.Users}
	//_, err := userService.GetEmployees("my_customer", "email", 500)
	//_, err = userService.ChangeOrgUnit(user, "社員・委託社員・派遣社員・アルバイト")
	user, err := userService.GetUser("suzuki.kengo@moneyforward.co.jp")
	if err != nil {
		log.Fatalln("Failed Changing user's Organizaion unit.", err)
	}
	fmt.Println(user)

	reportService, err := report.New(c.Client)
	if err != nil {
		log.Fatalf("Unable to retrieve reports Client %v", err)
	}

	r, err := reportService.
	UserUsageReport.
		Get("suzuki.kengo@Moneyforward.co.jp", "2017-01-17").
		Parameters("accounts:is_2sv_enrolled,accounts:is_2sv_enforced").
		Do()
	//r, err := reportService.Activities.List("all", "login").MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve logins to domain.", err)
	}

	for _, param := range r.UsageReports[0].Parameters {
		fmt.Printf("%v: %v\n", param.Name , param.BoolValue)
	}

	//if len(r.Items) == 0 {
	//	fmt.Println("No logins found.")
	//} else {
	//	fmt.Println("Logins:")
	//	for _, a := range r.Items {
	//		t, err := time.Parse(time.RFC3339Nano, a.Id.Time)
	//		if err != nil {
	//			fmt.Println("Unable to parse login time.")
	//			// Set time to zero.
	//			t = time.Time{}
	//		}
	//		fmt.Printf("%s: %s %s\n", t.Format(time.RFC822), a.Actor.Email,
	//			a.Events[0].Name)
	//	}
	//}

	//orgUnitService := &organizations.Service{srv.Orgunits}
	//_, err = orgUnitService.CreateOrganizationUnits("CISO室", []string{"セキュリティ推進グループ", "サービスインフラグループ", "社内インフラグループ", "情報セキュリティ管理部"})
	//_, err = orgUnitService.GetOrganizationUnit("CISO室/セキュリティ推進グループ")
	//_, err = orgUnitService.UpdateOrganizationUnit(r, "CISO室")

	//payload := constructPayload("/Users/suzuki/Desktop/org_structure.csv")
	//fmt.Println(payload)
	//url := "https://www.googleapis.com/batch"
	//
	//req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
	//req.Header.Add("content-type", "multipart/mixed; boundary=batch_0123456789")
	//req.Header.Add("authorization", "Bearer someToken")
	//res, _ := c.Do(req)
	//
	//defer res.Body.Close()
	//_, err = ioutil.ReadAll(res.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}
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
	return method + " " + "https://www.googleapis.com/admin/directory/v1/users/" +  email + "\n" +
		"Content-Type: application/json\n\n" + Body()
}

func Body() string {
	return "{\n" + "\"orgUnitPath\": \"/社員・委託社員・派遣社員・アルバイト\"\n" + "}\n"
}
