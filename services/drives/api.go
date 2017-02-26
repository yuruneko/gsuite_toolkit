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
// // https://developers.google.com/drive/v3/reference/files/list?authuser=1
func (s *Service) GetFiles() (*drive.FileList, error) {
	call := s.FilesService.
		List().
		Corpus("domain").
		Fields("*").
		OrderBy("modifiedTime").
		PageSize(1000)
		//Spaces("drive,photos")
	var r *drive.FileList
	var e error
	var i int
	for {
		r, e = call.Do()
		if e != nil {
			break
		} else if r.NextPageToken == "" {
			fmt.Println("Im the last")
			break
		} else {
			if len(r.Files) > 0 {
				for _, f := range r.Files {
					//fmt.Printf("%s (%s)\n", f.Name, f.Id)
					for _, o := range f.Owners {
						if o.EmailAddress == "nakade.takuya@moneyforward.co.jp" {
							fmt.Println(i)
							i++
							fmt.Printf("%s (%s)\n", f.Name, f.Id)
						}
					}
				}
			} else {
				fmt.Println("No files found.")
			}

			//fmt.Printf("File Size: %v\n", len(r.Files))
			//fmt.Printf("Page Token: %v\n", r.NextPageToken)
			call.PageToken(r.NextPageToken)
		}
	}
	return r, e
}