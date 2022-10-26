package resource

type (
	ClientAdminConfigRequestResource struct {
		Token string `json:"token"`
	}

	ClientAdminConfigReplyResource struct {
		Status bool    `json:"status"`
		Data   []Group `json:"data"`
	}

	Group struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)
