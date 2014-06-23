package shirolet

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"strings"
	"testing"
)

var (
	result Pgroup
)

func TestNewPgroup(t *testing.T) {
	Convey("空字符串>>>权限单元组", t, func() {
		result = NewPgroup("", "|")
		So(result, ShouldResemble, Pgroup{NewPunit("")})
	})

	Convey("单一权限字符串>>>权限单元组", t, func() {
		uus := []string{"table", "access", "bucket"}
		for i := 0; i < 3; i++ {
			Convey(strconv.Itoa(i+1)+"节", func() {
				pstr := strings.Join(uus[0:i], ":")
				result = NewPgroup(pstr, "|")
				So(result, ShouldResemble, Pgroup{NewPunit(pstr)})
			})
		}
	})

	Convey("组合权限字符串>>>权限单元组", t, func() {
		us := []string{"img:upload", "user", "a:b:c", "d:e"}
		pus := Pgroup{NewPunit("img:upload"), NewPunit("user"), NewPunit("a:b:c"), NewPunit("d:e")}
		for i := 1; i < 4; i++ {
			Convey(strconv.Itoa(i+1)+"组", func() {
				pstr := strings.Join(us[0:i], "|")
				result = NewPgroup(pstr, "|")
				So(result, ShouldResemble, pus)
			})
		}
	})
}

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

func TestPermitSatisfiedByComposed(t *testing.T) {
	p1 := "a:b:c + d:e + f | g:h + i | j"
	r1ts := []string{"a:b:c | d:e | f", "g:h | i", "j", "a|d|f", "g|i", "j|h", "*"}
	r1fs := []string{"a:b:false | d:e | f", "g | i:j", "j:*:k"}
	permit1 := NewPermitRaw(p1)
	Convey("资源复合权限:"+p1, t, func() {
		for i := 0; i < 2; i++ {
			Convey(fmt.Sprintf("第%s次验证》》》", strconv.Itoa(i+1)), func() {
				for _, r1t := range r1ts {
					r1 := NewHoldsRaw(r1t)
					Convey("请求权限应当成功:"+r1t, func() {
						So(permit1.SatisfiedBy(r1), ShouldBeTrue)
					})
				}
				for _, r1f := range r1fs {
					r1 := NewHoldsRaw(r1f)
					Convey("请求权限应当失败:"+r1f, func() {
						So(permit1.SatisfiedBy(r1), ShouldBeFalse)
					})
				}
			})
		}
	})
}
