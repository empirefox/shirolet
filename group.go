package shirolet

import (
	"strings"
)

//+gen methods:"Any,All"
type group units

func newGroup(perms string) group {
	units := strings.Split(perms, "+")
	lp := len(units)
	result := make(group, lp, lp)
	for i, unit := range units {
		result[i] = newUnit(unit)
	}
	return result
}

//groups

//资源要求的权限
func newGroups(permstring string) groups {
	plusStrSlice := strings.Split(permstring, "|")
	lr := len(plusStrSlice)
	result := make(groups, lr, lr)
	for i, plusStr := range plusStrSlice {
		result[i] = newGroup(plusStr)
	}
	return result
}

//每个单元的满足方式
//请求中的多个perm，有一个满足即可
//使用ImpliesAnon
var satisfiedHold = func(holds Holds) func(*unit) bool {
	return func(required *unit) bool {
		return holds.Any(func(one *Hold) bool {
			return one.impliesAnon(required)
		})
	}
}

//每一段的满足方式
//里面的每个单元要全部满足
var satisfiedPlus = func(holds Holds) func(group) bool {
	return func(pluses group) bool {
		return units(pluses).All(satisfiedHold(holds))
	}
}

//判断用户请求是否足够权限
//p代表资源的需求，holds代表用户的请求
//其中一段满足即可
func (p groups) SatisfiedBy(holds Holds) bool {
	return groups(p).Any(satisfiedPlus(holds))
}
