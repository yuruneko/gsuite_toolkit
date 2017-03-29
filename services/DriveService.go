package services

import (
	"google.golang.org/api/drive/v3"
	"net/http"
	"fmt"
)

// DriveService provides Drive related administration tasks.
// Details are available in a following link.
// https://developers.google.com/drive/v3/web/about-sdk
type DriveService struct {
	*drive.FilesService
	*http.Client
	Files []*drive.File
	Call  *drive.FilesListCall
}

// Initialize DriveService
func Init() (s *DriveService) {
	return &DriveService{}
}

// SetClient sets a client
func (s *DriveService) SetClient(client *http.Client) (error) {
	srv, err := drive.New(client)
	if err != nil {
		return err
	}
	s.FilesService = srv.Files

	s.Client = client
	return nil
}

// GetDriveMaterialsWithTitle retrieve all files of which contains {name} in name of file
// Note returning files are ones of which authorized user can see
// Refer to the following link for supported mimeType: https://developers.google.com/drive/v3/web/mime-types?authuser=0
// https://developers.google.com/drive/v3/reference/files/list?authuser=1
func (s *DriveService) GetDriveMaterialsWithTitle(title, mimeType string) ([]*drive.File, error) {
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
			reports = append(reports, f)
		}
		if r.NextPageToken == "" {
			return reports, nil
		}
		call.PageToken(r.NextPageToken)
	}
}

// GetFilesWithinDir searches files within a directory by regular expression
func (s *DriveService) GetFilesWithinDir(parentsId string) ([]*drive.File, error) {
	s.Call = s.FilesService.
		List().
		OrderBy("modifiedTime").
		Fields("*").
		// Refer formats fof Drive query from following link.
		// https://developers.google.com/drive/v3/web/search-parameters
		Q(fmt.Sprintf("'%v' in parents", parentsId))

	if e := s.RepeatCallerUntilNoPageToken(); e != nil {
		return nil, e
	}
	return s.Files, nil
}

func (s *DriveService) GetParents(parentsId string) (*drive.File, error) {
	return s.FilesService.Get(parentsId).Fields("name").Do()
}

func (s *DriveService) RepeatCallerUntilNoPageToken() error {
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