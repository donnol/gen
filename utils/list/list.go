package list

// Filter 从列表中过滤指定元素
func Filter(keys []string, s string) []string {
	newKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		if key == s {
			continue
		}
		newKeys = append(newKeys, key)
	}
	return newKeys
}
