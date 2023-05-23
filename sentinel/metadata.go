package sentinel

import (
	"github.com/arkeonetwork/arkeo/sentinel/conf"
)

var Version = "0.0.0"

type Metadata struct {
	Configuration conf.Configuration `json:"config"`
	Version       string             `json:"version"`
}

func NewMetadata(config conf.Configuration) Metadata {
	return Metadata{
		Version:       Version,
		Configuration: config,
	}
}
