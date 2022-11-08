package switchd

import "mercury/switch/conf"

type Metadata struct {
	Version string `json:"version"`
}

func NewMetadata(config conf.Configuration) Metadata {
	return Metadata{
		Version: "v1",
	}
}
