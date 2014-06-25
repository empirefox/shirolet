package shirolet

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	result Holds
)

func TestPermitSatisfiedByNull(t *testing.T) {
	Convey("空资源权限，应当全部通过", t, func() {
		p := NewPermitRaw("")
		pts := []string{"", "*", "a", "a:b", "a:b:c", "a:*:c", "a|b", "a:a|b:c|d:e"}
		for _, rt := range pts {
			r := NewHoldsRaw(rt)
			So(p.SatisfiedBy(r), ShouldBeTrue)
		}
	})
	Convey("请求权限为空，应当全部失败", t, func() {
		ps := []string{"*", "a", "a:b", "a:b:c", "a:*:c", "a|b", "a:a|b:c|d:e"}
		r := NewHoldsRaw("")
		for _, p := range ps {
			permit := NewPermitRaw(p)
			So(permit.SatisfiedBy(r), ShouldBeFalse)
		}
	})
}

func TestPermitSatisfiedByOr(t *testing.T) {
	ps := pairs{
		{"a|b", "a"},
		{"a|b|c", "b"},
		{"a|b|c|d", "a:*"},
		{"a|b|c|e", "a:*:*"},
		{"a|b|c|d|e", "d"},
	}
	Convey("单一的 '|' 分割格式权限测试", t, func() {
		for _, p := range ps {
			func(k, v string) {
				Convey(fmt.Sprintf("[%s] 应当满足 [%s]", v, k), func() {
					So(NewPermit(k).SatisfiedBy(NewHolds(v)), ShouldBeTrue)
				})
			}(p.k, p.v)
		}
	})
}

func TestPermitSatisfiedByComposed(t *testing.T) {
	p1 := "a:b:c + d:e + f | g:h + i | j"
	success := []string{"a:b:c | d:e | f", "g:h | i", "j", "a|d|f", "g|i", "j|h", "*"}
	failure := []string{"a:b:false | d:e | f", "g | i:j", "j:*:k"}
	permit1 := NewPermitRaw(p1)
	Convey("资源复合权限: "+p1, t, func() {
		Convey("请求权限应当成功", func() {
			for _, p := range success {
				func(p string) {
					Convey(p, func() {
						fmt.Println(p)
						So(permit1.SatisfiedBy(NewHoldsRaw(p)), ShouldBeTrue)
					})
				}(p)
			}
		})
		Convey("请求权限应当失败", func() {
			for _, p := range failure {
				func(p string) {
					Convey(p, func() {
						fmt.Println(p)
						So(permit1.SatisfiedBy(NewHoldsRaw(p)), ShouldBeFalse)
					})
				}(p)
			}
		})
	})
}

func TestDeprived(t *testing.T) {
	success := map[string]string{
		"a:b":     "a-a:c",
		"a:b,c":   "a-a:d",
		"a:b|c":   "a-a:c|d-d:c",
		"a:b|c:d": "a-a:b:c|c:d-c:d:e",
	}
	failure := map[string]string{
		"a:b":     "a-a:b",
		"a:b+c":   "a-a:b",
		"a:b|c":   "a-a:d,b,c",
		"a:b|c:d": "a-a:d,b,c|c-c:a,b,c,d,e",
	}
	Convey("剥夺权限", t, func() {
		Convey("应当通过的测试", func() {
			for k, v := range success {
				func(k, v string) {
					Convey(fmt.Sprintf("[%s]可以由[%s]访问", k, v), func() {
						So(NewPermit(k).SatisfiedBy(NewHolds(v)), ShouldBeTrue)
					})
				}(k, v)
			}
		})
		Convey("应当失败的测试", func() {
			for k, v := range failure {
				func(k, v string) {
					Convey(fmt.Sprintf("[%s]拒绝[%s]的访问", k, v), func() {
						So(NewPermit(k).SatisfiedBy(NewHolds(v)), ShouldBeFalse)
					})
				}(k, v)
			}
		})
	})
}
