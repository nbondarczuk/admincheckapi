package resource

type (
	Stat struct {
		Alloc      uint64 `json:"alloc"`
		TotalAlloc uint64 `json:"totalalloc"`
		Sys        uint64 `json:"sys"`
		NumGC      uint32 `json:"numgc"`
	}

	StatResource struct {
		Status bool `json:"status"`
		Data   Stat `json:"data"`
	}

	Version struct {
		Version string `json:"version"`
	}

	VersionResource struct {
		Status bool    `json:"status"`
		Data   Version `json:"data"`
	}
)
