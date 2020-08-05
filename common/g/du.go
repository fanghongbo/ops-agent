package g

import "sync"

const DuBs = "du.bs"

var (
	duPaths     []string
	duPathsLock = new(sync.RWMutex)
)

func DuPathMeta() []string {
	duPathsLock.RLock()
	defer duPathsLock.RUnlock()
	return duPaths
}

func SetDuPathMeta(paths []string) {
	duPathsLock.Lock()
	defer duPathsLock.Unlock()
	duPaths = paths
}
