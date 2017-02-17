package reports

import (
	"google.golang.org/api/admin/reports/v1"
	"net/http"
	"time"
	"strings"
	"google.golang.org/api/googleapi"
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
func (s *Service) NewService(client *http.Client) (error) {
	reportService, err := admin.New(client)
	if err != nil {
		return err
	}

	s.UserUsageReportService = reportService.UserUsageReport
	s.ActivitiesService = reportService.Activities
	s.ChannelsService = reportService.Channels
	s.CustomerUsageReportsService = reportService.CustomerUsageReports
	s.Client = client
	return nil
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

// Get2StepVerifiedStatusReport returns reports about 2 step verification status.
// date Must be in ISO 8601 format, yyyy-mm-dd
// https://developers.google.com/admin-sdk/reports/v1/guides/manage-usage-users
// Example: Get2StepVerifiedStatusReport("2017-01-01")
func (s *Service) Get2StepVerifiedStatusReport() (*admin.UsageReports, error) {
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
func (s *Service) GetLoginActivities(daysAgo int) ([]*admin.Activity, error) {
	time30DaysAgo := time.Now().Add(-time.Duration(daysAgo) * time.Hour * 24)
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