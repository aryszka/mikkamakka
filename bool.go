package mikkamakka

var (
	False = &Val{mbool, false}
	True  = &Val{mbool, true}
)

var InvalidBoolString = ErrorFromRawString("invalid bool string")

func BoolFromRawString(s string) *Val {
	switch s {
	case "true":
		return True
	case "false":
		return False
	default:
		return InvalidBoolString
	}
}

func BoolFromString(s *Val) *Val {
	return BoolFromRawString(RawString(s))
}

func BoolToString(b *Val) *Val {
	if b == True {
		return StringFromRaw("true")
	}

	return StringFromRaw("false")
}

func IsBool(a *Val) *Val {
	return is(a, mbool)
}

func Not(a *Val) *Val {
	checkType(a, mbool)

	if a == False {
		return True
	}

	return False
}
