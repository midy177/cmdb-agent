package pkg

// HostInfo 记录主机信息
type HostInfo struct {
	Server string `json:"server,omitempty"`
	Auth   string `json:"auth,omitempty"`
	Id     uint64 `json:"id"`
}
