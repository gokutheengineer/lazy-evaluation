package main

import (
	"fmt"
)

type Any interface{}
type EvalFunc func(Any) (Any, Any)

func main() {
	// function to create even numbers
	evenFunc := func(state Any) (Any, Any) {
		os := state.(int)
		ns := os + 2
		return os, ns
	}

	// create lazy evaluator for even numbers starting from 0
	// every time even is called, even func will generate the next even number, lazy generation
	evenInts := BuildLazyEvaluator(evenFunc, 0)

	for i := 0; i < 10; i++ {
		val := evenInts().(int)
		fmt.Println(val)
	}

}

// type generalized lazy value generator builder
func BuildLazyEvaluator(evalFunc EvalFunc, initState Any) func() Any {
	// generated values will be pushed to retValChan
	retValChan := make(chan Any)
	// loopFunc will continuie to run in a gorotuine
	// and will run evalFunc once the value is consumed when consumer calls the retFunc()
	loopFunc := func() {
		var actState Any = initState
		var retVal Any
		for {
			retVal, actState = evalFunc(actState)
			retValChan <- retVal
		}
	}

	// function that will consume from the retValChan
	retFunc := func() Any {
		return <-retValChan
	}

	// start running loop function, which keeps running to lazy generate values
	go loopFunc()
	return retFunc

}
