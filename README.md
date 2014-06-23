shirolet
========

A wildcard permission validate tool for go inspired by apache shiro

# usage
	presource := shirolet.NewPermit("printer:print:no1")
	paccount := shirolet.NewHolds("printer")
	ok := presource.SatisfiedBy(paccount)
	// ok => true

# fomate
## pstring on resource
	printer
	printer:*:no2
	a:b|c:e:d ('|' means 'or')
	a:b+c:d ('+' means 'and' that everyone is required)
	a+b | c | d+e (means 'a+b' or 'c' or 'd+e')
## pstring on account
	printer
	printer:*:no2
	a:b|c:e:d ('|' means having all)
	'+' is illegal
	
