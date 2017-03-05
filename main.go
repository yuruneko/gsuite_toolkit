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
	reportService "github.com/ken5scal/gsuite_toolkit/services/reports"
	userService "github.com/ken5scal/gsuite_toolkit/services/users"
	driveService "github.com/ken5scal/gsuite_toolkit/services/drives"
	"github.com/BurntSushi/toml"
	"github.com/ken5scal/gsuite_toolkit/models"
	"github.com/ken5scal/gsuite_toolkit/services"
	"fmt"
	"errors"
	"net/http"
	"github.com/ken5scal/gsuite_toolkit/actions"
)

const (
	ClientSecretFileName = "client_secret.json"
	CommandReport        = "report"
	CommandLogin         = "login"
	CommandDrive         = "drive"
)

type network struct {
	Name string
	Ip []string
}

func main() {
	var tomlConf models.TomlConfig
	var s services.Service
	var gsuiteClient *http.Client

	_, err := toml.DecodeFile("gsuite_config.toml", &tomlConf)
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

	gsuiteClient, err = client.CreateConfig().
		SetClientSecretFilename(ClientSecretFileName).
		SetScopes(tomlConf.Scopes).
		Build()
	app.Commands = []cli.Command{
		{
			Name:     CommandDrive,
			Category: CommandDrive,
			Usage:    "Audit files within Google Drive.",
			Before: func(*cli.Context) error {
				s = driveService.Init()
				err := s.SetClient(gsuiteClient)
				return err
			},
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "get all Files",
					Action: func(context *cli.Context) error {
						if _, ok := s.(*driveService.Service); !ok {
							return errors.New(fmt.Sprintf("Invalid type: %T", s))
						}
						dc := actions.NewDriveController(s.(*driveService.Service))
						r, err := dc.GetFiles()
						//r, err := s.GetFilesWithTitle("Google", "application/vnd.google-apps.folder")
						if err !=nil {
							return err
						}

						for _, f := range r {
							if len(f.Parents) > 0 {
								hoge, _ := s.GetParents(f.Parents[0])
								fmt.Printf(hoge.Name + " > ")
							}
							fmt.Print(f.Name + "\n")
							actions.GetPermissions(f)

							if r, err = s.GetFilesWithinDir(f.Id); err !=nil {
								return err
							}
							if err = actions.GetParameters(r); err != nil {
								return err
							}
						}

						return nil
					},
				},
			},
		},
		{
			Name:     CommandLogin,
			Category: CommandLogin,
			Usage:    "Create user profiles, manage user information, even add administrators.",
			Before: func(*cli.Context) error {
				if gsuiteClient == nil {
					return errors.New("No Client. Execute `gsuite_tookit setup`")
				}

				s = userService.Init()
				err := s.SetClient(gsuiteClient)
				return err
			},
			Subcommands: []cli.Command{
				{
					Name:  "rare-login",
					Usage: "get employees who have not logged in for a while",
					Action: func(context *cli.Context) error {
						r, err := s.(*userService.Service).GetUsersWithRareLogin(30, tomlConf.Owner.DomainName)
						if err != nil {
							return err
						}
						for _, user := range r {
							fmt.Println(user.PrimaryEmail)
						}
						return nil
					},
				},
			},
		},
		{
			Name:     CommandReport,
			Category: CommandReport,
			Usage:    "Gain insights on content management with Google Drive activity reports. Audit administrator actions. Generate customer and user usage reports.",
			Before: func(*cli.Context) error {
				s = reportService.Init()
				err := s.(*reportService.Service).SetClient(gsuiteClient)
				return err
			},
			Subcommands: []cli.Command{
				{
					Name:  "2sv",
					Usage: "get employees who have not enabled 2sv",
					Action: func(context *cli.Context) error {
						r, err := s.(*reportService.Service).Get2StepVerifiedStatusReport()
						if err != nil {
							return err
						}
						err = actions.GetNon2StepVerifiedUsers(r)
						if err != nil {
							return err
						}
						return nil
					},
				},
				{
					Name:  "untrusted_login",
					Usage: "get employees who have not been office for 30 days, but accessing",
					Action: func(c *cli.Context) error {
						r, err := s.(*reportService.Service).GetLoginActivities(30)
						if err != nil {
							return err
						}
						err = actions.GetIllegalLoginUsersAndIp(r, tomlConf.GetAllIps())
						if err != nil {
							return err
						}
						return nil
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)

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