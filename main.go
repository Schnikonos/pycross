package main

import (
	"pycross/baseIo"
	"pycross/pycrossBasic"
)

func main() {
	id := 1
	rows, cols := baseIo.GetData(id)
	pycross := pycrossBasic.Init(rows, cols)
	print(pycross)
	res := pycrossBasic.Solve(pycross)
	//res := []string{"*   *", "**  *", "**  ", "*** ", " *** "}
	baseIo.CheckRes(id, res)
}
