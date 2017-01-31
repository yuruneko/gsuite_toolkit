package reports

import (
	"google.golang.org/api/admin/reports/v1"
	"net/http"
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

//
func (s *Service) GetUserUsage(key, date, params string) (*admin.UsageReports, error) {
	return s.UserUsageReportService.
		Get(key, date).
		Parameters(params).
		Do()
}

// Date Must be in ISO 8601 format, yyyy-mm-dd
func (s *Service) GetNon2StepVerifiedUsers(date string) ([]string, error) {
	var users []string

	usageReports, err := s.GetUserUsage("all", date, "accounts:is_2sv_enrolled")
	if err != nil {
		return users, err
	}

	for _, r := range usageReports.UsageReports {
		if !r.Parameters[0].BoolValue {
			users = append(users, r.Entity.UserEmail)
		}
	}
	return users, nil
}

// Example:
func (s *Service) GetLoginActivities() (*admin.Activities, error) {
	return s.ActivitiesService.List("all", "login").MaxResults(10).Do()
}
