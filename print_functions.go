package main

import (
	"fmt"
)

// Print a pair in user-friendly form
func printPair(pair Pair) {
	fmt.Print(" < ")
	printVector(pair.V1)
	fmt.Print(" , ")
	printVector(pair.V2)
	fmt.Print(" > ")
}

// Print the map in user-friendly form
func printMap(P map[string]map[string]Pair) {
	fmt.Println("{")
	for key,special_vector := range P {
		fmt.Print(" ", key,": ") // Stampa il nome dell'elemento

		for process,pair := range special_vector {
			fmt.Print(process, ": ")
			printPair(pair)
		}

		fmt.Println()
	}
	fmt.Print("}")
}

// Print the vector in user-friendly form
func printVector(vector map[string]SignedElement) {
	fmt.Print("[")
	for _,signedElement := range vector {
		fmt.Print(signedElement.Element,",")
	}
	fmt.Print("]")
}