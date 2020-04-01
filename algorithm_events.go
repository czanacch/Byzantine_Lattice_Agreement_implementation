package main

import (
	"time"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
)


// EVENT 1: SENDING PROPOSAL
func SendingPROPOSAL() {
	status = "active"
	localVector[my_address] = SignedElement{
		Element: my_value,
		Signature: my_address}
	proposedVector = localVector
	proposalNumber++

	arrivedACK[proposalNumber] = make(map[string]bool)
	arrivedNACK[proposalNumber] = make(map[string]bool)
	for _,process := range processes {
		arrivedACK[proposalNumber][process] = false
		arrivedNACK[proposalNumber][process] = false
	}

	fmt.Println("I broadcast Proposal # ", proposalNumber)
	time.Sleep(2 * time.Second)
	for _,process := range processes {
		jsonData := PROPOSAL{
			Sender: my_address,
			V: proposedVector,
			P: lattice_proofs,
			S: proposalNumber }
		jsonValue, _ := json.Marshal(jsonData)

		go http.Post("http://"+process+"/proposalArrival", "application/json", bytes.NewBuffer(jsonValue))
	}

	mutex.Unlock()
}

// EVENT 2: PROPOSAL ARRIVAL
func handleProposalArrival(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	var proposal PROPOSAL
	_ = json.NewDecoder(r.Body).Decode(&proposal)

	go manageProposal()
	channel_Proposal <- proposal
}

func manageProposal() {
	mutex.Lock()
	proposal := <-channel_Proposal

	V := proposal.V
	P := proposal.P
	s := proposal.S

	fmt.Print("I received Proposal # ", proposalNumber, " from process ", proposal.Sender, " containing " )
	printVector(V)
	fmt.Println(" ")

	if Verify(V, P) {
		oldVector = localVector

		localVector = JoinVector(localVector, V, L)
		AddPair(V, P)
			
		if EqualVector(localVector, V) {
			time.Sleep(2 * time.Second)

			fmt.Println("I send an ACK to ", proposal.Sender, " related to Proposal # ",  proposalNumber)
			jsonData := ACK{
				Sender: my_address,
				S: s}
			jsonValue, _ := json.Marshal(jsonData)
			go http.Post("http://"+proposal.Sender+"/receivingACK", "application/json", bytes.NewBuffer(jsonValue))

		} else {
			time.Sleep(2 * time.Second)

			fmt.Print("I send a NACK to ", proposal.Sender, " related to Proposal # ",  proposalNumber, " containing ")
			printVector(localVector)
			fmt.Println(" ")

			jsonData := NACK{
				Sender: my_address,
				V: localVector,
				P: lattice_proofs,
				S: s }
			jsonValue, _ := json.Marshal(jsonData)
	
			go http.Post("http://"+proposal.Sender+"/receivingNACK", "application/json", bytes.NewBuffer(jsonValue))
		}

	}
	mutex.Unlock()
}

// EVENT 3: RECEIVING ACK
func handleReceivingACK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	var ack ACK
	_ = json.NewDecoder(r.Body).Decode(&ack)

	go manageACK()
	channel_Ack <- ack
}

func manageACK() {
	mutex.Lock()
	ack := <-channel_Ack
	s := ack.S
	j := ack.Sender

	if s == proposalNumber && arrivedACK[s][j] == false {
		fmt.Println("I received a correct ACK from ", j, " for my Proposal # ", s)
		arrivedACK[s][j] = true
		ackCount++
	}
	mutex.Unlock()
}

// EVENT 4: RECEIVING NACK
func handleReceivingNACK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	var nack NACK
	_ = json.NewDecoder(r.Body).Decode(&nack)

	go manageNACK()
	channel_Nack <- nack
}

func manageNACK() {
	mutex.Lock()
	nack:=  <-channel_Nack
	V := nack.V
	P := nack.P
	s := nack.S
	j := nack.Sender

	if Verify(V, P) && s == proposalNumber && arrivedNACK[s][j] == false && LessVector(nack_arrivals[j], V, L) {  
		fmt.Println("I received a correct NACK from ", j, " per my Proposal # ", s)
		arrivedNACK[s][j] = true
		nack_arrivals[j] = deepCopyVector(V)
		oldVector = localVector
		localVector = JoinVector(localVector, V, L)
		AddPair(V, P)
		nackCount++
	}
	mutex.Unlock()
}

// EVENT 5: FEEDBACK EXECUTION
func FeedbackExecution() {
	for {
		if (ackCount + nackCount) >= f + 2 && status == "active" {
			mutex.Lock()
			if ackCount == ackCount + nackCount {
				Decide()
			} else {
				Refine()
			}
			mutex.Unlock()
		}
		time.Sleep(1 * time.Millisecond)
	}
}


func Refine() {
	proposedVector = localVector
	proposalNumber++
	ackCount = 0
	nackCount = 0

	arrivedACK[proposalNumber] = make(map[string]bool)
	arrivedNACK[proposalNumber] = make(map[string]bool)
	for _,process := range processes {
		arrivedACK[proposalNumber][process] = false
		arrivedNACK[proposalNumber][process] = false
	}

	time.Sleep(2 * time.Second)

	for _,process := range processes {
		jsonData := PROPOSAL{
			Sender: my_address,
			V: proposedVector,
			P: lattice_proofs,
			S: proposalNumber }
		jsonValue, _ := json.Marshal(jsonData)

		go http.Post("http://"+process+"/proposalArrival", "application/json", bytes.NewBuffer(jsonValue))
	}
}

func Decide() {
	outputVector = proposedVector
	status = "passive"

	for j,_ := range proposedVector {
		decision = Join(decision, proposedVector[j].Element,L).Id
	}
	fmt.Println("I DECIDE FOR ", decision)
}