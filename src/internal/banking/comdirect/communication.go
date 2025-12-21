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
	config       Config
	creds        Credentials
	accessToken  string
	refreshToken string
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
	return c.getAccessTokens()
}

func (c *Communication) EndSession() error {
	// For Comdirect is no explicit session termination required
	return nil
}

func (c *Communication) getAccessTokens() (string, error) {

	fmt.Println("Start workflow for getting access token...")

	ctx := context.Background()
	client, _ := NewClient(15 * time.Second)

	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	client.SetDefaultHeaders(hdrs)

	//step 1: OAuth2 Token anfordern (OAuth2 Resource Owner Password Credentials Flow)
	var token ResponseOAuth2Flow

	data := url.Values{}
	data.Set("client_id", c.creds.ClientID)
	data.Set("client_secret", c.creds.ClientSecret)
	data.Set("grant_type", "password")
	data.Set("username", c.creds.AccountID)
	data.Set("password", c.creds.Pin)

	response, err := client.PostForm(ctx, c.config.OAuthURL, data, &token)
	if err != nil {
		return "", err
	}

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

	var sessionObjs []SessionObject

	//Headers setzen
	reqHdr := http.Header{}
	reqHdr.Set("Content-Type", "application/json")
	reqHdr.Set("Authorization", "Bearer "+token.AccessToken)
	reqHdr.Set("x-http-request-info", infoValue)

	ctx = WithHeaders(ctx, reqHdr)

	response, err = client.GetJSON(ctx, c.config.URL+"/session/clients/user/v1/sessions", &sessionObjs)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("failed to get session id, status code: %d", response.StatusCode)
	}

	//step 3 Anlage Validierung einer Session-TAN

	// Um auf ein anderes TAN-Verfahren zu wechseln (Beispiel: P_TAN), wird diese Validation-Schnittstelle
	// erneut aufgerufen. Dabei muss im Header folgende Informationen hinzugefügt werden:
	// x-once-authentication-info: {"typ":"P_TAN"}

	var responseSessionObj SessionObject
	var requestSessionObj SessionObject
	requestSessionObj.Identifier = sessionObjs[0].Identifier
	requestSessionObj.SessionTanActive = true
	requestSessionObj.Activated2FA = true

	//Headers bleiben gleich
	sessionUuid := sessionObjs[0].Identifier
	validateUrl := fmt.Sprintf("%s/session/clients/user/v1/sessions/%s/validate", c.config.URL, sessionUuid)
	response, err = client.PostJSON(ctx, validateUrl, &requestSessionObj, &responseSessionObj)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 201 {
		return "", fmt.Errorf("failed to validate session, status code: %d", response.StatusCode)
	}

	//x-once-authentication-info Header auswerten
	authInfoHeader := response.Header.Get("x-once-authentication-info")
	if authInfoHeader == "" {
		return "", fmt.Errorf("missing x-once-authentication-info header in response")
	}

	var onceAuthInfo OnceAuthenticationInfoHeader
	err = json.Unmarshal([]byte(authInfoHeader), &onceAuthInfo)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal x-once-authentication-info header: %v", err)
	}

	// Prüfen
	if !responseSessionObj.SessionTanActive {
		return "", fmt.Errorf("session tan not active after validation")
	}
	if (onceAuthInfo.Id == "") || (onceAuthInfo.Typ == "") {
		return "", fmt.Errorf("invalid once authentication info received")
	}

	fmt.Println()

	//Warten auf Phototan Push
	fmt.Println("\nPress Enter after approving the PhotoTAN Push on your device...")
	fmt.Scanln()

	// Step 4: Aktivierung einer Session-Tan

	// Bei Phototan Push muss kein
	// x-once-authentication:123456 //die eigentliche TAN
	// Header gesetzt werden.

	// Setzen des once authentication info headers
	// x-once-authentication-info:{"id":"7654321"} //id der Challenge
	reqHdr.Set("x-once-authentication-info", fmt.Sprintf("{\"id\":\"%s\"}", onceAuthInfo.Id))
	ctx = WithHeaders(ctx, reqHdr)
	activateUrl := fmt.Sprintf("%s/session/clients/user/v1/sessions/%s", c.config.URL, sessionUuid)
	response, err = client.PatchJSON(ctx, activateUrl, &requestSessionObj, &responseSessionObj)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("failed to activate session tan, status code: %d", response.StatusCode)
	}

	// Prüfen
	if !responseSessionObj.SessionTanActive {
		return "", fmt.Errorf("session tan not active after validation")
	}

	// Step5: Zugriffsrechte innerhalb der Session erweitern
	hdrs = http.Header{}
	hdrs.Set("Accept", "application/json")
	client.SetDefaultHeaders(hdrs)

	data = url.Values{}
	data.Set("client_id", c.creds.ClientID)
	data.Set("client_secret", c.creds.ClientSecret)
	data.Set("grant_type", "cd_secondary")
	data.Set("token", token.AccessToken)

	response, err = client.PostForm(ctx, c.config.OAuthURL, data, &token)
	if err != nil {
		return "", err
	}

	c.accessToken = token.AccessToken
	c.refreshToken = token.RefreshToken

	fmt.Println("Session-TAN activated successfully.")
	fmt.Println("Scope:", token.Scope)

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
