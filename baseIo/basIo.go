package baseIo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type rawData [][][]int

type stat map[string]float64

var timeStart time.Time

func GetData(i int) (rows [][]int, cols [][]int) {
	jsonFile, _ := os.Open(fmt.Sprintf("testFiles/test%v_in.json", i))
	defer jsonFile.Close()
	var data rawData
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &data)

	timeStart = time.Now()
	return data[0], data[1]
}

func CheckRes(name string, i int, res []string) {
	ellapse := time.Since(timeStart)

	for _, line := range res {
		fmt.Print(line + "\n")
	}

	expectedFile, _ := os.Open(fmt.Sprintf("testFiles/test%v_out.txt", i))
	defer expectedFile.Close()

	scanner := bufio.NewReader(expectedFile)
	for j := range res {
		expectedLine, _ := scanner.ReadString('\n')
		if strings.TrimRight(res[j], " ") != strings.TrimRight(expectedLine, " \r\n") {
			log.Fatal("Different results!")
		}
	}

	statFile, _ := os.Open(fmt.Sprintf("testFiles/test%v_stat.json", i))
	var data stat
	byteValue, _ := ioutil.ReadAll(statFile)
	statFile.Close()

	if len(byteValue) > 0 {
		_ = json.Unmarshal(byteValue, &data)
	} else {
		_, _ = os.Create(fmt.Sprintf("testFiles/test%v_stat.json", i))
		data = make(stat)
	}

	statFile, _ = os.Create(fmt.Sprintf("testFiles/test%v_stat.json", i))
	data[name] = ellapse.Seconds()
	result, _ := json.MarshalIndent(data, "", "  ")
	_, _ = statFile.Write(result)

	fmt.Printf("-------------- OK - Ellapse: %v -----------\n\n", ellapse)
}
