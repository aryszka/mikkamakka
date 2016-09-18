package mikkamakka

func SliceToVector(s []*Val) *Val {
	return newVal(Vector, s)
}

func ListToVector(l *Val) *Val {
	var items []*Val
	for {
		if l == NilVal {
			break
		}

		items, l = append(items, Car(l)), Cdr(l)
	}

	return SliceToVector(items)
}

func IsVector(a *Val) *Val {
	return is(a, Vector)
}

func VectorLen(v *Val) *Val {
	checkType(v, Vector)
	return SysIntToNumber(len(v.value.([]*Val)))
}

func VectorRef(v, i *Val) *Val {
	checkType(v, Vector)
	return v.value.([]*Val)[NumberToSysInt(i)]
}
