package drives

import (
	"google.golang.org/api/drive/v3"
	"net/http"
)

// Service provides Drive related administration tasks.
// Details are available in a following link.
// https://developers.google.com/drive/v3/web/about-sdk
type Service struct {
	*drive.FilesService
	*http.Client
}

// Initialize Service
func Init() (s *Service) {
	return &Service{}
}

// NewService sets a client
func (s *Service) NewService(client *http.Client) (error) {
	srv, err := drive.New(client)
	if err != nil {
		return nil
	}
	s.FilesService = srv.Files

	s.Client = client
	return nil
}

func (s *Service) GetFiles() (*drive.FileList, error) {
	return s.FilesService.
		List().
		PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
}