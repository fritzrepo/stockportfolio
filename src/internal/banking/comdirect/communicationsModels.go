package comdirect

type Credentials struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	AccountID    string `json:"AccountID"` //Zugangsnummer
	Pin          string `json:"pin"`
}

type ResponseOAuth2Flow struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	KdNr         string `json:"kdnr"`
	BpId         int    `json:"bpid"`
	KontaktId    int    `json:"kontaktId"`
}

type SessionObject struct {
	Identifier       string `json:"identifier"`
	SessionTanActive bool   `json:"sessionTanActive"`
	Activated2FA     bool   `json:"activated2FA"`
}

type OnceAuthenticationInfoHeader struct {
	Id             string   `json:"id"`
	Typ            string   `json:"typ"`
	AvailableTypes []string `json:"availableTypes"`
	Link           struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
		Type   string `json:"type"`
	} `json:"link"`
}
