package main


func deepCopyVector(original map[string]SignedElement) map[string]SignedElement {
	
	new := make(map[string]SignedElement)

 	for k,v := range original {
		newElement := v

 		new[k] = newElement
	}
	return new
}


/**** COMPONENT-WISE OPERATIONS BETWEEN VECTORS *****/

// Join between vectors
func JoinVector(v1 map[string]SignedElement, v2 map[string]SignedElement, lattice *Lattice) map[string]SignedElement {
	vector_result := make(map[string]SignedElement)
	
	for _, process := range processes {

		final_signature := ""
		if v1[process].Signature == v2[process].Signature {
			final_signature = v1[process].Signature
		} else {
			final_signature = my_address // If a Byzantine sent two justified values, my signature is after the join
		}

		vector_result[process] = SignedElement {
			Element: Join(v1[process].Element,v2[process].Element,lattice).Id,
			Signature: final_signature}

	}
	return vector_result
}

// Equal operator between vectors
func EqualVector(v1 map[string]SignedElement, v2 map[string]SignedElement) bool {
	for _, process := range processes {
		if v1[process].Element != v2[process].Element {
			return false
		}
	}
	return true
}

// Less operator between vectors
func LessVector(v1 map[string]SignedElement, v2 map[string]SignedElement, lattice *Lattice) bool {
	for _, process := range processes {
		if LessEqual(v1[process].Element, v2[process].Element, lattice) == false {
			return false
		}
	}
	for _, process := range processes {
		if Less(v1[process].Element, v2[process].Element, lattice) == false {
			return false
		}
	}
	return true
}