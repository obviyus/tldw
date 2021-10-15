package mutex

import "sync"

var (
	Db = sync.Mutex{}
)
