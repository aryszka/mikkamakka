package mikkamakka

type Vector []*Val

func VectorFromSlice(s Vector) *Val {
	return &Val{vector, s}
}

func VectorFromList(l *Val) *Val {
	var items Vector
	for {
		if l == Nil {
			break
		}

		items, l = append(items, Car(l)), Cdr(l)
	}

	return VectorFromSlice(items)
}

func IsVector(a *Val) *Val {
	return is(a, vector)
}

func VectorLen(v *Val) *Val {
	checkType(v, vector)
	return NumberFromRawInt(len(v.value.(Vector)))
}

func VectorRef(v, i *Val) *Val {
	checkType(v, vector)
	return v.value.(Vector)[RawInt(i)]
}
