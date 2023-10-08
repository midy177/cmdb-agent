package executor

type ExecReq struct {
	// IsCron
	IsCron bool `json:"is_cron"`
	// GroupId
	GroupId uint64 `json:"group_id,optional"` // 唯一标识
	// Uuid
	Uuid string `json:"uuid,optional"`
	// Content
	Content string `json:"content" validate:"required"`
}
