package xurl

type Regname struct {
	// TODO: implement proper parsing: https://datatracker.ietf.org/doc/html/rfc3986#section-3.2.2
	v string
}

func (rn *Regname) String() string {
	if rn == nil {
		return ""
	}
	return rn.v
}

func (rn Regname) empty() bool {
	return rn.v == ""
}
