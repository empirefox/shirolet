// Generated by: gen
// TypeWriter: gen
// Directive: +gen on *shirolet.unit

// See http://clipperhouse.github.io/gen for documentation

package shirolet

// units is a slice of type *unit. Use it where you would use []*unit.
type units []*unit

// All verifies that all elements of units return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func (rcv units) All(fn func(*unit) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Any verifies that one or more elements of units return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func (rcv units) Any(fn func(*unit) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}
