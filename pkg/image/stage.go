package image

import (
	"fmt"
	"strconv"
	"time"
)

type StageID struct {
	DependenciesDigest string `json:"dependenciesDigest"`
	UniqueID  int64  `json:"uniqueID"`
}

func (id StageID) String() string {
	return fmt.Sprintf("dependenciesDigest:%s uniqueID:%d", id.DependenciesDigest, id.UniqueID)
}

func (id StageID) UniqueIDAsTime() time.Time {
	return time.Unix(id.UniqueID/1000, id.UniqueID%1000)
}

type StageDescription struct {
	StageID *StageID `json:"stageID"`
	Info    *Info    `json:"info"`
}

func ParseUniqueIDAsTimestamp(uniqueID string) (int64, error) {
	if timestamp, err := strconv.ParseInt(uniqueID, 10, 64); err != nil {
		return 0, err
	} else {
		return timestamp, nil
	}
}
