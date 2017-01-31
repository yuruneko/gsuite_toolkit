package main

import (
	"log"

	"github.com/ken5scal/gsuite_toolkit/client"

	"encoding/csv"
	"fmt"
	admin "google.golang.org/api/admin/directory/v1"
	report "google.golang.org/api/admin/reports/v1"
	"io"
	"os"
	"strings"
	"github.com/ken5scal/gsuite_toolkit/services/reports"
)

const (
	clientSecretFileName = "client_secret.json"
)

func main() {
	scopes := []string{admin.AdminDirectoryOrgunitScope, admin.AdminDirectoryUserScope, report.AdminReportsAuditReadonlyScope, report.AdminReportsUsageReadonlyScope}
	c := client.NewClient(clientSecretFileName, scopes)

	s, err := reports.NewService(c.Client)
	if err != nil {
		log.Fatalln(err)
	}

	r, err := reports.GetUserUsage(s)
	if err != nil {
		log.Fatalln(err)
	}

	count := 0

	for _, reports := range r.UsageReports {
		//fmt.Println(reports.Entity.UserEmail)
		//
		if reports.Parameters[0].BoolValue {
			fmt.Println(reports.Entity.UserEmail)
			count++
		}
	}
	fmt.Println(count)

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
	//
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
	return method + " " + "https://www.googleapis.com/admin/directory/v1/users/" + email + "\n" +
		"Content-Type: application/json\n\n" + Body()
}

func Body() string {
	return "{\n" + "\"orgUnitPath\": \"/社員・委託社員・派遣社員・アルバイト\"\n" + "}\n"
}
