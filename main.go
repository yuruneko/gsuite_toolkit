package gsuite_admin_dir_api

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"os/user"
	"path/filepath"
	"os"
	"net/url"
	"github.com/kotakanbe/go-cve-dictionary/log"
	"encoding/json"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	to, err := tokenFromFile(cacheFile)
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

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encounterd.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil,err
	}
	t  := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}
