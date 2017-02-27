package drives

import (
	"google.golang.org/api/drive/v3"
	"net/http"
	"fmt"
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

// SetClient sets a client
func (s *Service) SetClient(client *http.Client) (error) {
	srv, err := drive.New(client)
	if err != nil {
		return err
	}
	s.FilesService = srv.Files

	s.Client = client
	return nil
}

// GetFiles retrieve all files within the domain

// Refer to the following link for supported mimeType: https://developers.google.com/drive/v3/web/mime-types?authuser=0
// To Get Child: Q('PARENT-ID' in parents)
// https://developers.google.com/drive/v3/reference/files/list?authuser=1
func (s *Service) GetFiles(name, mimeType string) ([]*drive.File, error) {
	call := s.FilesService.
		List().
		//Corpus("domain").
		//Fields("*").
		OrderBy("modifiedTime").
		// 本来は'Googleフォーム'で検索したいが、検索結果が帰ってこない
		Q(fmt.Sprintf("name contains '%v' and mimeType = '%v'", name, mimeType)).
		PageSize(1000)

	var reports []*drive.File
	for {
		r, e := call.Do()
		if e != nil {
			return nil, e
		}
		reports = append(reports, r.Files...)
		if r.NextPageToken == "" {
			return reports, nil
		}
		call.PageToken(r.NextPageToken)
	}
}