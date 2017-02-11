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

//type Hoges []*LoginInformation
//func (h Hoges) containHoge(email string) bool {
//	set := make(map[string]struct{}, len(h))
//	for _, s := range h {
//		set[s.Email] = struct{}{}
//	}
//	 _, ok := set[email]
//	return ok
//}
//func (h Hoges) Len() int {
//	return len(h)
//}
//func (h Hoges) Swap(i, j int) {
//	h[i], h[j] = h[j], h[i]
//}
//func (h Hoges) Less(i, j int) bool {
//	return h[i].Email < h[j].Email
//}

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

	loginData, _ := s.GetEmployeesNotLogInFromOfficeIP()

	for key, value := range loginData {
		if !value.OfficeLogin {
			fmt.Println(key)
			fmt.Print("     IP: ")
			fmt.Println(value.LoginIPs)
		}
	}


	//
	//payload := constructPayload("/users/suzuki/Desktop/org_structure.csv")
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
