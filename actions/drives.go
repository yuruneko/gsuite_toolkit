package actions

import (
	"github.com/ken5scal/gsuite_toolkit/services/drives"
	"google.golang.org/api/drive/v3"
	"fmt"
	"strconv"
	"errors"
)

func GetFiles(s *drives.Service) error {
	return nil
}

func GetParents() error {
	return nil
}

func GetParameters(r []*drive.File) error {
	for _, report := range r {
		if report.Capabilities.CanShare {
			continue
		}
		fmt.Println("	" + report.Name + " - " + strconv.FormatBool(report.Capabilities.CanShare))
		fmt.Println("		LastModifier: " + report.LastModifyingUser.EmailAddress)
		if len(report.Permissions) < 0 {
			return errors.New("Supposed to be no permission")
		}

		for _, o := range report.Owners {
			fmt.Println("		" + "I'm owner!" + ": " + o.EmailAddress)
		}

		for _, p := range report.Permissions {
			fmt.Println("		" + p.Role + ": " + p.EmailAddress)
		}
	}
	return nil
}