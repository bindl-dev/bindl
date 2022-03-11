package internal

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var errHeader = color.HiRedString("ERROR")

func ErrorMsg(err error) {
	fmt.Fprintf(os.Stderr, "%s - %s\n", errHeader, err.Error())
}
