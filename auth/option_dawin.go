//go:build !ios && darwin && !appstore

package auth

func isOnAppStore() bool {
	return false
}
