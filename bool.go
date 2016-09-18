package mikkamakka

var (
	False = newVal(Bool, false)
	True  = newVal(Bool, true)
)

var InvalidBoolString = SysStringToError("invalid bool string")

func SysStringToBool(s string) *Val {
	switch s {
	case "true":
		return True
	case "false":
		return False
	default:
		return InvalidBoolString
	}
}

func StringToBool(s *Val) *Val {
	return SysStringToBool(StringToSysString(s))
}

func BoolToString(b *Val) *Val {
	if b == True {
		return SysStringToString("true")
	}

	return SysStringToString("false")
}

func IsBool(a *Val) *Val {
	return is(a, Bool)
}

func Yes(a *Val) *Val {
	if a == False {
		return False
	}

	return True
}

func Not(a *Val) *Val {
	if a == False {
		return True
	}

	return False
}
