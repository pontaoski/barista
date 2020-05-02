package util

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func OutputError(err error) {
	fmt.Println(Sprintf("%s %s %s"), Red(Bold("Error")), Bold("==>"), err.Error())
}
