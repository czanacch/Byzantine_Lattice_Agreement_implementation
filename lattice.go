package main


/***************************************** LATTICE **********************************************/
type Lattice struct {
	Elements map[string]*Element // Set of elements contained in the lattice
}

func (l *Lattice) constructorEmptyLattice() {
	l.Elements = make(map[string]*Element)
	bottom := new(Element)
	top := new(Element)
	bottom.Id = "⊥"
	top.Id = "⊤"
	bottom.Uppers = append(bottom.Uppers, top)
	top.Lowers = append(top.Lowers, bottom)
	
	l.Elements["⊤"] = top
	l.Elements["⊥"] = bottom
}

type Element struct {
	Id string // Element identifier
	Uppers []*Element // Element upperbounds
	Lowers []*Element // Element lowerbounds
}


func (e *Element) AddElement(id string, lowers []*Element, uppers []*Element, lattice *Lattice) {
	e.Id = id
	e.Uppers = append(e.Uppers, uppers...)
	for _, element := range uppers {
		element.Lowers = append(element.Lowers, e)
	}

	e.Lowers = append(e.Lowers, lowers...)
	for _, element := range lowers {
		element.Uppers = append(element.Uppers, e)
	}

	lattice.Elements[e.Id] = e
}

func all_upperbounds(upperbounds []*Element, e *Element, lattice *Lattice) []*Element{

	if e == lattice.Elements["⊤"] {
		return []*Element{}
	}

	for _, upper := range e.Uppers {
		if !contains(upperbounds,upper) {
			upperbounds = append(upperbounds, upper)
		}
		for _, element := range all_upperbounds(upperbounds, upper, lattice) { // Deep recursion
			if !contains(upperbounds,element) {
				upperbounds = append(upperbounds, element)
			}
		}
	}

	return upperbounds	
}

// It returns true if the slice s contains e
func contains(s []*Element, e *Element) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

// Implementation of join operation between elements
func Join(value1 string, value2 string, lattice *Lattice) *Element {

	e1 := lattice.Elements[value1]
	e2 := lattice.Elements[value2]

	e1_upperbounds := all_upperbounds([]*Element{}, e1 , lattice)
	e2_upperbounds := all_upperbounds([]*Element{}, e2 , lattice)

	// Base cases
	if e1 == e2 { // If the elements are the same, I return one of the two
		return e1
	}
	// If one is ⊥, I return the other element
	if e1 == lattice.Elements["⊥"] {
		return e2
	}
	if e2 == lattice.Elements["⊥"] {
		return e1
	}

	// If one is ⊤, I return ⊤
	if e1 == lattice.Elements["⊤"] || e2 == lattice.Elements["⊤"] {
		return lattice.Elements["⊤"]
	}

	common_upperBounds := []*Element{} // Lista di tutti gli upperbound comuni
	for _, element1 := range e1_upperbounds {
		for _, element2 := range e2_upperbounds {
			if element1.Id == element2.Id {
				common_upperBounds = append(common_upperBounds, element1)
			}
		}
	}

	max := "⊤"
	for _,element := range common_upperBounds {
		if Less(element.Id, max, lattice) {
			max = element.Id
		}
	}
	return lattice.Elements[max]
}

// Implementation of < operation between elements
func Less(value1 string, value2 string, lattice *Lattice) bool {
	
	e1 := lattice.Elements[value1]
	e2 := lattice.Elements[value2]
	
	e1_upperbounds := all_upperbounds([]*Element{}, e1 , lattice)

	if contains(e1_upperbounds,e2) {
		return true
	} else {
		return false
	}

}

// Implementation of <= operation between elements
func LessEqual(value1 string, value2 string, lattice *Lattice) bool {
	if Less(value1,value2,lattice) || value1 == value2 {
		return true
	} else {
		return false
	}

}


// Create an example of lattice
func exampleLattice() *Lattice {
	my_lattice := new(Lattice)
	my_lattice.constructorEmptyLattice()
	
	a_element := new(Element)
	a_element.AddElement("a", []*Element{my_lattice.Elements["⊥"]}, []*Element{}, my_lattice)

	b_element := new(Element)
	b_element.AddElement("b", []*Element{my_lattice.Elements["⊥"]}, []*Element{}, my_lattice)

	c_element := new(Element)
	c_element.AddElement("c", []*Element{my_lattice.Elements["⊥"]}, []*Element{}, my_lattice)

	d_element := new(Element)
	d_element.AddElement("d", []*Element{my_lattice.Elements["⊥"]}, []*Element{}, my_lattice)
	

	e_element := new(Element)
	e_element.AddElement("e", []*Element{my_lattice.Elements["a"], my_lattice.Elements["b"]}, []*Element{}, my_lattice)

	f_element := new(Element)
	f_element.AddElement("f", []*Element{my_lattice.Elements["a"], my_lattice.Elements["c"]}, []*Element{}, my_lattice)

	g_element := new(Element)
	g_element.AddElement("g", []*Element{my_lattice.Elements["a"], my_lattice.Elements["d"]}, []*Element{}, my_lattice)

	h_element := new(Element)
	h_element.AddElement("h", []*Element{my_lattice.Elements["b"], my_lattice.Elements["c"]}, []*Element{}, my_lattice)
	
	i_element := new(Element)
	i_element.AddElement("i", []*Element{my_lattice.Elements["b"], my_lattice.Elements["d"]}, []*Element{}, my_lattice)

	l_element := new(Element)
	l_element.AddElement("l", []*Element{my_lattice.Elements["c"], my_lattice.Elements["d"]}, []*Element{}, my_lattice)


	m_element := new(Element)
	m_element.AddElement("m", []*Element{my_lattice.Elements["e"], my_lattice.Elements["f"], my_lattice.Elements["h"]}, []*Element{my_lattice.Elements["⊤"]}, my_lattice)

	n_element := new(Element)
	n_element.AddElement("n", []*Element{my_lattice.Elements["e"], my_lattice.Elements["g"], my_lattice.Elements["i"]}, []*Element{my_lattice.Elements["⊤"]}, my_lattice)
	
	o_element := new(Element)
	o_element.AddElement("o", []*Element{my_lattice.Elements["f"], my_lattice.Elements["g"], my_lattice.Elements["l"]}, []*Element{my_lattice.Elements["⊤"]}, my_lattice)

	p_element := new(Element)
	p_element.AddElement("p", []*Element{my_lattice.Elements["h"], my_lattice.Elements["i"], my_lattice.Elements["l"]}, []*Element{my_lattice.Elements["⊤"]}, my_lattice)
	

	return my_lattice
}