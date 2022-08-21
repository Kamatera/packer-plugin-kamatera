package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

type Kamatera struct {
	apiUrl      string
	apiClientID string
	apiSecret   string
	ui          packersdk.Ui
}

func NewKamateraClient(apiUrl, apiClientID, apiSecret string, ui packersdk.Ui) *Kamatera {
	return &Kamatera{
		apiUrl:      apiUrl,
		apiClientID: apiClientID,
		apiSecret:   apiSecret,
		ui:          ui,
	}
}

func (k *Kamatera) Request(method string, path string, body interface{}, silent bool) (interface{}, error) {
	url := fmt.Sprintf("%s/%s", k.apiUrl, path)
	buf := new(bytes.Buffer)
	if body == nil {
		if ! silent {
			k.ui.Say(fmt.Sprintf("Kamatera Request: %s %s", method, url))
		}
	} else {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, fmt.Errorf("cannot encode body %+v", err)
		}
		if ! silent {
			k.ui.Say(fmt.Sprintf("Kamatera Request: %s %s %s", method, url, buf.String()))
		}
	}

	req, _ := http.NewRequest(method, url, buf)
	req.Header.Add("AuthClientId", k.apiClientID)
	req.Header.Add("AuthSecret", k.apiSecret)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := cleanhttp.DefaultClient()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
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

		result, e := k.Request("GET", fmt.Sprintf("service/queue?id=%s", commandID), nil, true)
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
