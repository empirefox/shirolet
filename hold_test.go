package shirolet

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestNewHold(t *testing.T) {
	hs := map[string]*Hold{
		"a":           &Hold{grant: newUnit("a")},
		"a:c":         &Hold{grant: newUnit("a:c")},
		"a-a:c":       &Hold{grant: newUnit("a"), deprive: newUnit("a:c")},
		"a:b-a:b:c":   &Hold{grant: newUnit("a:b"), deprive: newUnit("a:b:c")},
		"a:*:c-a:b:c": &Hold{grant: newUnit("a:*:c"), deprive: newUnit("a:b:c")},
		"a-b":         &Hold{grant: newUnit("a")},
		"a:b-a:c":     &Hold{grant: newUnit("a:b")},
		"a:*:c-a:b:x": &Hold{grant: newUnit("a:*:c")},
	}
	Convey("一般解析测试", t, func() {
		for s, o := range hs {
			func(s string, o *Hold) {
				Convey(s, func() {
					So(newHold(s), ShouldResemble, o)
				})
			}(s, o)
		}
	})
}

func TestHoldImpliesAnon(t *testing.T) {
	success := pairs{
		{"a", "a"},
		{"a:b", "a:b"},
		{"a:b:c", "a:b:c"},
		{"a:*:c", "a:*:c"},

		{"b", "b:c"},
		{"b:c", "b:c:d"},
		{"b:*:c", "b:o:c"},

		{"a-a:c", "a:b"},
		{"a-a:c,d,e", "a:b"},
	}
	failure := pairs{
		{"a-a:b", "a:b"},
		{"a-a:b,c,d,e", "a:b"},
	}
	Convey("隐含+Anon", t, func() {
		for _, p := range success {
			func(p pair) {
				Convey(fmt.Sprintf("Hold [%s] 隐含 [%s]", p.k, p.v), func() {
					So(newHold(p.k).impliesAnon(newUnit(p.v)), ShouldBeTrue)
				})
			}(p)
		}
		for _, p := range failure {
			func(p pair) {
				Convey(fmt.Sprintf("Hold [%s] 不隐含 [%s]", p.k, p.v), func() {
					So(newHold(p.k).impliesAnon(newUnit(p.v)), ShouldBeFalse)
				})
			}(p)
		}
	})
}

func TestNewHolds(t *testing.T) {
	Convey("空字符串>>>权限单元组", t, func() {
		result = NewHolds("")
		So(result, ShouldResemble, Holds{newHold("")})
	})

	Convey("单一权限字符串>>>权限单元组", t, func() {
		hus := []string{"table", "access", "bucket"}
		for i := 0; i < 3; i++ {
			pstr := strings.Join(hus[0:i], ":")
			result = NewHolds(pstr)
			So(result, ShouldResemble, Holds{newHold(pstr)})
		}
	})

	Convey("组合权限字符串>>>权限单元组", t, func() {
		hs := []string{"img:upload", "user", "a:b:c", "d:e"}
		holds := Holds{newHold("img:upload"), newHold("user"), newHold("a:b:c"), newHold("d:e")}
		for i := 1; i < 4; i++ {
			pstr := strings.Join(hs[0:i], "|")
			result = NewHolds(pstr)
			So(result, ShouldResemble, holds[0:i])
		}
	})
}
