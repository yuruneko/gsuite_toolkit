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
	"time"
)

const (
	clientSecretFileName = "client_secret.json"
)

func main() {
	scopes := []string{
		admin.AdminDirectoryOrgunitScope, admin.AdminDirectoryUserScope,
		report.AdminReportsAuditReadonlyScope, report.AdminReportsUsageReadonlyScope,
	}
	c := client.NewClient(clientSecretFileName, scopes)

	s, err := reports.NewService(c.Client)
	if err != nil {
		log.Fatalln(err)
	}

	t := time.Now().Add(-time.Duration(2 * time.Hour * 24))
	ts := strings.Split(t.Format(time.RFC3339), "T") // yyyy-mm-dd
	r, err := s.GetNon2StepVerifiedUsers(ts[0])
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(len(r))

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
