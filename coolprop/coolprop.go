package main

import (
	"fmt"

	"CoolProp"
)

func main() {
	waterTCrit := CoolProp.Props1SI("Pcrit", "Water")
	fmt.Printf("Water TCrit : %v degC \n", waterTCrit-273.16)
	H1 := CoolProp.PropsSI("H", "T", 275, "P", 101325, "Water")
	H2 := CoolProp.PropsSI("H", "T", 280, "P", 101325, "Water")
	fmt.Printf("Cambio de entalpía es de: %v J/kg\n", H2-H1)
}
