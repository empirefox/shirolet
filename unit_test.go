package shirolet

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type pair struct {
	k string
	v string
}
type pairs []pair

type row struct {
	t    string
	oks  pairs
	bads pairs
}

func TestImplies(t *testing.T) {
	var rows []row
	equals := pairs{
		{"a", "a"},
		{"a:b", "a:b"},
		{"a:b:c", "a:b:c"},
		{"a:*:c", "a:*:c"},
	}
	badequals := pairs{
		{"a", "b"},
		{"a:b", "a:v"},
		{"a:b:c", "a:b:v"},
		{"a:b:c", "a:v:c"},
		{"a:*:c", "a:*:v"},
	}
	rows = append(rows, row{"相等测试", equals, badequals})

	shorts := pairs{
		{"a", "a:b"},
		{"a:b", "a:b:c"},
		{"a:*:b", "a:o:b"},
	}
	badshorts := pairs{
		{"a:b", "a"},
		{"a:b:c", "a:b"},
		{"a:o:b", "a:*:b"},
	}
	rows = append(rows, row{"长短测试", shorts, badshorts})

	parts := pairs{
		{"a,b", "b"},
		{"a,b", "a"},
		{"a,b", "a:*"},
		{"a,b", "a:c"},
		{"a,b", "a:c:*"},
		{"a,b", "a:c:d"},
		{"a,b", "a:*:c"},
		{"a:b,c", "a:b"},
		{"a:b,c", "a:c"},
		{"a:b,c", "a:b:d"},
		{"a:b,c", "a:b:*"},
		{"a:b:c,d", "a:b:c"},
		{"a:b:c,d", "a:b:d"},
		{"a:b,c:d", "a:b:d"},
		{"a:b,c:d", "a:c:d"},
		{"a1,a2:b1,b2:c1,c2", "a1:b1:c1"},
		{"a1,a2:b1,b2:c1,c2", "a2:b1:c2"},
	}
	badparts := pairs{
		{"b", "a,b"},
		{"a", "a,b"},
		{"a:b", "a,b"},
		{"b:a", "a,b"},
		{"a,b:c", "a,b"},
		{"a,b:c:d", "a,b"},
		{"d:a,b:c", "a,b"},
		{"d:c:a,b", "a,b"},
		{"*:c:a,b", "a,b"},
		{"d:*:a,b", "a,b"},
	}
	rows = append(rows, row{"多节测试", parts, badparts})

	wildcards := pairs{
		{"*", "*"},
		{"*", "a"},
		{"*", "a:b"},
		{"*", "a:b:c"},
		{"*", "a:b:*"},
		{"*", "a:*:c"},
		{"*", "a:*:*"},
		{"a", "a:*"},
		{"a", "a:*:*"},
		{"a:*", "a"},
		{"a:*", "a:b"},
		{"a:*", "a:b:c"},
		{"a:*", "a:b:*"},
		{"a:*", "a:*:c"},
		{"a:*", "a:*:*"},
		{"a:b", "a:b:*"},
		{"a:b:*", "a:b"},
		{"a:b:*", "a:b:c"},
		{"a:b:*", "a:b:c,d,e"},
		{"a:*:c", "a:b:c"},
		{"a:*:c", "a:b,d,e:c"},
	}
	badwildcards := pairs{
		{"a", "*"},
		{"a:b", "*"},
		{"a:b:c", "*"},
		{"a:b:*", "*"},
		{"a:*:*", "*"},
		{"a:*:c", "*"},
		{"b", "a:*"},
		{"a:b", "a:*"},
		{"a:b:c", "a:*"},
		{"a:b:*", "a:*"},
		{"a:*:c", "a:*"},
		{"a:c", "a:b:*"},
		{"a:b:c", "a:b:*"},
		{"a:*:b", "a:b:*"},
	}
	rows = append(rows, row{"通配符测试", wildcards, badwildcards})

	anonWithoutAnon := pairs{
		{"", ""},
		{"anon", ""},
		{"", "anon"},
		{"anon", "anon"},
		{"*", "anon"},
	}
	badAnonWithoutAnon := pairs{
		{"a", "anon"},
		{"a:b", "anon"},
		{"a:b:c", "anon"},
		{"a:*", "anon"},
		{"a:*:c", "anon"},
		{"a:b:*", "anon"},
		{"a:b,c", "anon"},

		{"anon", "a"},
		{"anon", "a:b"},
		{"anon", "a:b:c"},
		{"anon", "a:b:*"},
		{"anon", "a:*:c"},
		{"anon", "a,b"},
	}
	rows = append(rows, row{"不包含anon时的anon测试", anonWithoutAnon, badAnonWithoutAnon})

	for _, r := range rows {
		Convey(r.t, t, func() {
			for _, p := range r.oks {
				func(p pair) {
					Convey(fmt.Sprintf("[%s] 隐含 [%s]", p.k, p.v), func() {
						So(newUnit(p.k).implies(newUnit(p.v)), ShouldBeTrue)
					})
				}(p)
			}
			for _, p := range r.bads {
				func(p pair) {
					Convey(fmt.Sprintf("[%s] 不隐含 [%s]", p.k, p.v), func() {
						So(newUnit(p.k).implies(newUnit(p.v)), ShouldBeFalse)
					})
				}(p)
			}
		})
	}

	anonWithAnon := pairs{
		{"", ""},
		{"anon", ""},
		{"", "anon"},
		{"anon", "anon"},
		{"a", "anon"},
		{"a:b", "anon"},
		{"a:b:c", "anon"},
		{"a:*", "anon"},
		{"a:*:c", "anon"},
		{"a:b:*", "anon"},
		{"a:b,c", "anon"},
		{"*", "anon"},
		{"a1,a2:b1,b2:c1,c2", ""},
	}
	badAnonWithAnon := pairs{
		{"anon", "a"},
		{"anon", "a:b"},
		{"anon", "a:b:c"},
		{"anon", "a:b:*"},
		{"anon", "a:*:c"},
		{"anon", "a,b"},
	}
	Convey("包含anon时的anon测试", t, func() {
		for _, p := range anonWithAnon {
			func(p pair) {
				Convey(fmt.Sprintf("[%s] 隐含 [%s]", p.k, p.v), func() {
					So(newUnit(p.k).impliesAnon(newUnit(p.v)), ShouldBeTrue)
				})
			}(p)
		}
		for _, p := range badAnonWithAnon {
			func(p pair) {
				Convey(fmt.Sprintf("[%s] 不隐含 [%s]", p.k, p.v), func() {
					So(newUnit(p.k).impliesAnon(newUnit(p.v)), ShouldBeFalse)
				})
			}(p)
		}
	})

}

func TestSetPartsWidthErr(t *testing.T) {
	ps := map[string]*unit{
		"a": &unit{
			[]WordSet{NewWordSet("a")},
		},
		"a,b": &unit{
			[]WordSet{NewWordSet("a", "b")},
		},
		"a,b:c": &unit{
			[]WordSet{NewWordSet("a", "b"), NewWordSet("c")},
		},
		"a,b:c,d": &unit{
			[]WordSet{NewWordSet("a", "b"), NewWordSet("c", "d")},
		},
		"a,b:c,d:e": &unit{
			[]WordSet{NewWordSet("a", "b"), NewWordSet("c", "d"), NewWordSet("e")},
		},
		"a,b:c,d:e,f": &unit{
			[]WordSet{NewWordSet("a", "b"), NewWordSet("c", "d"), NewWordSet("e", "f")},
		},
		"a,b:*:e,f": &unit{
			[]WordSet{NewWordSet("a", "b"), NewWordSet("*"), NewWordSet("e", "f")},
		},
	}
	Convey("生成Parts逻辑", t, func() {
		Convey("手动构建对比测试", func() {
			for k, v := range ps {
				So(newUnit(k), ShouldResemble, v)
			}
		})
	})
}
