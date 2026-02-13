package comdirect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
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
	mu            sync.RWMutex
	accessToken   string
	refreshToken  string
	ticker        *time.Ticker
	stopChan      chan struct{}
	isTimerActive bool
	client        *Client
	sessionCtx    context.Context
	sessionCancel context.CancelFunc
}

func NewCommunication(pathToConfig string, pathToCreds string) (*Communication, error) {
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

	httpClient, err := NewClient(15 * time.Second)
	if err != nil {
		return nil, err
	}

	hdrs := http.Header{}
	hdrs.Set("Accept", "application/json")
	httpClient.SetDefaultHeaders(hdrs)

	return &Communication{
		config:        config,
		creds:         creds,
		accessToken:   "",
		refreshToken:  "",
		stopChan:      nil,
		isTimerActive: false,
		client:        httpClient,
	}, nil
}

func (c *Communication) StartSession() error {
	// Implementation for starting a session with Comdirect
	if c.sessionCtx != nil {
		return fmt.Errorf("session already started")
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.sessionCtx = ctx
	c.sessionCancel = cancel

	err := c.getAccessTokens()
	if err != nil {
		c.EndSession()
		return err
	}
	return nil
}

func (c *Communication) EndSession() error {
	// For Comdirect is no explicit session termination required
	c.stopTimer()
	c.revokeToken()
	if c.sessionCancel != nil {
		c.sessionCancel()
		c.sessionCancel = nil
		c.sessionCtx = nil
	}
	return nil
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

// IsTimerActive returns whether the timer is currently running
func (c *Communication) IsTimerActive() bool {
	return c.isTimerActive
}

func (c *Communication) getAccessTokens() error {

	fmt.Println("Start workflow for getting access token...")
	var response *http.Response
	var err error

	//step 1: OAuth2 Token anfordern (OAuth2 Resource Owner Password Credentials Flow)
	var token ResponseOAuth2Flow

	data := url.Values{}
	data.Set("client_id", c.creds.ClientID)
	data.Set("client_secret", c.creds.ClientSecret)
	data.Set("grant_type", "password")
	data.Set("username", c.creds.AccountID)
	data.Set("password", c.creds.Pin)

	// Code block, damit "defer cancel()" möglichst nah am letzten Gebrauch des Contexts ist.
	{
		ctx, cancel := context.WithTimeout(c.sessionCtx, 5*time.Second)
		defer cancel()

		response, err = c.client.PostForm(ctx, c.config.OAuthURL, data, &token)
		if err != nil {
			return err
		}
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("failed to get access token, status code: %d", response.StatusCode)
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

	{
		ctx2, cancel := context.WithTimeout(c.sessionCtx, 5*time.Second)
		defer cancel()

		ctx2 = WithHeaders(ctx2, reqHdr)

		response, err = c.client.GetJSON(ctx2, c.config.URL+"/session/clients/user/v1/sessions", &sessionObjs)
		if err != nil {
			return err
		}
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("failed to get session id, status code: %d", response.StatusCode)
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

	sessionUuid := sessionObjs[0].Identifier

	{
		//Headers bleiben gleich
		ctx3, cancel := context.WithTimeout(c.sessionCtx, 5*time.Second)
		defer cancel()

		ctx3 = WithHeaders(ctx3, reqHdr)

		validateUrl := fmt.Sprintf("%s/session/clients/user/v1/sessions/%s/validate", c.config.URL, sessionUuid)
		response, err = c.client.PostJSON(ctx3, validateUrl, &requestSessionObj, &responseSessionObj)
		if err != nil {
			return err
		}
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("failed to validate session, status code: %d", response.StatusCode)
	}

	//x-once-authentication-info Header auswerten
	authInfoHeader := response.Header.Get("x-once-authentication-info")
	if authInfoHeader == "" {
		return fmt.Errorf("missing x-once-authentication-info header in response")
	}

	var onceAuthInfo OnceAuthenticationInfoHeader
	err = json.Unmarshal([]byte(authInfoHeader), &onceAuthInfo)
	if err != nil {
		return fmt.Errorf("failed to unmarshal x-once-authentication-info header: %v", err)
	}

	// Prüfen
	if !responseSessionObj.SessionTanActive {
		return fmt.Errorf("session tan not active after validation")
	}
	if (onceAuthInfo.Id == "") || (onceAuthInfo.Typ == "") {
		return fmt.Errorf("invalid once authentication info received")
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

	{
		ctx4, cancel := context.WithTimeout(c.sessionCtx, 5*time.Second)
		defer cancel()

		ctx4 = WithHeaders(ctx4, reqHdr)

		activateUrl := fmt.Sprintf("%s/session/clients/user/v1/sessions/%s", c.config.URL, sessionUuid)
		response, err = c.client.PatchJSON(ctx4, activateUrl, &requestSessionObj, &responseSessionObj)
		if err != nil {
			return err
		}
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("failed to activate session tan, status code: %d", response.StatusCode)
	}

	// Prüfen
	if !responseSessionObj.SessionTanActive {
		return fmt.Errorf("session tan not active after validation")
	}

	// Step5: Zugriffsrechte innerhalb der Session erweitern

	data = url.Values{}
	data.Set("client_id", c.creds.ClientID)
	data.Set("client_secret", c.creds.ClientSecret)
	data.Set("grant_type", "cd_secondary")
	data.Set("token", token.AccessToken)

	{
		ctx5, cancel := context.WithTimeout(c.sessionCtx, 5*time.Second)
		defer cancel()

		response, err = c.client.PostForm(ctx5, c.config.OAuthURL, data, &token)
		if err != nil {
			return err
		}
	}

	c.mu.Lock()
	c.accessToken = token.AccessToken
	c.refreshToken = token.RefreshToken
	c.mu.Unlock()

	fmt.Println("Session-TAN activated successfully.")
	fmt.Println("Scope:", token.Scope)
	fmt.Printf("Token expires in %d seconds\n", token.ExpiresIn)

	return nil
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
	c.mu.Lock()
	if c.isTimerActive {
		c.mu.Unlock()
		return fmt.Errorf("timer is already active")
	}
	c.isTimerActive = true
	c.stopChan = make(chan struct{})
	c.mu.Unlock()

	c.ticker = time.NewTicker(interval)

	go func() {
		defer func() {
			c.mu.Lock()
			c.isTimerActive = false
			c.mu.Unlock()
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
	c.mu.RLock()
	if !c.isTimerActive {
		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()

	if c.ticker != nil {
		c.ticker.Stop()
	}

	// Signal the goroutine to stop
	close(c.stopChan)

	// Recreate the stop channel for potential future use
	//c.stopChan = make(chan struct{})

	c.mu.Lock()
	c.isTimerActive = false
	c.mu.Unlock()
	fmt.Println("Timer stopped")
}

func (c *Communication) refreshAccessToken() error {

	ctx, cancel := context.WithTimeout(c.sessionCtx, 5*time.Second)
	defer cancel()

	data := url.Values{}
	data.Set("client_id", c.creds.ClientID)
	data.Set("client_secret", c.creds.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.accessRefreshTokenSnapshot())

	var token ResponseOAuth2Flow

	response, err := c.client.PostForm(ctx, c.config.OAuthURL, data, &token)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("failed to refresh access token, status code: %d", response.StatusCode)
	}

	c.mu.Lock()
	c.accessToken = token.AccessToken
	c.refreshToken = token.RefreshToken
	c.mu.Unlock()

	fmt.Printf("Token expires in %d seconds\n", token.ExpiresIn)

	return nil
}

func (c *Communication) revokeToken() error {
	ctx, cancel := context.WithTimeout(c.sessionCtx, 5*time.Second)
	defer cancel()

	reqHdr := http.Header{}
	reqHdr.Set("Content-Type", "application/x-www-form-urlencoded")
	reqHdr.Set("Authorization", "Bearer "+c.accessTokenSnapshot())
	ctx = WithHeaders(ctx, reqHdr)

	response, err := c.client.Delete(ctx, c.config.OAuthURL+"/revoke")
	if err != nil {
		return err
	}

	if response.StatusCode != 204 {
		return fmt.Errorf("failed to revoke token, status code: %d", response.StatusCode)
	}

	fmt.Println("Access token revoked successfully.")
	return nil
}

func (c *Communication) accessTokenSnapshot() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.accessToken
}

func (c *Communication) accessRefreshTokenSnapshot() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.refreshToken
}
