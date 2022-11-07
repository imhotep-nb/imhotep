package utils

import (
	"fmt"
)

func HandleLog(logger *[]string, format string, v ...interface{}) {
	aGuardar := fmt.Sprintf(format, v...)
	*logger = append(*logger, aGuardar)
	fmt.Print(aGuardar + "\n")
}
