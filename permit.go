package shirolet

import (
	"regexp"
	"strings"
)

//+gen methods:"Any,All"
type Pgroup Punits

func NewPgroup(perms, split string) Pgroup {
	units := strings.Split(perms, split)
	lp := len(units)
	result := make(Pgroup, lp, lp)
	for i, unit := range units {
		result[i] = NewPunit(unit)
	}
	return result
}

//用户请求权限
func NewHolds(perms string) Pgroup {
	return NewPgroup(perms, "|")
}

func NewHoldsRaw(ps string) Pgroup {
	return NewHolds(Fmt(ps))
}

type Permit Pgroups

//资源要求的权限
func NewPermit(permstring string) Permit {
	plusStrSlice := strings.Split(permstring, "|")
	lr := len(plusStrSlice)
	result := make(Permit, lr, lr)
	for i, plusStr := range plusStrSlice {
		result[i] = NewPgroup(plusStr, "+")
	}
	return result
}

func NewPermitRaw(ps string) Permit {
	return NewPermit(Fmt(ps))
}

var (
	fmtReg1 = regexp.MustCompile(`^\s*\||\s|\|\s*$`)
	fmtReg2 = regexp.MustCompile(`\|{2,}`)
)

//格式化权限字符串，一般在存入数据库前使用
func Fmt(ps string) string {
	ps = fmtReg1.ReplaceAllString(ps, "")
	return fmtReg2.ReplaceAllString(ps, "|")
}

//每个单元的满足方式
//请求中的多个perm，有一个满足即可
var satisfiedUnit = func(holds Pgroup) func(*Punit) bool {
	return func(required *Punit) bool {
		return Punits(holds).Any(func(one *Punit) bool {
			return one.ImpliesAnon(required)
		})
	}
}

//每一段的满足方式
//里面的每个单元要全部满足
var satisfiedPlus = func(holds Pgroup) func(Pgroup) bool {
	return func(pluses Pgroup) bool {
		return Punits(pluses).All(satisfiedUnit(holds))
	}
}

//判断用户请求是否足够权限
//p代表资源的需求，holds代表用户的请求
//其中一段满足即可
func (p Permit) SatisfiedBy(holds Pgroup) bool {
	return Pgroups(p).Any(satisfiedPlus(holds))
}
