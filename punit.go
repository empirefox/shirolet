package shirolet

import (
	"errors"
	"github.com/deckarep/golang-set"
	"github.com/golang/glog"
	"strings"
)

const (
	WildCard      = "*"
	TokenDelim    = ":"
	SubTokenDelim = ","
)

var (
	anon  = "anon"
	panon Punit
)

func init() {
	SetAnno(anon)
}

func NewSet(s []string) mapset.Set {
	a := mapset.NewSet()
	for _, v := range s {
		a.Add(v)
	}
	return a
}

//+gen * methods:"All,Any"
type Punit struct {
	Parts []mapset.Set
}

func NewPunit(pstr string) *Punit {
	p := &Punit{}
	p.setParts(pstr)
	return p
}

func SetAnno(pstr string) {
	err := panon.setPartsWidthErr(pstr)
	if err != nil {
		glog.Errorf(err.Error())
		panon.setPartsWidthErr(anon)
		return
	}
	anon = pstr
}

func (p *Punit) setParts(expr string) {
	err := p.setPartsWidthErr(expr)
	if err != nil {
		glog.Infoln(err.Error())
		p.setPartsWidthErr(anon)
	}
}

func (p *Punit) setPartsWidthErr(expr string) error {
	s := strings.Trim(expr, " "+TokenDelim)
	if len(s) == 0 {
		return errors.New("权限单元格式错误,不能只包含token分界符(:)")
	}
	s = strings.ToLower(s)
	parts := strings.Split(s, TokenDelim)
	p.Parts = make([]mapset.Set, len(parts))
	for i, v := range parts {
		v = strings.Trim(v, " "+SubTokenDelim)
		if len(v) == 0 {
			return errors.New("权限单元格式错误,不能只包含sub token分界符(,)")
		}
		subparts := strings.Split(v, SubTokenDelim)
		p.Parts[i] = NewSet(subparts)
	}
	return nil
}

func (p *Punit) Implies(req *Punit) bool {
	available := p.Parts
	requested := req.Parts
	i := 0
	for _, v := range requested {
		if len(available)-1 < i {
			return true
		} else {
			part := available[i]
			if !part.Contains(WildCard) && !part.IsSubset(v) {
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

func (p *Punit) ImpliesAnon(req *Punit) bool {
	return panon.Implies(req) || p.Implies(req)
}
