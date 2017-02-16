package reports

import (
	"github.com/ken5scal/gsuite_toolkit/services/reports"
	"fmt"
	"net/http"
	"errors"
)

func GetNon2StepVerifiedUsers(client *http.Client) error {
	s, err := reports.NewService(client)
	if err != nil {
		return err
	}
	report, err := s.Get2StepVerifiedStatusReport()
	if err != nil {
		return err
	}

	if len(report.UsageReports) == 0 {
		return errors.New("No Report Available")
	}

	var paramIndex int
	fmt.Println("Latest Report: " + report.UsageReports[0].Date)
	for i, param := range report.UsageReports[0].Parameters {
		// https://developers.google.com/admin-sdk/reports/v1/guides/manage-usage-users
		// Parameters: https://developers.google.com/admin-sdk/reports/v1/reference/usage-ref-appendix-a/users-accounts
		if param.Name == "accounts:is_2sv_enrolled" {
			paramIndex = i
			break
		}
	}

	for _, r := range report.UsageReports {
		if !r.Parameters[paramIndex].BoolValue {
			fmt.Println(r.Entity.UserEmail)
		}
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