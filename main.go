package main

import (
	"log"
	"github.com/ken5scal/gsuite_toolkit/client"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"github.com/urfave/cli"
	"sort"
	"github.com/ken5scal/gsuite_toolkit/mapper/reports"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/ken5scal/gsuite_toolkit/services"
	reports2 "github.com/ken5scal/gsuite_toolkit/services/reports"
)

const (
	ClientSecretFileName = "client_secret.json"
	subCommandReport = "report"
)

func main() {
	b, e := ioutil.ReadFile("gsuite_config.yml")
	if e != nil {
		log.Fatal(e)
	}

	conf := struct {
		Trusted struct {
			Ip []string `yaml:",flow"`
		}
		Domain string
	}{}

	err := yaml.Unmarshal(b, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	app := cli.NewApp()
	app.Name = "gsuite"
	app.Usage = "help managing gsuite"
	app.Version = "0.1"
	app.Authors = []cli.Author{{Name: "Kengo Suzuki", Email:"kengoscal@gmai.com"}}
	app.Action  = func(c *cli.Context) error {
		if c.NArg() == 0 {
			cli.ShowAppHelp(c)
		}
		return nil
	}

	gsuiteClient := client.CreateConfig().
		SetClientSecretFilename(ClientSecretFileName).
		SetScopes([]string{
			client.AdminReportsUsageReadonlyScope.String(),
			client.AdminReportsAuditReadonlyScope.String(), }).
		Build()

	var s services.Service

	app.Commands = []cli.Command{
		{
			Name: subCommandReport,
			Category: subCommandReport,
			Usage: "get report on your G Suite",
			UsageText: "Check potential security risks, trace user login activities, audit admin activities, etc",
			Before: func(*cli.Context) error {
				service := &reports2.Service{}
				err := service.NewService(gsuiteClient)
				return err
			},
			Subcommands: []cli.Command{
				{
					Name:  "2sv",
					Usage: "get employees who have not enabled 2sv",
					Action: func(context *cli.Context) error {
						r, err := s.(*reports2.Service).Get2StepVerifiedStatusReport()
						handleError(err)
						return reports.GetNon2StepVerifiedUsers(r)
					},
				},
				{
					Name:  "illegal_login",
					Usage: "get employees who have not been office for 30 days, but accessing",
					Action: func(c *cli.Context) error {
						r, err := s.(*reports2.Service).GetLoginActivities(30)
						handleError(err)
						return reports.GetIllegalLoginUsersAndIp(r, conf.Trusted.Ip)
					},
				},
			},
		},
		{
			Name: "login",
			Category: subCommandReport,
		},
	}

	//app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("flagfilename"))
	//app.Flags = flags

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)

	//scopes := []string{
	//	admin.AdminDirectoryOrgUnitScope, admin.AdminDirectoryUserScope,
	//	report.AdminReportsAuditReadonlyScope, report.AdminReportsUsageReadonlyScope,
	//}
	//c := gsuite_Client.NewClient(clientSecretFileName, scopes)
	//goneUsers, err := users.GetUsersWhoHasNotLoggedInFor30Days(c.Client)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//for _, user := range goneUsers {
	//	fmt.Println(user.PrimaryEmail)
	//}

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

		if strings.Contains(row[5], "@") && !strings.Contains(payload, row[5]) {
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

func handleError(err error) {
	if err != nil {
		cli.NewExitError(err, 1)
	}
}