// +build !go1.3 plan9 solaris

package proxy

import "github.com/ayushbpl10/oauth2_proxy/logger"

func WatchForUpdates(filename string, done <-chan bool, action func()) {
	logger.Printf("file watching not implemented on this platform")
	go func() { <-done }()
}
