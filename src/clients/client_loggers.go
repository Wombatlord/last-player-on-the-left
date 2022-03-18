package clients

import "log"

const (
	RSSLog  = "RssClient"
	DLLog   = "Downloader"
	TCPALog = "TCPAudio"
)

var loggers = map[string]*log.Logger{
	RSSLog:  nil,
	DLLog:   nil,
	TCPALog: nil,
}

// InitLoggers should be called by the bootstrapping application to
// assign log.Logger instances to the various clients
func InitLoggers(factory func(prefix string) *log.Logger) {
	loggers[RSSLog] = factory(RSSLog)
	loggers[DLLog] = factory(DLLog)
	loggers[TCPALog] = factory(TCPALog)
}
