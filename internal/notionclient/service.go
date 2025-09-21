// Package notionclient provides a thin wrapper around the Notion SDK. The
// goal is to keep API calls in one place so the rest of the application can be
// tested without hitting the network.
package notionclient

import (
	"context"

	"github.com/jomei/notionapi"
)

// Service wraps a Notion API client and exposes a small set of convenience
// methods used by the renderer and writer.
type Service struct {
	client *notionapi.Client
}

// New creates a Service initialized with the provided Notion integration token.
func New(token string) *Service {
	return &Service{client: notionapi.NewClient(notionapi.Token(token))}
}

// FetchPages queries the given Notion database and returns the list of pages
// (results) returned by the API.
func (s *Service) FetchPages(databaseID string) ([]notionapi.Page, error) {
	resp, err := s.client.Database.Query(context.Background(), notionapi.DatabaseID(databaseID), &notionapi.DatabaseQueryRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// GetChildren retrieves child blocks for the provided block or page ID.
func (s *Service) GetChildren(id notionapi.BlockID) ([]notionapi.Block, error) {
	resp, err := s.client.Block.GetChildren(context.Background(), id, nil)
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}
