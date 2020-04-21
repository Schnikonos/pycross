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

var timeStart time.Time

func GetData(i int) (rows [][]int, cols [][]int) {
	jsonFile, _ := os.Open(fmt.Sprintf("testFiles/test%v_in.json", i))
	defer jsonFile.Close()
	fmt.Print(jsonFile)
	var data rawData
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &data)

	timeStart = time.Now()
	return data[0], data[1]
}

func CheckRes(i int, res []string) {
	ellapse := time.Since(timeStart)
	fmt.Printf("******* Ellapse: %v ********\n\n", ellapse)

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
}
