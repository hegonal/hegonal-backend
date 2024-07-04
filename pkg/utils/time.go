package utils

import "time"

func TimeNow() time.Time {
	return time.Now().UTC()
}