package dns

func try(nameserver []string) bool {
	for i := range HostTestTxtLi {
		txt := ResolveTxt(nameserver, HostTestTxtLi[i], func(txt string) bool {
			return true
		})
		if txt != nil {
			return true
		}
	}
	return false
}
