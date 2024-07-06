package utils

import "github.com/lib/pq"

func StringContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func UnorderedEqual(first, second pq.StringArray) bool {
    if len(first) != len(second) {
        return false
    }
    exists := make(map[string]bool)
    for _, value := range first {
        exists[value] = true
    }
    for _, value := range second {
        if !exists[value] {
            return false
        }
    }
    return true
}