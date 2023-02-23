package service

var Pool *PoolT

func CreatePool() {
	Pool = &PoolT{
		backends: make([]*Backend, 2),
		current:  0,
	}
	Pool.initPool()
}
