package main


func VerifyValue(V map[string]SignedElement, j string) bool {
	if V[j].Element == "⊥" || V[j].Signature == j  {
		return true
	} else {
		return false
	}
}

func VerifyPairSet(P map[string]map[string]Pair, j string, key string) bool {
	
	if P[key][j].V1 == nil && P[key][j].V2 == nil || P[key][j].V1[j].Element == P[key][j].V2[j].Element ||  Join(P[key][j].V1[j].Element,P[key][j].V2[j].Element,L).Id != key {
		return false
	}

	if VerifyValue(P[key][j].V1, j) && VerifyValue(P[key][j].V2, j) {
		return true
	}

	if VerifyValue(P[key][j].V1, j) {
		return VerifyPairSet(P, j, P[key][j].V2[j].Element)
	}
	if VerifyValue(P[key][j].V2, j) {
		return VerifyPairSet(P, j, P[key][j].V1[j].Element)
	}

	return VerifyPairSet(P, j, P[key][j].V2[j].Element) && VerifyPairSet(P, j, P[key][j].V1[j].Element)

}

func UnionProofs(V map[string]SignedElement, j string, P map[string]map[string]Pair) {
	lattice_proofs[V[j].Element][j] = P[V[j].Element][j]
	
	if VerifyValue(lattice_proofs[V[j].Element][j].V1, j) == false {
		UnionProofs(lattice_proofs[V[j].Element][j].V1, j, P)
	}
	if VerifyValue(lattice_proofs[V[j].Element][j].V2, j) == false {
		UnionProofs(lattice_proofs[V[j].Element][j].V2, j, P)
	}
}

func AddPair(V map[string]SignedElement, P map[string]map[string]Pair) {
	for j,_ := range localVector {
		if Less(oldVector[j].Element, localVector[j].Element, L) {
			// Add the pair
			lattice_proofs[localVector[j].Element][j] = Pair {
				V1: deepCopyVector(oldVector),
				V2: deepCopyVector(V)	}

			if VerifyValue(V, j) == false {
				UnionProofs(V, j, P)
			}
		}
	}
}

func Verify(V map[string]SignedElement, P map[string]map[string]Pair) bool {
	for j,_ := range V {		
		if VerifyValue(V, j) == false && VerifyPairSet(P, j, V[j].Element) == false {
			return false
		}
	}
	return true
}