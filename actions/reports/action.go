package reports

import (
	"github.com/ken5scal/gsuite_toolkit/services/reports"
	"fmt"
	"net/http"
)

func GetNon2StepVerifiedUsers(client *http.Client) error {
	s, err := reports.NewService(client)
	if err != nil {
		return err
	}
	non2svUserReports, err := s.GetNon2StepVerifiedUsers()
	if err != nil {
		return err
	}

	fmt.Println("Latest Report: " + non2svUserReports.TimeStamp.String())
	for _, user := range non2svUserReports.Users {
		fmt.Println(user.Entity.UserEmail)
	}
	return nil
}

func GetIllegalLoginUsersAndIp(client *http.Client) error {
	s, err := reports.NewService(client)
	if err != nil {
		return err
	}
	loginData, err := s.GetEmployeesNotLogInFromOfficeIP()
	if err != nil {
		return err
	}
	for key, value := range loginData {
		if !value.OfficeLogin {
			fmt.Println(key)
			fmt.Print("     IP: ")
			fmt.Println(value.LoginIPs)
		}
	}
	return nil
}