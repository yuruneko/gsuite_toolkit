package services

import (
	"google.golang.org/api/admin/reports/v1"
	"net/http"
	"time"
	"strings"
	"google.golang.org/api/googleapi"
)

// ReportService provides following functions.
// Content management with Google Drive activity reports.
// Audit administrator actions.
// Generate customer and user usage reports.
// Details are available in a following link
// https://developers.google.com/admin-sdk/reports/
type ReportService struct {
	*admin.UserUsageReportService
	*admin.ActivitiesService
	*admin.ChannelsService
	*admin.CustomerUsageReportsService
	*http.Client
	Call *admin.ActivitiesListCall
	Activities []*admin.Activity
}

// Initialize ReportService
func InitReportService() (s *ReportService) {
	return &ReportService{}
}

// SetClient creates instance of Report related Services
func (s *ReportService) SetClient(client *http.Client) (error) {
	srv, err := admin.New(client)
	if err != nil {
		return err
	}

	s.UserUsageReportService = srv.UserUsageReport
	s.ActivitiesService = srv.Activities
	s.ChannelsService = srv.Channels
	s.CustomerUsageReportsService = srv.CustomerUsageReports
	s.Client = client
	return nil
}

// GetUserUsage returns G Suite service activities across your account's Users
// key should be either "all" or primary id
// params should be one or combination of user report parameters
// https://developers.google.com/admin-sdk/reports/v1/guides/manage-usage-users
// Example:GetUserUsage("all", "2017-01-01", "accounts:is_2sv_enrolled,"accounts:last_name"")
func (s *ReportService) GetUserUsage(key, date, params string) (*admin.UsageReports, error) {
	return s.UserUsageReportService.
		Get(key, date).
		Parameters(params).
		Do()
}

// Get2StepVerifiedStatusReport returns reports about 2 step verification status.
// date Must be in ISO 8601 format, yyyy-mm-dd
// https://developers.google.com/admin-sdk/reports/v1/guides/manage-usage-users
// Example: Get2StepVerifiedStatusReport("2017-01-01")
func (s *ReportService) Get2StepVerifiedStatusReport() (*admin.UsageReports, error) {
	var usageReports *admin.UsageReports
	var err error
	max_retry := 10

	var timeStamp time.Time
	for i := 0; i < max_retry; i++ {
		timeStamp = time.Now().Add(-time.Duration(time.Duration(i) * time.Hour * 24))
		ts := strings.Split(timeStamp.Format(time.RFC3339), "T") // yyyy-mm-dd
		usageReports, err = s.GetUserUsage("all", ts[0], "accounts:is_2sv_enrolled")
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == http.StatusForbidden {
				return nil, err
			}
		} else if err == nil {
			break
		}
	}
	return usageReports, err
}

// GetLoginActivities reports login activities of all Users within organization
// daysAgo: number of past days you are interested from present time
// EX: GetLoginActivities(30)
func (s *ReportService) GetLoginActivities(daysAgo int) ([]*admin.Activity, error) {
	time30DaysAgo := time.Now().Add(-time.Duration(daysAgo) * time.Hour * 24)
	s.Call = s.ActivitiesService.
		List("all", "login").
		EventName("login_success").
		StartTime(time30DaysAgo.Format(time.RFC3339))

	if e := s.RepeatCallerUntilNoPageToken(); e != nil {
		return nil, e
	}
	return s.Activities, nil
}

func (s *ReportService) RepeatCallerUntilNoPageToken() error {
	s.Activities =  []*admin.Activity{}
	for {
		r, e := s.Call.Do()
		if e != nil {
			return e
		}
		s.Activities = append(s.Activities, r.Items...)
		if r.NextPageToken == "" {
			return nil
		}
		s.Call.PageToken(r.NextPageToken)
	}
}