package buildinfo

import "fmt"

// GetInfo returns a formatted build info.
func GetInfo(version, date, commit string) string {
	info := fmt.Sprintf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s",
		replaceUnknown(version),
		replaceUnknown(date),
		replaceUnknown(commit),
	)
	fmt.Println(info)
	return info
}

func replaceUnknown(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}
