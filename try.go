package dns

func try(nameserver []string) bool {

	txt := ResolveTxt(nameserver, HostTestTxt, func(txt string) bool {
		return true
	})
	return txt != nil
}
