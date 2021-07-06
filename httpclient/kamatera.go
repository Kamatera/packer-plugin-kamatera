package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

type Kamatera struct {
	apiUrl      string
	apiClientID string
	apiSecret   string
}

func NewKamateraClient(apiUrl, apiClientID, apiSecret string) *Kamatera {
	return &Kamatera{
		apiUrl:      apiUrl,
		apiClientID: apiClientID,
		apiSecret:   apiSecret,
	}
}

func (k *Kamatera) Request(method string, path string, body interface{}) (interface{}, error) {
	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, fmt.Errorf("cannot encode body %+v", err)
		}
	}

	req, _ := http.NewRequest(method, fmt.Sprintf("%s/%s", k.apiUrl, path), buf)
	req.Header.Add("AuthClientId", k.apiClientID)
	req.Header.Add("AuthSecret", k.apiSecret)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := cleanhttp.DefaultClient()
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when doing http request %+v", err)
	}
	defer res.Body.Close()

	var result interface{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("bad status code from Kamatera API: %d", res.StatusCode)
		} else {
			return nil, fmt.Errorf("invalid response from Kamatera API: %+v", result)
		}
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error response from Kamatera API (%d): %+v", res.StatusCode, result)
	}
	return result, nil
}

func (k *Kamatera) WaitCommand(commandID string) (map[string]interface{}, error) {
	startTime := time.Now()
	time.Sleep(2 * time.Second)

	for {
		if startTime.Add(40*time.Minute).Sub(time.Now()) < 0 {
			return nil, errors.New("timeout waiting for Kamatera command to complete")
		}

		time.Sleep(2 * time.Second)

		result, e := k.Request("GET", fmt.Sprintf("service/queue?id=%s", commandID), nil)
		if e != nil {
			return nil, e
		}

		commands := result.([]interface{})
		if len(commands) != 1 {
			return nil, errors.New("invalid response from Kamatera queue API: invalid number of command responses")
		}

		command := commands[0].(map[string]interface{})
		status, hasStatus := command["status"]
		if hasStatus {
			switch status.(string) {
			case "complete":
				return command, nil
			case "error":
				log, hasLog := command["log"]
				if hasLog {
					return nil, fmt.Errorf("kamatera command failed: %s", log)
				} else {
					return nil, fmt.Errorf("kamatera command failed: %v", command)
				}
			}
		}
	}
}
