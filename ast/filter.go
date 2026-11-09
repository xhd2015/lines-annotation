package ast

// FilterLoadInfo keeps filtering before annotation construction, not just
// before serialization. The backing AST is not copied and traversal stop
// semantics are preserved.
func FilterLoadInfo(info LoadInfo, match func(string) bool) LoadInfo {
	if match == nil {
		return info
	}
	return filteredLoadInfo{LoadInfo: info, match: match}
}

type filteredLoadInfo struct {
	LoadInfo
	match func(string) bool
}

func (f filteredLoadInfo) RangeFiles(handle func(File) bool) {
	f.LoadInfo.RangeFiles(func(file File) bool {
		if !f.match(file.RelPath()) {
			return true
		}
		return handle(file)
	})
}
