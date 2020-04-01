package main

import (
	"sync"
	"net/http"
	"fmt"
)

var processes []string  // IP addresses of all processes (remote machines) (its length returns n)
var my_address string // IP address of the current process
var my_value string // Value proposed by the current process
var L *Lattice // Global lattice known by all processes
var arrivedACK map[int]map[string]bool // variable indicating that an ACK came from the p_j process or not for the current proposal
var arrivedNACK map[int]map[string]bool // variable indicating that a NACK came from the p_j process or not for the current proposal

var f int // maximum Byzantine tolerance
var mutex sync.Mutex
var channel_Proposal chan PROPOSAL
var channel_Ack chan ACK
var channel_Nack chan NACK

var status string
var ackCount int
var nackCount int
var proposalNumber int
var decision string
var localVector map[string]SignedElement
var oldVector map[string]SignedElement
var proposedVector map[string]SignedElement
var outputVector map[string]SignedElement

var nack_arrivals map[string]map[string]SignedElement

var lattice_proofs map[string]map[string]Pair


func main() {
	mutex = sync.Mutex{}
	channel_Proposal = make(chan PROPOSAL)
	channel_Ack = make(chan ACK)
	channel_Nack = make(chan NACK)

	initialization() // Inizializza tutte le variabili

	fmt.Println("Activation of the process", my_address)

	go SendingPROPOSAL()
	go FeedbackExecution()

	http.HandleFunc("/proposalArrival", handleProposalArrival)

	http.HandleFunc("/receivingACK", handleReceivingACK)

	http.HandleFunc("/receivingNACK", handleReceivingNACK)
	
	http.ListenAndServe(":8000", nil)
}

func initialization() {

	mutex.Lock()

	L = exampleLattice() // Create an example of lattice
	my_address = "127.0.0.1:8000"
	processes = []string{"127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003"}

	f = len(processes)/2 - 1

	my_value = "a" // Element proposed by p_i
	
	status = "passive"
	ackCount = 0
	nackCount = 0
	proposalNumber = 0
	decision  = "⊥"
	arrivedACK = make(map[int]map[string]bool)
	arrivedNACK = make(map[int]map[string]bool)

	localVector = make(map[string]SignedElement)
	oldVector = make(map[string]SignedElement)
	proposedVector = make(map[string]SignedElement)
	outputVector = make(map[string]SignedElement)

	// Initialization of arrivedACK and arrivedNACK
	arrivedACK[0] = make(map[string]bool)
	arrivedNACK[0] = make(map[string]bool)
	for _,process := range processes {
		arrivedACK[0][process] = false
		arrivedNACK[0][process] = false
	}
	

	// Vectors initialization
	for _,process := range processes {
		localVector[process] = SignedElement{
			Element: "⊥",
			Signature: ""}
		
		oldVector[process] = SignedElement{
			Element: "⊥",
			Signature: ""}

		proposedVector[process] = SignedElement{
			Element: "⊥",
			Signature: ""}

		outputVector[process] = SignedElement{
			Element: "⊥",
			Signature: ""}
	}

	// Initialization of nack_arrivals matrix
	nack_arrivals = make(map[string]map[string]SignedElement)
	for _, process := range processes {
		nack_arrivals[process] = make(map[string]SignedElement)
	}

	for _, processRow := range processes {
		for _, processColumn := range processes {
			nack_arrivals[processRow][processColumn] = SignedElement{
				Element: "⊥",
				Signature: ""}
		}
	}

	// Initialization of lattice_proofs map
	lattice_proofs = make(map[string]map[string]Pair)
	
	for _,element := range L.Elements {
		lattice_proofs[element.Id] = make(map[string]Pair)
	}

	for _, keys := range L.Elements {
		for _, processColumn := range processes {
			lattice_proofs[keys.Id][processColumn] = Pair{
				V1: nil,
				V2: nil }
		}
	}

}