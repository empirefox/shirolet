shirolet
========

A wildcard permission validate tool for go inspired by apache shiro

# usage
	//if string contains space, use NewPermitRaw
	presource := shirolet.NewPermit("printer:print:no1")
	
	//if string contains space, use NewHoldsRaw
	paccount := shirolet.NewHolds("printer")
	
	ok := presource.SatisfiedBy(paccount)
	// ok => true

# fomate
## pstring on resource
	printer
	printer:*:no2
	a:b | c:e:d ('|' means 'or'. 'a:b' or 'c:e:d' can operate.)
	a:b + c:d ('+' means 'must'. 'a:b' and 'c:d' both are required.)
	a+b | c | d+e (means 'a+b' or 'c' or 'd+e'.)
**'-' is not supported on resource**

## pstring on account
	printer
	printer:*:no2
	a:b | c:e:d (have 2 permission cards 'a:b' and 'c:e:d', '|' means having all.)
	book - book:write (can read every book, but cannot write any. '-' means deprive.)
**'+' is not supported on account**
