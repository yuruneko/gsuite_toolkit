package services

import "net/http"

type Service interface {
	NewService(client *http.Client) (error)
}