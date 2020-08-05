package g

import "sync"

const UrlCheckHealth = "url.check.health"

var (
	reportUrls     map[string]string
	reportUrlsLock = new(sync.RWMutex)
)

func ReportUrlMeta() map[string]string {
	reportUrlsLock.RLock()
	defer reportUrlsLock.RUnlock()
	return reportUrls
}

func SetReportUrlMeta(urls map[string]string) {
	reportUrlsLock.RLock()
	defer reportUrlsLock.RUnlock()
	reportUrls = urls
}
