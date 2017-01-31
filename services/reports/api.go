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

func GetUserUsage(s *Service) (*admin.UsageReports, error) {
	return s.UserUsageReportService.
		Get("all", "2017-01-26").
		Parameters("accounts:is_2sv_enrolled").
		Do()
}

// Example:
func GetLoginActivities(s *Service) (*admin.Activities, error) {
	return s.ActivitiesService.List("all", "login").MaxResults(10).Do()
}

