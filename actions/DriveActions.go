package actions

import (
	"github.com/ken5scal/gsuite_toolkit/services/drives"
	"google.golang.org/api/drive/v3"
	"fmt"
	"strconv"
	"errors"
	"github.com/ken5scal/gsuite_toolkit/services"

)

type DriveController struct {
	*drives.Service
}

const (
	FolderMimeType = "application/vnd.google-apps.folder"
)

func NewDriveController(s *drives.Service) *DriveController {
	return &DriveController{s}
}

func SearchFolders(s services.Service, title string) error {
	if _, ok := s.(*drives.Service); !ok {
		return errors.New(fmt.Sprintf("Invalid type: %T", s))
	}

	dc := NewDriveController(s.(*drives.Service))
	// 本来は'Googleフォーム'で検索したいが、検索結果が帰ってこない
	if r, err := dc.GetDriveMaterialsWithTitle(title, FolderMimeType); err !=nil {
		return  err
	} else {
		for _, f := range r {
			if len(f.Parents) > 0 {
				parent, _ := dc.GetParents(f.Parents[0])
				fmt.Printf(parent.Name + " > ")
			}

			fmt.Print(f.Name + "\n")
			GetPermissions(f)

			if r, err = dc.GetFilesWithinDir(f.Id); err !=nil {
				return  err
			}
			if err = GetParameters(r); err != nil {
				return  err
			}
		}
	}
	return  nil
}

func SearchAllFolders(s services.Service) error {
	if _, ok := s.(*drives.Service); !ok {
		return errors.New(fmt.Sprintf("Invalid type: %T", s))
	}

	dc := NewDriveController(s.(*drives.Service))

	if r, err := dc.GetDriveMaterialsWithTitle("*", FolderMimeType); err !=nil {
		return  err
	} else {
		for _, f := range r {
			if len(f.Parents) > 0 {
				parent, _ := dc.GetParents(f.Parents[0])
				fmt.Printf(parent.Name + " > ")
			}

			fmt.Print(f.Name + "\n")
			GetPermissions(f)

			if r, err = dc.GetFilesWithinDir(f.Id); err !=nil {
				return  err
			}
			if err = GetParameters(r); err != nil {
				return  err
			}
		}
	}
	return  nil
}

func GetPermissions(f *drive.File) {
	for _, p := range f.Permissions {
		fmt.Println("	" + p.Role + ": " + p.EmailAddress)
	}
}

func GetPermissions2(f *drive.File) {
	for _, p := range f.Permissions {
		fmt.Println("		" + p.Role + ": " + p.EmailAddress)
	}
}

func GetParameters(r []*drive.File) error {
	for _, report := range r {
		fmt.Println("	" + report.Name + " - " + strconv.FormatBool(report.Capabilities.CanShare))
		fmt.Println("		LastModifier: " + report.LastModifyingUser.EmailAddress)
		if len(report.Permissions) < 0 {
			return errors.New("Supposed to be no permission")
		}

		for _, o := range report.Owners {
			fmt.Println("		" + "I'm owner!" + ": " + o.EmailAddress)
		}

		GetPermissions2(report)
	}
	return nil
}