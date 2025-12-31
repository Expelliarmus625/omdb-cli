package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Expelliamus625/omdb-cli/internal/config"
	"github.com/Expelliamus625/omdb-cli/internal/logger"
)

type Client struct {
	config *config.Config
}

func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) GetMovie(movieName string) (*Movie, error) {
	// Validate empty movie name
	if movieName == "" {
		logger.Log.Error("Movie name cannot be empty")
		return nil, fmt.Errorf("Movie name cannot be empty")
	}
	// baseUrl := "https://www.omdbapi.com/"
	baseUrl := c.config.APIBaseUrl
	params := url.Values{}
	params.Add("t", movieName)
	params.Add("apikey", c.config.APIKey)
	url := baseUrl + "?" + params.Encode()

	// Call OMDB Api
	timeout, err := time.ParseDuration(c.config.APITimeout)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("Error calling omdb api", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Decode movie from response json
	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, err
	}

	// Parse 'Response' field from response to find out if api returned an error
	response, err := strconv.ParseBool(movie.Response)
	if err != nil {
		return nil, err
	}
	if !response {
		logger.Log.Warn("Movie not found", "movieName", movieName)
		return nil, fmt.Errorf("Movie not found")
	}

	logger.Log.Info("Movie fetched", "Title", movie.Title)
	return &movie, nil
}
