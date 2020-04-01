package main

/****MESSAGES*****/
type PROPOSAL struct {
	Sender string `json:"sender"`
	V map[string]SignedElement `json:"v"`
	P map[string]map[string]Pair `json:"p"`
	S int `json:"s"`
}

type NACK struct {
	Sender string `json:"sender"`
	V map[string]SignedElement `json:"v"`
	P map[string]map[string]Pair `json:"p"`
	S int `json:"s"`
}

type ACK struct {
	Sender string `json:"sender"`
	S int `json:"s"`
}

/****USEFUL OBJECTS*****/
// SignedElement is an element of the lattice with the signature of a process
type SignedElement struct {
	Element string `json:"element"`
	Signature string `json:"signature"`
}

// A pair of vectors can be found in the entries of lattice_proofs
type Pair struct {
	V1 map[string]SignedElement `json:"v1"`
	V2 map[string]SignedElement `json:"v2"`
}