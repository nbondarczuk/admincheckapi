package resource

type (
	ClientAdminAuthRequestResource struct {
		Data Claim `json:"data"`
	}

	Claim struct {
		ClientID            string   `json:"client_id"`
		Authority           string   `json:"authority"`
		Scopes              []string `json:"scopes"`
		Username            string   `json:"username"`
		Password            string   `json:"password"`
		RedirectURI         string   `json:"redirect_uri"`
		CodeChallenge       string   `json:"code_challenge"`
		CodeChallengeMethod string   `json:"code_challenge_method"`
		State               string   `json:"state"`
		ClientSecret        string   `json:"client_secret"`
		Thumbprint          string   `json:"thumbprint"`
		PemData             string   `json:"pem_data"`
	}

	ClientAdminAuthReplyResource struct {
		Status bool   `json:"status"`
		Token  string `json:"token"`
	}
)
