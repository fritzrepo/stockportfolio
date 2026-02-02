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
	config        Config
	creds         Credentials
	accessToken   string
	refreshToken  string
	ticker        *time.Ticker
	stopChan      chan struct{}
	isTimerActive bool
	client        *Client
	ctx           context.Context
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

	httpCtx := context.Background()
	httpClient, _ := NewClient(15 * time.Second)

	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	httpClient.SetDefaultHeaders(hdrs)

	return &Communication{
		config:        config,
		creds:         creds,
		accessToken:   "",
		refreshToken:  "",
		stopChan:      make(chan struct{}),
		isTimerActive: false,
		client:        httpClient,
		ctx:           httpCtx,
	}, nil
}

func (c *Communication) StartSession() (string, error) {
	// Implementation for starting a session with Comdirect
	return c.getAccessTokens()
}

func (c *Communication) EndSession() error {
	// For Comdirect is no explicit session termination required
	c.stopTimer()
	return nil
}

func (c *Communication) getAccessTokens() (string, error) {

	fmt.Println("Start workflow for getting access token...")

	//step 1: OAuth2 Token anfordern (OAuth2 Resource Owner Password Credentials Flow)
	var token ResponseOAuth2Flow

	data := url.Values{}
	data.Set("client_id", c.creds.ClientID)
	data.Set("client_secret", c.creds.ClientSecret)
	data.Set("grant_type", "password")
	data.Set("username", c.creds.AccountID)
	data.Set("password", c.creds.Pin)

	response, err := c.client.PostForm(c.ctx, c.config.OAuthURL, data, &token)
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

	c.ctx = WithHeaders(c.ctx, reqHdr)

	response, err = c.client.GetJSON(c.ctx, c.config.URL+"/session/clients/user/v1/sessions", &sessionObjs)
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
	response, err = c.client.PostJSON(c.ctx, validateUrl, &requestSessionObj, &responseSessionObj)
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
	c.ctx = WithHeaders(c.ctx, reqHdr)
	activateUrl := fmt.Sprintf("%s/session/clients/user/v1/sessions/%s", c.config.URL, sessionUuid)
	response, err = c.client.PatchJSON(c.ctx, activateUrl, &requestSessionObj, &responseSessionObj)
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
	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	c.client.SetDefaultHeaders(hdrs)

	data = url.Values{}
	data.Set("client_id", c.creds.ClientID)
	data.Set("client_secret", c.creds.ClientSecret)
	data.Set("grant_type", "cd_secondary")
	data.Set("token", token.AccessToken)

	response, err = c.client.PostForm(c.ctx, c.config.OAuthURL, data, &token)
	if err != nil {
		return "", err
	}

	c.accessToken = token.AccessToken
	c.refreshToken = token.RefreshToken

	fmt.Println("Session-TAN activated successfully.")
	fmt.Println("Scope:", token.Scope)
	fmt.Printf("Token expires in %d seconds\n", token.ExpiresIn)

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

// startTimer starts a periodic timer that calls the specified function at the given interval
func (c *Communication) startTimer(interval time.Duration, callback func()) error {
	if c.isTimerActive {
		return fmt.Errorf("timer is already active")
	}

	c.ticker = time.NewTicker(interval)
	c.isTimerActive = true

	go func() {
		defer func() {
			c.isTimerActive = false
		}()

		for {
			select {
			case <-c.ticker.C:
				// Call the callback function
				callback()
			case <-c.stopChan:
				// Timer was stopped
				return
			}
		}
	}()

	fmt.Printf("Timer started with interval: %v\n", interval)
	return nil
}

// stopTimer stops the periodic timer
func (c *Communication) stopTimer() {
	if !c.isTimerActive {
		return
	}

	if c.ticker != nil {
		c.ticker.Stop()
	}

	// Signal the goroutine to stop
	close(c.stopChan)

	// Recreate the stop channel for potential future use
	c.stopChan = make(chan struct{})

	c.isTimerActive = false
	fmt.Println("Timer stopped")
}

// IsTimerActive returns whether the timer is currently running
func (c *Communication) IsTimerActive() bool {
	return c.isTimerActive
}

// RefreshTokenPeriodically starts a timer that periodically refreshes the access token
func (c *Communication) RefreshTokenPeriodically(interval time.Duration) error {
	return c.startTimer(interval, func() {
		fmt.Println("Refreshing access token...")
		err := c.refreshAccessToken()
		if err != nil {
			c.stopTimer()
			fmt.Printf("Error refreshing token: %v\n", err)
		} else {
			fmt.Println("Access token refreshed successfully")
		}
	})
}

func (c *Communication) refreshAccessToken() error {

	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	c.client.SetDefaultHeaders(hdrs)

	data := url.Values{}
	data.Set("client_id", c.creds.ClientID)
	data.Set("client_secret", c.creds.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.refreshToken)

	var token ResponseOAuth2Flow

	response, err := c.client.PostForm(c.ctx, c.config.OAuthURL, data, &token)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("failed to refresh access token, status code: %d", response.StatusCode)
	}

	c.accessToken = token.AccessToken
	c.refreshToken = token.RefreshToken

	fmt.Printf("Token expires in %d seconds\n", token.ExpiresIn)

	return nil
}
