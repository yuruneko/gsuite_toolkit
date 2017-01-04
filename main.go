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
	"fmt"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFileName, err := tokenCacheFileName()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	token, err := tokenFromFile(cacheFileName)
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken(cacheFileName, token)
	}
	return config.Client(ctx, token)
}
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to : %s\n", file)
	f,err := os.Create(file)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

// getTokenFromWeb uses Config to request a Token.
// Ig returns  the retieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil{
		log.Fatalf("Unable to read authorization code %v", err)
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return token
}

// tokenCacheFiele generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFileName() (string, error) {
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
