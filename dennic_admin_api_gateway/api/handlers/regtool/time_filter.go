package v1

func UpdateTimeFilter(up string) string {
	if up != "0001-01-01 00:00:00 +0000 UTC" {
		return up
	}
	return ""
}
