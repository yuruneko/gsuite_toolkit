package organizations

import "google.golang.org/api/admin/directory/v1"

// Service provides Organization Units related functionality
// Details are available in a following link
// https://developers.google.com/admin-sdk/directory/v1/guides/manage-org-units#create_ou
type Service struct {
	*admin.OrgunitsService
}

