package main

import (
	"fmt"

	"units/quantity"
)

func p1() {
	fmt.Printf("-------------------Prueba 1--------------------------\n")
	g, e := quantity.Parse("9.8 N.cd")
	fmt.Printf("%v, %v\n", g, e)
}

func p2() {
	// Vamos a probar
	fmt.Printf("-------------------Prueba 2--------------------------\n")
	a, e := quantity.Parse("1 m")
	b, e := quantity.Parse("4 in")
	fmt.Printf("Están las dos magnitudes: %v y %v\n", a, b)
	fmt.Printf("El símbolo de la unidad b es: %v\n", b.Symbol())
	sePuede := a.HasCompatibleUnit(b.Symbol())
	fmt.Printf("Se puede convertir? %v\n", sePuede)
	fmt.Printf("%v, %v\n", a, e)
}

func p3() {
	// Vamos a probar errores con el parseo de las unidades
	fmt.Printf("-------------------Prueba 3--------------------------\n")
	a, e := quantity.Parse("5 kkm")
	fmt.Printf("La unidad y el error son: %v y %v\n", a, e)
}

func p4() {
	// Vamos a probar con unidades de temperatura
	fmt.Printf("-------------------Prueba 4--------------------------\n")
	temp, e := quantity.Parse("100 degC")
	fmt.Printf("A ver, la temperatura es %v y el error %v\n", temp, e)
	tempK, success := temp.ConvertTo("K")
	fmt.Printf("La conversión fue de %v y el éxito es %v\n", tempK, success)
	fmt.Print(temp.ToSI())
}

func p5() {
	// Vamos a probar dimensionalidad
	fmt.Printf("-------------------Prueba 5--------------------------\n")
	cantidad, e := quantity.Parse("100 yd.degF")
	fmt.Printf("La cantidad %v y el error %v\n", cantidad, e)
	dimensiones := cantidad.Dimensionality()
	fmt.Printf("%v", dimensiones)
}

func main() {
	p1()
	p2()
	p3()
	p4()
	p5()
}
