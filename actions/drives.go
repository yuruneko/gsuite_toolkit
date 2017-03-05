package actions

import (
	"github.com/ken5scal/gsuite_toolkit/services/drives"
	"google.golang.org/api/drive/v3"
	"fmt"
	"strconv"
	"errors"
)

type DriveController struct {
	*drives.Service
}

func NewDriveController(s *drives.Service) *DriveController {
	return &DriveController{s}
}

func (dc DriveController) GetFiles() ([]*drive.File, error) {
	// 本来は'Googleフォーム'で検索したいが、検索結果が帰ってこない
	title := "Google"
	mimeType := "application/vnd.google-apps.folder"
	if r, err := dc.GetFilesWithTitle(title, mimeType); err !=nil {
		return nil, err
	} else {
		return r, nil
	}
}

func GetParents() error {
	return nil
}

func GetPermissions(f *drive.File) {
	for _, p := range f.Permissions {
		fmt.Println("	" + p.Role + ": " + p.EmailAddress)
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

		GetPermissions(report)
	}
	return nil
}