package resource

import (
	"admincheckapi/api/model"
)

type (
	ClientGroupAdmin struct {
		Admin bool `json:"admin"`
	}

	ClientGroupAdminReplyResource struct {
		Status bool             `json:"status"`
		Data   ClientGroupAdmin `json:"data"`
	}

	ClientAdminGroups struct {
		Count int64                    `json:"count"`
		Data  []model.ClientAdminGroup `json:"data"`
	}

	ClientAdminGroupReplyResource struct {
		Status bool              `json:"status"`
		Data   ClientAdminGroups `json:"data"`
	}
)
