package conf

var (
	LenStackBuf = 4096

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int
)
