package pycrossObject

type state int

const (
	na     state = -1
	empty        = 0
	filled       = 1
)

type cell struct {
	state state
}

func (c *cell) String() string {
	switch c.state {
	case 0:
		return " "
	case 1:
		return "*"
	default:
		return "?"
	}
}

func initCell() *cell {
	return &cell{state: na}
}

func newCell(state state) *cell {
	return &cell{state: state}
}

func sum(a []int) int {
	res := 0
	for _, el := range a {
		res += el
	}
	return res
}

func array(el state, length int) []state {
	res := make([]state, length)
	for i := range res {
		res[i] = el
	}
	return res
}

func concat(slices [][]state) []state {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]state, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func combine(res *[][]state, lineLen int, remainingDef []int, workLine []state, separator []state) {
	remainingLen := lineLen - len(workLine)
	if len(remainingDef) == 0 {
		*res = append(*res, append(workLine, array(empty, remainingLen)...))
		return
	}

	block, remainingDef := remainingDef[0], remainingDef[1:]
	moveRange := remainingLen - block - len(separator) - sum(remainingDef) - len(remainingDef) + 1

	for i := 0; i < moveRange; i++ {
		combine(res, lineLen, remainingDef,
			concat([][]state{workLine, separator, array(empty, i), array(filled, block)}),
			[]state{empty})
	}
	remainingDef = append(remainingDef, block)
}

func getLineCombs(lineLength int, lineDef []int) *[][]state {
	result := make([][]state, 0)
	combine(&result, lineLength, lineDef, []state{}, []state{})
	return &result
}

type line struct {
	definition   []int
	combinations [][]state
}

func newLine(lineLength int, definition []int) *line {
	res := line{definition: definition}
	res.combinations = *getLineCombs(lineLength, definition)
	return &res
}

type lines struct {
	length int
	defs   []*line
}

func newLines(length int, linesDef *[][]int) *lines {
	res := lines{length: len(*linesDef)}
	res.defs = make([]*line, len(*linesDef))

	for i, lineDef := range *linesDef {
		res.defs[i] = newLine(length, lineDef)
	}
	return &res
}

type combChannel struct {
	index  int
	result *[][]state
}

func checkLine(currentLine *[]*cell, pattern *[]state) bool {
	for i, el := range *currentLine {
		if el.state != -1 && (*pattern)[i] != el.state {
			return false
		}
	}
	return true
}

func filter(combs *[][]state, line *[]*cell, index int, channel chan<- *combChannel) {
	result := [][]state{}
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

func (l *lines) filterComb(matrix *[][]*cell, getLine func(int, *[][]*cell) *[]*cell) bool {
	oldComb := 0
	newComb := 0

	channel := make(chan *combChannel)

	for i, line := range l.defs {
		oldComb += len(line.combinations)
		matrixLine := getLine(i, matrix)
		go filter(&line.combinations, matrixLine, i, channel)
	}

	for range l.defs {
		channelComb := <-channel
		l.defs[channelComb.index].combinations = *channelComb.result
		newComb += len(*channelComb.result)
	}

	return oldComb != newComb
}

type Picross struct {
	rows   *lines
	cols   *lines
	matrix [][]*cell
}

func Init(rowDef, colDef [][]int) *Picross {
	rowLength := len(rowDef)
	colLength := len(colDef)
	picross := &Picross{rows: newLines(colLength, &rowDef), cols: newLines(rowLength, &colDef)}

	picross.matrix = make([][]*cell, rowLength)
	for i := 0; i < rowLength; i++ {
		picross.matrix[i] = make([]*cell, colLength)
		for j := 0; j < colLength; j++ {
			picross.matrix[i][j] = initCell()
		}
	}

	return picross
}

func getGlobalState(index int, combinations *[][]state) state {
	c := (*combinations)[0][index]
	for _, el := range *combinations {
		if el[index] != c {
			return na
		}
	}
	return c
}

func threadCompute(i, j int, p *Picross, channel chan<- bool) {
	if p.matrix[i][j].state == na {
		p.matrix[i][j].state = getGlobalState(j, &p.rows.defs[i].combinations)
	}
	if p.matrix[i][j].state == na {
		p.matrix[i][j].state = getGlobalState(i, &p.cols.defs[j].combinations)
	}
	channel <- true
}

func (p *Picross) reduce() bool {
	channel := make(chan bool)
	for i := 0; i < p.rows.length; i++ {
		for j := 0; j < p.cols.length; j++ {
			go threadCompute(i, j, p, channel)
		}
	}

	for i := 0; i < p.rows.length*p.cols.length; i++ {
		<-channel
	}

	rowsChange := p.rows.filterComb(&p.matrix, func(i int, matrix *[][]*cell) *[]*cell { return &(*matrix)[i] })
	colsChange := p.cols.filterComb(&p.matrix, func(i int, matrix *[][]*cell) *[]*cell {
		res := make([]*cell, len(*matrix))
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

func lineToString(line []*cell) string {
	res := ""
	for _, el := range line {
		switch el.state {
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
