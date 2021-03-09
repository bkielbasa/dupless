package funcs

func dupa() {} // want "the function name contains the forbidden pattern"

func functionThatContainsDupa() {} // want "the function name contains the forbidden pattern"

func correctFuncName() {}
