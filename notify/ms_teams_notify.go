package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type MSTeamsNotify struct {
	Username          string `json:"username"`
	ChannelName       string `json:"channelName"` //Not mandatory field
	ChannelWebhookURL string `json:"channelWebhookURL"`
	IconUrl           string `json:"iconUrl"`
}

type msPostMessage struct {
	// Channel  string `json:"channel"`
	// Username string `json:"username"`
	Text string `json:"text,omitempty"`
	// Icon_url string `json:"icon_url"`
}

func (msTeamsNotify MSTeamsNotify) GetClientName() string {
	return "MSTeams"
}

func (msTeamsNotify MSTeamsNotify) Initialize() error {

	if len(strings.TrimSpace(msTeamsNotify.Username)) == 0 {
		return errors.New("MS Teams: Username is a required field")
	}

	if len(strings.TrimSpace(msTeamsNotify.ChannelWebhookURL)) == 0 {
		return errors.New("MS Teams: channelWebhookURL is a required field")
	}

	return nil
}

func (msTeamsNotify MSTeamsNotify) SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) error {

	message := getMessageFromResponseTimeNotification(responseTimeNotification)

	payload, jsonErr := msTeamsNotify.getJsonParamBody(message)

	if jsonErr != nil {
		return jsonErr
	}
	fmt.Printf("payload: |%v|\n", payload)
	getResponse, respErr := http.Post(msTeamsNotify.ChannelWebhookURL, "application/json", payload)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New("MS Teams : Send notifaction failed. Response code " + strconv.Itoa(getResponse.StatusCode))
	}

	return nil
}

func (msTeamsNotify MSTeamsNotify) SendErrorNotification(errorNotification ErrorNotification) error {

	message := getMessageFromErrorNotification(errorNotification)

	payload, jsonErr := msTeamsNotify.getJsonParamBody(message)

	if jsonErr != nil {
		return jsonErr
	}

	getResponse, respErr := http.Post(msTeamsNotify.ChannelWebhookURL, "application/json", payload)

	if respErr != nil {
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		return errors.New("MS Teams : Send notifaction failed. Response code " + strconv.Itoa(getResponse.StatusCode))
	}

	return nil
}

func (msTeamsNotify MSTeamsNotify) getJsonParamBody(message string) (io.Reader, error) {

	data, jsonErr := json.Marshal(msPostMessage{message})
	if jsonErr != nil {
		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())
		return nil, jsonErr
	}

	return bytes.NewBuffer(data), nil
}
