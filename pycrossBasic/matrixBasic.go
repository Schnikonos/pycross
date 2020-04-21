package pycrossBasic

type Picross struct {
	rowDef  [][]int
	colDef  [][]int
	rowLen  int
	colLen  int
	matrix  [][]int
	rowComb [][][]int
	colComb [][][]int
}

func sum(a []int) int {
	res := 0
	for _, el := range a {
		res += el
	}
	return res
}

func array(el int, length int) []int {
	res := make([]int, length)
	for i := range res {
		res[i] = el
	}
	return res
}

func concat(slices [][]int) []int {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]int, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func combine(res *[][]int, lineLen int, remainingDef []int, workLine []int, separator []int) {
	remainingLen := lineLen - len(workLine)
	if len(remainingDef) == 0 {
		*res = append(*res, append(workLine, array(0, remainingLen)...))
		return
	}

	block, remainingDef := remainingDef[0], remainingDef[1:]
	moveRange := remainingLen - block - len(separator) - sum(remainingDef) - len(remainingDef) + 1

	for i := 0; i < moveRange; i++ {
		combine(res, lineLen, remainingDef,
			concat([][]int{workLine, separator, array(0, i), array(1, block)}),
			[]int{0})
	}
	remainingDef = append(remainingDef, block)
}

func getLineCombs(lineLength int, lineDef []int) *[][]int {
	result := make([][]int, 0)
	combine(&result, lineLength, lineDef, []int{}, []int{})
	return &result
}

func Init(rowDef, colDef [][]int) *Picross {
	picross := &Picross{rowDef: rowDef, colDef: colDef, rowLen: len(rowDef),
		colLen: len(colDef)}

	picross.matrix = make([][]int, picross.rowLen)
	for i := 0; i < picross.rowLen; i++ {
		picross.matrix[i] = make([]int, picross.colLen)
		for j := 0; j < picross.colLen; j++ {
			picross.matrix[i][j] = -1
		}
	}

	picross.rowComb = make([][][]int, picross.rowLen)
	picross.colComb = make([][][]int, picross.colLen)

	for i := 0; i < picross.rowLen; i++ {
		picross.rowComb[i] = *getLineCombs(picross.colLen, picross.rowDef[i])
	}

	for j := 0; j < picross.colLen; j++ {
		picross.colComb[j] = *getLineCombs(picross.rowLen, picross.colDef[j])
	}

	return picross
}

func checkLine(currentLine *[]int, pattern *[]int) bool {
	for i, el := range *currentLine {
		if el != -1 && (*pattern)[i] != el {
			return false
		}
	}
	return true
}

func filter(combs *[][]int, line *[]int, index int, channel chan<- *combChannel) {
	result := [][]int{}
	for i := range *combs {
		if checkLine(line, &(*combs)[i]) {
			result = append(result, (*combs)[i])
		}
	}

	channel <- &combChannel{
		index:  index,
		result: &result,
	}
}

type combChannel struct {
	index  int
	result *[][]int
}

func filterComb(linesComb *[][][]int, matrix *[][]int, getLine func(int, *[][]int) *[]int) bool {
	oldComb := 0
	newComb := 0

	channel := make(chan *combChannel)

	for i, lineComb := range *linesComb {
		oldComb += len(lineComb)
		line := getLine(i, matrix)
		go filter(&(*linesComb)[i], line, i, channel)
	}

	for range *linesComb {
		channelComb := <-channel
		(*linesComb)[channelComb.index] = *channelComb.result
		newComb += len((*linesComb)[channelComb.index])
	}

	return oldComb != newComb
}

func getGlobalState(index int, combinations *[][]int) int {
	c := (*combinations)[0][index]
	for _, el := range *combinations {
		if el[index] != c {
			return -1
		}
	}
	return c
}

func threadCompute(i, j int, p *Picross, channel chan<- bool) {
	if p.matrix[i][j] == -1 {
		p.matrix[i][j] = getGlobalState(j, &p.rowComb[i])
	}
	if p.matrix[i][j] == -1 {
		p.matrix[i][j] = getGlobalState(i, &p.colComb[j])
	}
	channel <- true
}

func (p *Picross) reduce() bool {
	channel := make(chan bool)
	for i := 0; i < p.rowLen; i++ {
		for j := 0; j < p.colLen; j++ {
			go threadCompute(i, j, p, channel)
		}
	}

	for i := 0; i < p.rowLen*p.colLen; i++ {
		<-channel
	}

	rowsChange := filterComb(&p.rowComb, &p.matrix, func(i int, matrix *[][]int) *[]int { return &(*matrix)[i] })
	colsChange := filterComb(&p.colComb, &p.matrix, func(i int, matrix *[][]int) *[]int {
		res := make([]int, len(*matrix))
		for j, el := range *matrix {
			res[j] = el[i]
		}
		return &res
	})

	return rowsChange || colsChange
}

func (p *Picross) toString() []string {
	res := make([]string, len(p.matrix))
	for i, line := range p.matrix {
		res[i] = lineToString(line)
	}
	return res
}

func lineToString(line []int) string {
	res := ""
	for _, el := range line {
		switch el {
		case 0:
			res += " "
		case 1:
			res += "*"
		default:
			res += "?"
		}
	}
	return res
}

func Solve(picross *Picross) []string {
	for picross.reduce() {
		// pass
	}
	return picross.toString()
}
