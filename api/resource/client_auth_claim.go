package resource

type (
	ClientAdminAuthRequestResource struct {
		Claim
	}

	Claim struct {
		ClientID            string   `json:"client_id,omitempty"`
		Authority           string   `json:"authority,omitempty"`
		Scopes              []string `json:"scopes,omitempty"`
		Username            string   `json:"username,omitempty"`
		Password            string   `json:"password,omitempty"`
		RedirectURI         string   `json:"redirect_uri,omitempty"`
		CodeChallenge       string   `json:"code_challenge,omitempty"`
		CodeChallengeMethod string   `json:"code_challenge_method,omitempty"`
		State               string   `json:"state,omitempty"`
		ClientSecret        string   `json:"client_secret,omitempty"`
		Thumbprint          string   `json:"thumbprint,omitempty"`
		PemData             string   `json:"pem_data,omitempty"`
	}

	ClientAdminAuthReplyResource struct {
		Status bool   `json:"status"`
		Token  string `json:"token"`
	}
)
