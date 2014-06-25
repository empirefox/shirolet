package shirolet

import (
	"regexp"
)

type Permit interface {
	SatisfiedBy(Holds) bool
}

//资源要求的权限
func NewPermit(ps string) Permit {
	return newGroups(ps)
}

//资源要求的权限
func NewPermitRaw(ps string) Permit {
	return newGroups(Fmt(ps))
}

//用户请求权限
func NewHolds(ps string) Holds {
	return newHolds(ps)
}

func NewHoldsRaw(ps string) Holds {
	return newHolds(Fmt(ps))
}

//设置匿名权限或空字符串权限
func SetAnno(ps string) {
	setAnno(ps)
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
