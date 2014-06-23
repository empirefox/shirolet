package shirolet

import (
	"testing"
)

func should(this, other, err string, t *testing.T) {
	p1 := NewPunit(this)
	p2 := NewPunit(other)
	if !p1.Implies(p2) {
		t.Errorf(err, this, other)
	}
}

func shouldNot(this, other, err string, t *testing.T) {
	p1 := NewPunit(this)
	p2 := NewPunit(other)
	if p1.Implies(p2) {
		t.Errorf(err, this, other)
	}
}

func shouldWithAnon(this, other, err string, t *testing.T) {
	p1 := NewPunit(this)
	p2 := NewPunit(other)
	if !p1.ImpliesAnon(p2) {
		t.Errorf(err, this, other)
	}
}

func shouldNotWithAnon(this, other, err string, t *testing.T) {
	p1 := NewPunit(this)
	p2 := NewPunit(other)
	if p1.ImpliesAnon(p2) {
		t.Errorf(err, this, other)
	}
}

func equal(this, other string, t *testing.T) {
	should(this, other, "%s == %s 模式不相等", t)
}

func Test_Equal(t *testing.T) {
	pstr1 := "user"
	pstr2 := "user:read"
	pstr3 := "user:read:123"
	pstr4 := "user:*:123"

	equal(pstr1, pstr1, t)
	equal(pstr2, pstr2, t)
	equal(pstr3, pstr3, t)
	equal(pstr4, pstr4, t)
	equal("anon", "", t)
	equal("", "anon", t)
}

func notEqual(this, other string, t *testing.T) {
	shouldNot(this, other, "%s should not implies %s (euqal)", t)
}

func Test_NotEqual(t *testing.T) {
	pstr1 := "user"
	pstr2 := "user:read"
	pstr3 := "user:read:123"

	notEqual(pstr1, "admin", t)
	notEqual(pstr2, "user:write", t)
	notEqual(pstr2, "admin:read", t)
	notEqual(pstr3, "user:write:123", t)
	notEqual(pstr3, "user:read:456", t)
	notEqual(pstr3, "admin:read:123", t)
}

func shorter(this, other string, t *testing.T) {
	should(this, other, "%s should implies %s 短模式", t)
}

func Test_Shorter(t *testing.T) {
	pstr1 := "user"
	pstr2 := "user:read"
	pstr3 := "user:read:123"

	shorter(pstr1, pstr2, t)
	shorter(pstr1, pstr3, t)
	shorter(pstr2, pstr3, t)
}

func longer(this, other string, t *testing.T) {
	shouldNot(this, other, "错误: %s should not implies %s 长模式", t)
}

func Test_Longer(t *testing.T) {
	pstr1 := "user"
	pstr2 := "user:read"
	pstr3 := "user:read:123"

	longer(pstr1, "", t)
	longer(pstr2, pstr1, t)
	longer(pstr3, pstr1, t)
	longer(pstr3, pstr2, t)
}

func wildcard(this, other string, t *testing.T) {
	should(this, other, "错误: %s should implies %s *模式", t)
}

func Test_Wildcard(t *testing.T) {
	pstr1 := "user"
	pstr2 := "user:read"
	pstr3 := "user:read:123"

	wildcard("user:*", pstr1, t)
	wildcard(pstr1, "user:*", t)
	wildcard("user:*:*", pstr1, t)
	wildcard("user:*", pstr2, t)
	wildcard("user:*:*", pstr2, t)
	wildcard("user:read:*", pstr2, t)
	wildcard(pstr2, "user:read:*", t)
	wildcard("user:*", pstr3, t)
	wildcard("user:*:*", pstr3, t)
	wildcard("user:*:123", pstr3, t)
	wildcard("user:read:*", pstr3, t)
}

func wildcardNot(this, other string, t *testing.T) {
	shouldNot(this, other, "错误: %s should not implies %s *模式", t)
}

func Test_WildcardNot(t *testing.T) {
	pstr1 := "user"
	pstr2 := "user:read"
	pstr3 := "user:read:123"

	wildcardNot(pstr1, "*", t)
	wildcardNot(pstr2, "user:*", t)
	wildcardNot(pstr2, "user:*:*", t)
	wildcardNot(pstr2, "user:*:123", t)
	wildcardNot(pstr3, "user:*", t)
	wildcardNot(pstr3, "user:*:*", t)
	wildcardNot(pstr3, "user:read:*", t)
	wildcardNot(pstr3, "user:*:123", t)
}

func isAnon(this, other string, t *testing.T) {
	shouldWithAnon(this, other, "错误: %s should implies %s (anon)", t)
}

func Test_IsAnon(t *testing.T) {
	isAnon("", "anon", t)
	isAnon("user", "anon", t)
	isAnon("user", "", t)
	isAnon("user:read", "", t)
	isAnon("user:read:123", "anon", t)
}

func notAnon(this, other string, t *testing.T) {
	shouldNot(this, other, "错误: %s should not implies %s (not anon)", t)
}

func Test_WithOutAnon(t *testing.T) {
	notAnon("user", "anon", t)
	notAnon("user", "", t)
	notAnon("user:read", "", t)
	notAnon("user:read:123", "anon", t)
}
