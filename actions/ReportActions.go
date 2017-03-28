package actions

import (
	"fmt"
	"errors"
	"google.golang.org/api/admin/reports/v1"
	"github.com/ken5scal/gsuite_toolkit/services/reports"
	"github.com/ken5scal/gsuite_toolkit/services"
)

type ReportAction struct {
	*reports.Service
}

func NewReportAction(s services.Service) (*ReportAction, error) {
	if _, ok := s.(*reports.Service); !ok {
		return nil, errors.New(fmt.Sprintf("Invalid type: %T", s))
	}
	return &ReportAction{s.(*reports.Service)}, nil
}

func (action ReportAction) GetNon2StepVerifiedUsers() error {
	report, err := action.Get2StepVerifiedStatusReport()
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

func (action ReportAction) GetAllLoginActivities(daysAgo int) ([]*admin.Activity, error) {
	activities, err := action.GetLoginActivities(daysAgo)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

// GetIllegalLoginUsersAndIp
// Main purpose is to detect employees who have not logged in from office for 30days
func  (action ReportAction)  GetIllegalLoginUsersAndIp(activities []*admin.Activity, officeIPs []string) error {
	data := make(map[string]*LoginInformation)
	for _, activity := range activities {
		email := activity.Actor.Email
		ip := activity.IpAddress

		if value, ok := data[email]; ok {
			if !value.OfficeLogin {
				// If an user has logged in from not verified IP so far
				// then check if new IP is the one from office or not.
				value.OfficeLogin = containIP(officeIPs, ip)
			}
			value.LoginIPs = append(value.LoginIPs, ip)
		} else {
			data[email] = &LoginInformation{
				email,
				containIP(officeIPs, ip),
				[]string{ip}}
		}
	}

	for key, value := range data {
		if !value.OfficeLogin {
			fmt.Println(key)
			fmt.Print("     IP: ")
			fmt.Println(value.LoginIPs)
		}
	}
	return nil
}
type LoginInformation struct {
	Email       string
	OfficeLogin bool
	LoginIPs    []string
}

func containIP(ips []string, ip string) bool {
	set := make(map[string]struct{}, len(ips))
	for _, s := range ips {
		set[s] = struct{}{}
	}

	_, ok := set[ip]
	return ok
}