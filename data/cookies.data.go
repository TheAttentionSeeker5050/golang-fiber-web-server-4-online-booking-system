package data

import "time"

// GetCookieExpirationTime returns time in seconds
func GetCookieExpirationTime() time.Duration {
	return time.Duration(60 * 60 * 24 * 7)
}
