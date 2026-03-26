package grove_engine

import (
	"strings"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

const (
	customPrefix = "c:"
)

type GroveMetaData struct {
	Custom []string
}

func NewGroveMetaDataFrom(nodeMeta bullet_interface.NodeMetadata) GroveMetaData {
	var customStrings []string
	for k, _ := range nodeMeta {
		if strings.HasPrefix(k, customPrefix) {
			trimmed := strings.TrimPrefix(k, customPrefix)
			customStrings = append(customStrings, trimmed)
		}
	}
	return GroveMetaData{Custom: customStrings}
}

type GroveAggregateData struct {
	ActiveCount   int
	CompleteCount int
}

func NewGroveMetaData(customFlags []string) GroveMetaData {
	meta := GroveMetaData{
		Custom: customFlags,
	}
	return meta
}

// map[string]interface{}
func (g *GroveMetaData) ToGrove() *bullet_interface.NodeMetadata {
	var meta = make(bullet_interface.NodeMetadata)
	for _, c := range g.Custom {
		if c != "" {
			meta[customFieldToMetaKey(c)] = 1
		}
	}
	return &meta
}

func customFieldToMetaKey(field string) string {
	return customPrefix + field
}
