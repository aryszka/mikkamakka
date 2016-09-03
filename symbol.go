package main

type sym struct {
	val string
}

func sfromString(s string) *val {
	return &val{symbol, &sym{s}}
}

func sstringVal(s *val) string {
	return s.value.(*sym).val
}
