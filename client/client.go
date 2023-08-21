package client

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"riot-data-ingest/api"
	"os"
	"io"
	"time"
	"strconv"
	"fmt"
)

const (
	apiURLFormat      = "%s://%s.%s%s"
	baseURL           = "api.riotgames.com"
	scheme            = "https"
	apiTokenHeaderKey = "X-Riot-Token"
)

type Client struct{
	L      log.FieldLogger
	Region api.Region
	APIKey string
	Client Doer

}

type Doer interface {
	// Do processes an HTTP request and returns the response
	Do(r *http.Request) (*http.Response, error)
}

func NewClient(region api.Region, client Doer, logger log.FieldLogger) *Client {
	return &Client{
		L:      logger,
		Region: region,
		APIKey: os.Getenv("RIOT_TOKEN"),
		Client: client,
	}
}

// DoRequest processes a http.Request and returns the response.
// Rate-Limiting and retrying is handled via the corresponding response headers.
func (c *Client) DoRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	logger := c.Logger().WithField(
		"endpoint", endpoint,
	)
	request, err := c.NewRequest(method, endpoint, body)
	if err != nil {
		logger.Debug(err)
		return nil, err
	}
	response, err := c.Client.Do(request)
	if err != nil {
		logger.Debug(err)
		return nil, err
	}
	if response.StatusCode == http.StatusServiceUnavailable {
		logger.Info("service unavailable, retrying")
		time.Sleep(time.Second)
		response, err = c.Client.Do(request)
		if err != nil {
			logger.Debug(err)
			return nil, err
		}
	}
	if response.StatusCode == http.StatusTooManyRequests {
		retry := response.Header.Get("Retry-After")
		seconds, err := strconv.Atoi(retry)
		if err != nil {
			logger.Debug(err)
			return nil, err
		}
		logger.Infof("rate limited, waiting %d seconds", seconds)
		time.Sleep(time.Duration(seconds) * time.Second)
		return c.DoRequest(method, endpoint, body)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		logger.Debugf("error response: %v", response.Status)
		err, ok := api.StatusToError[response.StatusCode]
		if !ok {
			err = api.Error{
				Message:    "unknown error reason",
				StatusCode: response.StatusCode,
			}
		}
		return nil, err
	}
	return response, nil
}

// NewRequest returns a new http.Request with necessary headers et.
func (c *Client) NewRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	logger := c.Logger().WithField(
		"endpoint", endpoint,
	)
	request, err := http.NewRequest(method, fmt.Sprintf(apiURLFormat, scheme, c.Region, baseURL, endpoint), body)
	if err != nil {
		logger.Debug(err)
		return nil, err
	}
	request.Header.Add(apiTokenHeaderKey, c.APIKey)
	request.Header.Add("Accept", "application/json")
	return request, nil
}


// Logger returns a logger with client specific fields set.
func (c *Client) Logger() log.FieldLogger {
	return c.L.WithFields(log.Fields{
		"method":   CallerName(0),
		"region": c.Region,
	})
}