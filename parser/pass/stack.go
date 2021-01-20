package pass

// based on https://www.educative.io/edpresso/how-to-implement-a-stack-in-golang

// Variable represents the stack of a variable
type Variable struct {
	IsDefined bool
	Slot      int
}

// ScopeStack data structure
type ScopeStack []map[string]*Variable

// IsEmpty check if stack is empty
func (s *ScopeStack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *ScopeStack) Push(elem map[string]*Variable) {
	*s = append(*s, elem) // Simply append the new value to the end of the stack
}

// Pop Remove and return top element of stack. Return false if stack is empty.
func (s *ScopeStack) Pop() (map[string]*Variable, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	index := len(*s) - 1   // Get the index of the top most element.
	element := (*s)[index] // Index into the slice and obtain the element.
	*s = (*s)[:index]      // Remove it from the stack by slicing it off.
	return element, true
}

// Peek check top element of stack without removing it
func (s *ScopeStack) Peek() (map[string]*Variable, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	index := len(*s) - 1 // Get the index of the top most element.
	return (*s)[index], true
}
