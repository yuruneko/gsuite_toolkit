package gsuite_admin_dir_api

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"os/user"
	"path/filepath"
	"os"
	"net/url"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
}

// tokenCacheFiele generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
	url.QueryEscape("admin-directory_v1-go-quickstart.json")), err
}