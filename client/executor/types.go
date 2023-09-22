package executor

type ExecReq struct {
	// ID
	Id uint64 `json:"id"` // 唯一标识
	// Name
	Name string `json:"name" validate:"required"`
	// Content
	Content string `json:"content" validate:"required"`
	WithEnd bool   `json:"withEnd"`
}
