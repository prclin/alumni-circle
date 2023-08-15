package request

import "github.com/prclin/alumni-circle/model/entity"

type BreakAddImageVO struct {
	BreakId     int                        `json:"break_id"`
	BindingList []entity.ImageBreakBinding `json:"binding_list"`
}
