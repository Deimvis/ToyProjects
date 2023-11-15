package deck

import "reflect"

func repeat(cards []Card, times int) []Card {
	result := make([]Card, len(cards)*times)
	for i := 0; i < times; i++ {
		copy(result[i*len(cards):], cards)
	}
	return result
}

func genericRepeat(slice interface{}, times int) interface{} {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		panic("Input is not a slice")
	}

	elementType := sliceValue.Type().Elem()
	result := reflect.MakeSlice(reflect.SliceOf(elementType), 0, times*sliceValue.Len())

	for i := 0; i < times; i++ {
		result = reflect.AppendSlice(result, sliceValue)
	}

	return result.Interface()
}
