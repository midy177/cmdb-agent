package handler

type TailLogReq struct {
	Follow  bool   `json:"follow"` // 是否实时跟踪
	LogPath string `json:"log_path" validate:"required"`
	SeekEnd bool   `json:"seek_end"`
}

type StopOnRunningExec struct {
	Id uint64 `json:"id" validate:"required"`
}

type DownloadFileReq struct {
	Filepath string `json:"filepath" validate:"required"`
}

type UpgradeReq struct {
	UpgradeUrl string `json:"upgrade_url" validate:"required"`
}
