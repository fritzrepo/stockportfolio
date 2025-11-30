package comdirect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	URL      string `json:"url"`
	OAuthURL string `json:"oauthUrl"`
}

type Communication struct {
	config Config
	creds  Credentials
}

func GetCommunication(pathToConfig string, pathToCreds string) (*Communication, error) {
	var creds Credentials
	err := loadConfig(pathToCreds, &creds)
	if err != nil {
		return nil, err
	}
	var config Config
	err = loadConfig(pathToConfig, &config)
	if err != nil {
		return nil, err
	}
	return &Communication{
		config: config,
		creds:  creds,
	}, nil
}

func (c *Communication) StartSession() (string, error) {
	// Implementation for starting a session with Comdirect
	return c.getSessionStep1()
}

func (c *Communication) EndSession() error {
	// Implementation for ending a session with Comdirect
	return nil
}

func (c *Communication) getSessionStep1() (string, error) {
	// Implementation for session workflow

	ctx := context.Background()
	client, _ := NewClient(15 * time.Second)

	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	client.SetDefaultHeaders(hdrs)

	//step 1: OAuth2 Token anfordern (OAuth2 Resource Owner Password Credentials Flow)
	var result ResponseOAuth2Flow

	data := url.Values{}
	data.Set("client_id", c.creds.ClientID)
	data.Set("client_secret", c.creds.ClientSecret)
	data.Set("grant_type", "password")
	data.Set("username", c.creds.AccountID)
	data.Set("password", c.creds.Pin)

	response, err := client.PostForm(ctx, c.config.OAuthURL, data, &result)
	if err != nil {
		return "", err
	}

	fmt.Println("Access Token:", result.AccessToken)
	fmt.Println("Token Type:", result.TokenType)
	fmt.Println("Expires In:", result.ExpiresIn)
	fmt.Println("Refresh Token:", result.RefreshToken)
	fmt.Println("Scope:", result.Scope)
	fmt.Println("KdNr:", result.KdNr)
	fmt.Println("BpId:", result.BpId)
	fmt.Println("KontaktId:", result.KontaktId)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("failed to get access token, status code: %d", response.StatusCode)
	}

	//step 2 (Session status)

	// session id bauen
	sessionId := uuid.New().String()

	// Request ID bauen aus den letzten 9 Zeichen des aktuellen Timestamps in ms.
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	timestampStr := fmt.Sprintf("%d", timestamp)
	requestId := timestampStr[len(timestampStr)-9:]

	//{"clientRequestId":{"sessionId":"{{session_id}}","requestId":"{{request_id}}"}}
	infoValue := fmt.Sprintf("{\"clientRequestId\":{\"sessionId\":\"%s\",\"requestId\":\"%s\"}}", sessionId, requestId)

	var sessionObj []ResponseSessionObject

	//Headers setzen
	reqHdr := http.Header{}
	reqHdr.Set("Content-Type", "application/json")
	reqHdr.Set("Authorization", "Bearer "+result.AccessToken)
	reqHdr.Set("x-http-request-info", infoValue)

	ctx = WithHeaders(ctx, reqHdr)

	response, err = client.GetJSON(ctx, c.config.URL+"/session/clients/user/v1/sessions", &sessionObj)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("failed to get session id, status code: %d", response.StatusCode)
	}

	//step 3 hier weitermachen
	//Headers setzen
	reqHdr = http.Header{}
	reqHdr.Set("Authorization", "Bearer SOME_PER_REQUEST_TOKEN")
	ctx = WithHeaders(ctx, reqHdr)

	response, err = client.GetJSON(ctx, c.config.URL+"/session/clients/user/v1/sessions/", &sessionObj)
	if err != nil {
		return "", err
	}

	return response.Status, nil
}

func loadConfig(path string, out interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(out)
	if err != nil {
		return err
	}
	return nil
}
