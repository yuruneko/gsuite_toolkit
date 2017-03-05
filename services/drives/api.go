package drives

import (
	"google.golang.org/api/drive/v3"
	"net/http"
	"fmt"
	"strings"
)

// Service provides Drive related administration tasks.
// Details are available in a following link.
// https://developers.google.com/drive/v3/web/about-sdk
type Service struct {
	*drive.FilesService
	*http.Client
	Files []*drive.File
	Call  *drive.FilesListCall
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

// GetFilesWithTitle retrieve all files of which contains {name} in name of file
// Note returning files are ones of which authorized user can see
// Refer to the following link for supported mimeType: https://developers.google.com/drive/v3/web/mime-types?authuser=0
// https://developers.google.com/drive/v3/reference/files/list?authuser=1
func (s *Service) GetFilesWithTitle(title, mimeType string) ([]*drive.File, error) {
	call := s.FilesService.
		List().
		//Corpus("domain").
		Fields("*").
		OrderBy("modifiedTime").
		// Refer formats fof Drive query from following link.
		// https://developers.google.com/drive/v3/web/search-parameters
		Q(fmt.Sprintf("name contains '%v' and mimeType = '%v'", title, mimeType))

	var reports []*drive.File
	for {
		r, e := call.Do()
		if e != nil {
			return nil, e
		}

		for _, f := range r.Files {
			if strings.Contains(f.Name,"Googleフォーム") {
				reports = append(reports, f)
			}
		}
		if r.NextPageToken == "" {
			return reports, nil
		}
		call.PageToken(r.NextPageToken)
	}
}

// GetFilesWithinDir searches files within a directory by regular expression
func (s *Service) GetFilesWithinDir(parentsId string) ([]*drive.File, error) {
	s.Call = s.FilesService.
		List().
		OrderBy("modifiedTime").
		Fields("*").
		// Refer formats fof Drive query from following link.
		// https://developers.google.com/drive/v3/web/search-parameters
		Q(fmt.Sprintf("%v' in parents", parentsId))

	if e := s.RepeatCallerUntilNoPageToken(); e != nil {
		return nil, e
	}
	return s.Files, nil
}

func (s *Service) GetParents(parentsId string) (*drive.File, error) {
	return s.FilesService.Get(parentsId).Fields("name").Do()
}

func (s *Service) RepeatCallerUntilNoPageToken() error {
	s.Files = []*drive.File{}
	for {
		r, e := s.Call.Do()
		if e != nil {
			return e
		}
		s.Files = append(s.Files, r.Files...)
		if r.NextPageToken == "" {
			return nil
		}
		s.Call.PageToken(r.NextPageToken)
	}
}