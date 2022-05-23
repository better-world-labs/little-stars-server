package main

import "aed-api-server/internal/pkg/inject"

func main() {
	// load components
	new(inject.Component).
		Load(&A{}).
		Load(&B{}).
		Load(&ICImpl1{}, "impl1").
		Load(&ICImpl2{}, "impl2").
		Install()
}
