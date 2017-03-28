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
	CommandDrive         = "drive"
	FolderMimeType = "application/vnd.google-apps.folder"
	GeneralUsage = "Audit files within Google Drive"
	SubCommandList = "list"
	ListUsage = "list all of files"
	SubCommandSearch = "list"
	SearchUsage = "search a keyword buy specifying an argument"
)

func NewDriveController(s services.Service) (*DriveController, error) {
	if _, ok := s.(*drives.Service); !ok {
		return nil, errors.New(fmt.Sprintf("Invalid type: %T", s))
	}
	return &DriveController{s.(*drives.Service)}, nil
}

func SearchFolders(s services.Service, title string) error {
	dc, err := NewDriveController(s)
	if err != nil {
		return err
	}

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
	dc, err := NewDriveController(s)
	if err != nil {
		return err
	}

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