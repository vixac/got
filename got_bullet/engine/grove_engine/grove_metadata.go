package grove_engine

import (
	"strings"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
)

const (
	customPrefix = "c:"
)

type GroveMetaData struct {
	Custom  []string
	Created engine.DateTime
	//Scheduled *engine.DateTime
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

func NewGroveMetaData(customFlags []string, createdDate engine.DateTime, scheduled *engine.DateTime) GroveMetaData {
	meta := GroveMetaData{
		Custom:  customFlags,
		Created: createdDate,
	}
	return meta
}

func (g *GroveMetaData) ToGrove() *bullet_interface.NodeMetadata {
	var meta = make(bullet_interface.NodeMetadata)
	for _, c := range g.Custom {
		if c != "" {
			meta[customFieldToMetaKey(c)] = 1
		}
	}
	meta["created"] = g.Created.Date
	return &meta
}

func customFieldToMetaKey(field string) string {
	return customPrefix + field
}
