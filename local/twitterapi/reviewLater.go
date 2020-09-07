// TODO: transform intersectsortedusers into a more general function using
// reflections. start of code below.
// // expirement with reflection.
// func playSliceReflection(slice1, slice2 interface{}) interface{} {

// }

// return a third slice that includes only elements that are present in both
// func intersectSorted(slice1, slice2 interface{}, comp func(int, int) (bool, bool)) interface{} {
// 	finalSlice := reflect.New(reflect.TypeOf(slice1))

// 	sliceA := reflect.ValueOf(slice1)
// 	sliceB := reflect.ValueOf(slice2)

// 	for i, j := 0, 0; i < sliceA.Len() && j < sliceB.Len(); j = j + 1 {
// 		k := 0
// 		eq, less := comp(i, j)

// 		if less {
// 			i++
// 		} else if eq {
// 			finalSlice.Index(k).Set(sliceB.Index(i))
// 			i++
// 			j++
// 			k++
// 		} else {
// 			j++
// 		}
// 	}

// 	return finalSlice.Interface()
// }