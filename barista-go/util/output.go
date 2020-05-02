package util

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func OutputError(err error) {
	fmt.Println(Sprintf("%s %s\n%s", Red(Bold("Error")), Bold("==>"), err.Error()))
}
