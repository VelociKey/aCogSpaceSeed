module sov.fleet/plinth-filesystem

go 1.26.4

require (
	sov.fleet/logiclibrary v0.0.0
	sov.fleet/white3 v0.0.0
)

replace (
	sov.fleet/logiclibrary => ../logiclibrary
	sov.fleet/white3 => ../../00xper/white3
)
