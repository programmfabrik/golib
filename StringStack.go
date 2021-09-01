package lib

type StringStack []string

// IsEmpty: check if stack is empty
func (s StringStack) IsEmpty() bool {
	return len(s) == 0
}

// Push a new value onto the stack
func (s *StringStack) Push(str string) {
	*s = append(*s, str)
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *StringStack) Pop() string {
	if s.IsEmpty() {
		return ""
	}
	index := len(*s) - 1   // Get the index of the top most element.
	element := (*s)[index] // Index into the slice and obtain the element.
	*s = (*s)[:index]      // Remove it from the stack by slicing it off.
	return element
}

// Return the last element of the stack without removing it
func (s *StringStack) Peek() string {
	if s.IsEmpty() {
		return ""
	}
	return (*s)[len(*s)-1]
}
