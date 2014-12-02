package shirolet

import (
	"errors"
	"strings"

	"github.com/golang/glog"
)

const (
	WildCard      = "*"
	TokenDelim    = ":"
	SubTokenDelim = ","
)

var (
	anon  = "anon"
	panon unit
)

func init() {
	setAnno(anon)
}

//+gen * methods:"All,Any"
type unit struct {
	Parts []WordSet
}

func newUnit(pstr string) *unit {
	u := &unit{}
	u.setParts(pstr)
	return u
}

func setAnno(pstr string) {
	err := panon.setPartsWidthErr(pstr)
	if err != nil {
		glog.Errorf(err.Error())
		panon.setPartsWidthErr(anon)
		return
	}
	anon = pstr
}

func (u *unit) setParts(expr string) {
	err := u.setPartsWidthErr(expr)
	if err != nil {
		u.setPartsWidthErr(anon)
	}
}

func (u *unit) setPartsWidthErr(expr string) error {
	s := strings.Trim(expr, " "+TokenDelim)
	if len(s) == 0 {
		return errors.New("权限单元格式错误,不能只包含token分界符(:)," + s)
	}
	s = strings.ToLower(s)
	parts := strings.Split(s, TokenDelim)
	u.Parts = make([]WordSet, len(parts))
	for i, v := range parts {
		v = strings.Trim(v, " "+SubTokenDelim)
		if len(v) == 0 {
			return errors.New("权限单元格式错误,不能只包含sub token分界符(,)," + v)
		}
		subparts := strings.Split(v, SubTokenDelim)
		u.Parts[i] = NewWordSet(subparts...)
	}
	return nil
}

func (u *unit) implies(req *unit) bool {
	if req == nil {
		return true
	}
	available := u.Parts
	required := req.Parts
	i := 0
	for _, v := range required {
		if len(available)-1 < i {
			return true
		} else {
			part := available[i]
			if !part.Contains(WildCard) && !part.IsSuperset(v) {
				return false
			}
			i++
		}
	}
	for ; i < len(available); i++ {
		part := available[i]
		if !part.Contains(WildCard) {
			return false
		}
	}
	return true
}

func (u *unit) impliesAnon(req *unit) bool {
	return panon.implies(req) || u.implies(req)
}
