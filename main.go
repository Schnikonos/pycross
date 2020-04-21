package main

import (
	"pycross/baseIo"
	"pycross/pycrossBasic"
	"pycross/pycrossObject"
)

func main() {
	for i := 1; i < 7; i++ {
		execute1(i)
		execute2(i)
	}
}

func execute1(i int) {
	rows, cols := baseIo.GetData(i)
	pycross := pycrossBasic.Init(rows, cols)
	res := pycrossBasic.Solve(pycross)
	baseIo.CheckRes("go_basic", i, res)
}

func execute2(i int) {
	rows, cols := baseIo.GetData(i)
	pycross := pycrossObject.Init(rows, cols)
	res := pycrossObject.Solve(pycross)
	baseIo.CheckRes("go_object", i, res)
}
