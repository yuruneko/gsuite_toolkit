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
	"github.com/ken5scal/gsuite_toolkit/services/users"
)

const (
	clientSecretFileName = "client_secret.json"
)

var isChecked = make(map[string]bool)

func main() {
	scopes := []string{
		admin.AdminDirectoryOrgunitScope, admin.AdminDirectoryUserScope,
		report.AdminReportsAuditReadonlyScope, report.AdminReportsUsageReadonlyScope,
	}
	c := client.NewClient(clientSecretFileName, scopes)

	userService, err := users.NewService(c.Client)
	if err != nil {
		log.Fatalln(err)
	}

	users, err := userService.GetAllUsersInDomain("moneyforward.co.jp", 500)
	if err != nil {
		log.Fatalln(err)
	}
	for _, user := range users.Users {
		isChecked[user.PrimaryEmail] = false
	}

	fmt.Println(isChecked)

	s, err := reports.NewService(c.Client)
	if err != nil {
		log.Fatalln(err)
	}
	//
	//r, err := s.GetNon2StepVerifiedUsers()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//fmt.Println("Total User: ", r.TotalUser)
	//fmt.Println("Total Insecure User: ", len(r.InsecureUsers))
	//fmt.Println("Date: ", r.InsecureUsers[0].Date)
	//for _, insecure := range r.InsecureUsers {
	//	fmt.Println(insecure.Entity.UserEmail)
	//}



	a, err := s.GetLoginActivities()
	if err != nil {
		log.Fatalln(err)
	}
	for _, activity := range a.Items {
		fmt.Println(activity.Actor.Email)
		t, err := time.Parse(time.RFC3339Nano, activity.Id.Time)
		if err != nil {
			fmt.Println("Unable to parse login time.")
			// Set time to zero.
			t = time.Time{}
		}
		fmt.Printf("%s: %s %s\n", t.Format(time.RFC822), activity.Actor.Email,
			activity.Events[0].Name)

		for _, event := range activity.Events {
			fmt.Println("	" + event.Name)
			fmt.Println("	" + event.Type)
			fmt.Println(event.Parameters)
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
