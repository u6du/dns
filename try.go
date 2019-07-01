package dns

func try(nameserver []string) bool {
	for i := range HostTestTxtLi {
		txt := ResolveTxt(HostTestTxtLi[i], nameserver, func(txt *string) bool {
			return txt != nil
		})
		if txt != nil {
			return true
		}
	}
	return false
}
