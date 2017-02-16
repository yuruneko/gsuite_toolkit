package reports

import (
	"google.golang.org/api/admin/reports/v1"
	"net/http"
	"time"
	"strings"
)

// Service provides following functions.
// Content management with Google Drive activity reports.
// Audit administrator actions.
// Generate customer and user usage reports.
// Details are available in a following link
// https://developers.google.com/admin-sdk/reports/
type Service struct {
	*admin.UserUsageReportService
	*admin.ActivitiesService
	*admin.ChannelsService
	*admin.CustomerUsageReportsService
	*http.Client
}

// NewService creates instance of Report related Services
func NewService(client *http.Client) (*Service, error) {
	reportService, err := admin.New(client)
	if err != nil {
		return nil, err
	}

	return &Service{
		reportService.UserUsageReport,
		reportService.Activities,
		reportService.Channels,
		reportService.CustomerUsageReports,
		client}, nil
}

// GetUserUsage returns G Suite service activities across your account's Users
// key should be either "all" or primary id
// params should be one or combination of user report parameters
// https://developers.google.com/admin-sdk/reports/v1/guides/manage-usage-users
// Example:GetUserUsage("all", "2017-01-01", "accounts:is_2sv_enrolled,"accounts:last_name"")
func (s *Service) GetUserUsage(key, date, params string) (*admin.UsageReports, error) {
	return s.UserUsageReportService.
		Get(key, date).
		Parameters(params).
		Do()
}

// GetNon2StepVerifiedUsers returns emails of Users who have not yet enabled 2 step verification.
// date Must be in ISO 8601 format, yyyy-mm-dd
// Example: GetNon2StepVerifiedUsers("2017-01-01")
func (s *Service) GetNon2StepVerifiedUsers() (*Users, error) {
	var usageReports *admin.UsageReports
	var err error
	max_retry := 10

	for i := 0; i < max_retry; i++ {
		t := time.Now().Add(-time.Duration(time.Duration(i) * time.Hour * 24))
		ts := strings.Split(t.Format(time.RFC3339), "T") // yyyy-mm-dd
		usageReports, err = s.GetUserUsage("all", ts[0], "accounts:is_2sv_enrolled")
		if err == nil {
			break
		}
	}

	users := &Users{len(usageReports.UsageReports), make([]*admin.UsageReport, 0)}

	for _, r := range usageReports.UsageReports {
		if !r.Parameters[0].BoolValue {
			users.Users = append(users.Users, r)
		}
	}
	return users, err
}

// GetLoginActivities reports login activities of all Users within organization
func (s *Service) GetLoginActivities() ([]*admin.Activity, error) {
	time30DaysAgo := time.Now().Add(-time.Duration(30) * time.Hour * 24)
	layout := "2006-01-02T15:04:05.000Z"

	call := s.ActivitiesService.
		List("all", "login").
		EventName("login_success").
		StartTime(time30DaysAgo.Format(layout))
		//EndTime("2017-02-05T20:35:28.000Z")

	firstIteration := true
	token := "justrandomtoken"
	var activityList []*admin.Activity
	for token != "" {
		if !firstIteration {
			call.PageToken(token)
		}
		firstIteration = false
		activities, err := call.Do()
		if err != nil {
			return nil, err
		}
		activityList = append(activityList, activities.Items...)
		token = activities.NextPageToken
	}

	return activityList, nil
}

type LoginInformation struct {
	Email       string
	OfficeLogin bool
	LoginIPs    []string
}

// GetEmployeesNotLogInFromOfficeIP
// Main purpose is to detect employees who have not logged in from office for 30days
func (s *Service) GetEmployeesNotLogInFromOfficeIP() (map[string]*LoginInformation, error) {
	data := make(map[string]*LoginInformation)
	officeIPs := []string{"124.32.248.42", "210.130.170.193", "210.138.23.111", "210.224.77.186", "118.243.201.33", "122.220.198.115"}

	activities, err := s.GetLoginActivities()
	if err != nil {
		return nil, err
	}
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

	return data, nil
}

func containIP(ips []string, ip string) bool {
	set := make(map[string]struct{}, len(ips))
	for _, s := range ips {
		set[s] = struct{}{}
	}

	_, ok := set[ip]
	return ok
}

type Users struct {
	TotalUser     int
	Users []*admin.UsageReport
}