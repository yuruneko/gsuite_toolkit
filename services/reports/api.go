package reports

import (
	"google.golang.org/api/admin/reports/v1"
	"net/http"
	"time"
	"strings"
	"fmt"
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

// GetUserUsage returns G Suite service activities across your account's users
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

// GetNon2StepVerifiedUsers returns emails of users who have not yet enabled 2 step verification.
// date Must be in ISO 8601 format, yyyy-mm-dd
// Example: GetNon2StepVerifiedUsers("2017-01-01")
func (s *Service) GetNon2StepVerifiedUsers() (*users, error) {
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

	users := &users{len(usageReports.UsageReports), make([]*admin.UsageReport, 0)}

	for _, r := range usageReports.UsageReports {
		if !r.Parameters[0].BoolValue {
			users.InsecureUsers = append(users.InsecureUsers, r)
		}
	}
	return users, err
}

// GetLoginActivities reports login activities of all users within organization
func (s *Service) GetLoginActivities() (*admin.Activities, error) {
	time30DaysAgo := time.Now().Add(-time.Duration(30) * time.Hour * 24)
	layout := "2006-01-02T15:04:05.000Z"
	t := time30DaysAgo.Format(layout)
	return s.ActivitiesService.
		List("all", "login").
		EventName("login_success").
		//StartTime("2017-01-28T20:35:28.000Z").
		StartTime(t).
		Do()
}

func (s *Service) GetUsersNotLoggedIn() {

}

type users struct {
	TotalUser     int
	InsecureUsers []*admin.UsageReport
}