package handler

type TailLogReq struct {
	Follow  bool   `json:"follow"` // 是否实时跟踪
	LogPath string `json:"log_path" validate:"required"`
	SeekEnd bool   `json:"seek_end"`
}

type StopOnRunningExec struct {
	Name string `json:"name" validate:"required"`
}

type DownloadFileReq struct {
	Filepath string `json:"filepath" validate:"required"`
}

type UpgradeReq struct {
	UpgradeUrl string `json:"upgrade_url" validate:"required"`
}
