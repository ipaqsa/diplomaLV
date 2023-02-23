package service

type Backend struct {
	Address string
}

type PoolT struct {
	backends []*Backend
	current  uint64
}
