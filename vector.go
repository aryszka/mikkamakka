package main

import "log"

type vect struct {
	items []*val
}

func vectorFromList(l *val) *val {
	log.Println(l)

	var items []*val
	for {
		if l == vnil {
			break
		}

		items, l = append(items, car(l)), cdr(l)
	}

	return &val{vector, &vect{items}}
}

func isVector(a *val) *val {
	if a.mtype == vector {
		return vtrue
	}

	return vfalse
}

func vectorLength(v *val) *val {
	checkType(v, vector)
	return fromInt(len(v.value.(*vect).items))
}

func vectorRef(v, i *val) *val {
	checkType(v, vector)
	checkType(i, number)
	return v.value.(*vect).items[intVal(i)]
}
