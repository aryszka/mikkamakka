package mikkamakka

type vect struct {
	items []*Val
}

func vectorFromList(l *Val) *Val {
	var items []*Val
	for {
		if l == Nil {
			break
		}

		items, l = append(items, Car(l)), Cdr(l)
	}

	return &Val{vector, &vect{items}}
}

func isVector(a *Val) *Val {
	if a.mtype == vector {
		return True
	}

	return False
}

func vectorLength(v *Val) *Val {
	checkType(v, vector)
	return NumberFromRawInt(len(v.value.(*vect).items))
}

func vectorRef(v, i *Val) *Val {
	checkType(v, vector)
	checkType(i, number)
	return v.value.(*vect).items[RawInt(i)]
}

func VectorFromList(v *Val) *Val {
	return vectorFromList(v)
}
