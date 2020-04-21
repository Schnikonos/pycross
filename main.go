package main

import (
	"pycross/baseIo"
	"pycross/pycrossBasic"
)

func main() {
	for i := 1 ; i < 6 ; i++ {
		execute(i)
	}
}

func execute(i int) {
	rows, cols := baseIo.GetData(i)
	pycross := pycrossBasic.Init(rows, cols)
	res := pycrossBasic.Solve(pycross)
	baseIo.CheckRes(i, res)
}