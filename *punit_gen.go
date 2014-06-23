// Generated by: gen
// TypeWriter: gen
// Directive: +gen on *shirolet.Punit

// See http://clipperhouse.github.io/gen for documentation

package shirolet

// Punits is a slice of type *Punit. Use it where you would use []*Punit.
type Punits []*Punit

// All verifies that all elements of Punits return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func (rcv Punits) All(fn func(*Punit) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Any verifies that one or more elements of Punits return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func (rcv Punits) Any(fn func(*Punit) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}