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
	"github.com/ken5scal/gsuite_toolkit/services/users"
	"github.com/ken5scal/gsuite_toolkit/services/reports"
	"github.com/urfave/cli"
	"sort"
	"github.com/urfave/cli/altsrc"
)

const (
	clientSecretFileName = "client_secret.json"
)

func main() {
	app := cli.NewApp()
	app.Name = "gsuite"
	app.Usage = "help managing gsuite"
	app.Version = "0.1"

	var option string
	flags := []cli.Flag {
		cli.StringFlag{
			Name: "repot option",
			Value: "2sv, 2sv",
			Usage: "Get report about `2SV`",
			Destination: &option,
		},
	}

	app.Action  = func(c *cli.Context) error {
		arg := "repot"
		if c.NArg() >0 {
			arg = c.Args()[0]
		}

		switch arg {
		case "report":
			if option == "2sv" {
				// Get 2 sv report
			}
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name: "get 2sv",
			Category: "report",
		},
		{
			Name: "get login",
			Category: "report",
		},
	}

	app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("flagfilename"))
	app.Flags = flags

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)

	scopes := []string{
		admin.AdminDirectoryOrgunitScope, admin.AdminDirectoryUserScope,
		report.AdminReportsAuditReadonlyScope, report.AdminReportsUsageReadonlyScope,
	}
	c := client.NewClient(clientSecretFileName, scopes)
	goneUsers, err := users.GetUsersWhoHasNotLoggedInFor30Days(c.Client)
	if err != nil {
		log.Fatalln(err)
	}
	for _, user := range goneUsers {
		fmt.Println(user.PrimaryEmail)
	}


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

	non2SVuser, err := s.GetNon2StepVerifiedUsers()
	if err != nil {
		log.Fatalln(err)
	}

	for _, user := range non2SVuser.Users {
		fmt.Println(user.Entity.UserEmail)
	}

	//
	//payload := constructPayload("/non2SVuser/suzuki/Desktop/org_structure.csv")
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
