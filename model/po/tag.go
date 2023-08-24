package po

import "time"

type TTag struct {
	Id         uint32    `json:"id"`
	Name       string    `json:"name"`
	State      uint8     `json:"state"`
	Extra      *string   `json:"extra"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
