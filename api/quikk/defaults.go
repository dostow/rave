package quikk

type Quikk struct {
	ShortCode string `json:"shortcode"`
	Public    string `json:"public"`
	Secret    string `json:"secret"`
	PassKey   string `json:"passkey"`
	URL       string `json:"url"`
	Staging   bool
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var stagingAPIURL = "https://tryapi.quikk.dev/v1/mpesa"
var productionAPIURL = "https://tryapi.quikk.dev/v1/mpesa"

func New(shortCode, public, secret, passKey string, staging bool) *Quikk {
	url := productionAPIURL
	if staging {
		url = stagingAPIURL
	}
	return &Quikk{
		ShortCode: shortCode,
		Public:    public,
		Secret:    secret,
		PassKey:   passKey,
		URL:       url,
		Staging:   staging,
	}
}
