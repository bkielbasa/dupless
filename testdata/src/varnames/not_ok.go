package varnames

var xyz = 123 // want "the variable name contains the forbidden pattern"

func myFunc2() {
	xxx := "" // want "the variable name contains the forbidden pattern"
	if xxx == "" {
	}
}
