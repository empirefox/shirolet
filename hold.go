package shirolet

import (
	"regexp"
	"strings"
)

//+gen * methods:"Any,All"
type Hold struct {
	grant   *unit
	deprive *unit
}

func (h *Hold) implies(unit *unit) bool {
	if h.deprive != nil && h.deprive.implies(unit) {
		return false
	}
	return h.grant.implies(unit)
}

func (h *Hold) impliesAnon(unit *unit) bool {
	if h.deprive != nil && h.deprive.implies(unit) {
		return false
	}
	return h.grant.impliesAnon(unit)
}

var (
	newHoldReg1 = regexp.MustCompile(`\s`)
	newHoldReg2 = regexp.MustCompile(`^[^-]+-[^-]+$`)
	newHoldReg3 = regexp.MustCompile(`^[^-:]+:?.*$`)
)

//允许剥夺anon权限
func newHold(perms string) *Hold {
	h := &Hold{}
	perms = newHoldReg1.ReplaceAllString(perms, "")
	if !newHoldReg2.MatchString(perms) && !newHoldReg3.MatchString(perms) {
		h.grant = newUnit("")
		return h
	}
	ss := strings.Split(perms, "-")
	switch len(ss) {
	case 1:
		h.grant = newUnit(ss[0])
	case 2:
		h.grant = newUnit(ss[0])
		deprive := newUnit(ss[1])
		if h.grant.impliesAnon(deprive) {
			h.deprive = deprive
		}
	}
	return h
}

//holds

//用户请求权限
func newHolds(perms string) Holds {
	hs := strings.Split(perms, "|")
	l := len(hs)
	result := make(Holds, l, l)
	for i, h := range hs {
		result[i] = newHold(h)
	}
	return result
}
