package services

import (
	"reflect"
)

// Wait struct
// TODO: Assert reflect types are same
type Wait struct {
	ValueType reflect.Type
	Channel   chan interface{}
	Function  func(chan<- interface{})
}

// WaitForAll goroutines
func WaitForAll(waits []Wait) []interface{} {
	length := len(waits)
	values := make([]interface{}, length)
	// var wg sync.WaitGroup
	var selectCases []reflect.SelectCase
	// wg.Add(length)
	for _, wait := range waits {
		selectCases = append(selectCases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(wait.Channel),
		})

		go wait.Function(wait.Channel)
	}

	for i := 0; i < length; i++ {
		chosen, value, ok := reflect.Select(selectCases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			selectCases[chosen].Chan = reflect.ValueOf(nil)
			continue
		}

		values[chosen] = value.Interface()

	}

	return values

}
