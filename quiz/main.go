package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type score struct {
	correct  int
	mistaken int
}

type question struct {
	problem  string
	solution string
}

func loading(filename string) []question {
	data, _ := os.Open(filename + ".csv")
	csvFile := csv.NewReader(data)

	list := []question{}

	for {
		q, err := csvFile.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		list = append(list, question{q[0], q[1]})
	}

	return list
}

func readString(c chan string) {
	r := bufio.NewReader(os.Stdin)
	rstr, _ := r.ReadString('\n')

	c <- strings.TrimSpace(rstr)
	close(c)
}

func checkSolution(s string, i string, total *score) {
	if s == i {
		total.correct++
		fmt.Printf("Correct!\n")
	} else {
		total.mistaken++
		fmt.Printf("False!\n")
	}
}

func startQuiz(list []question, duration int) {
	total := score{0, 0}

	timer := time.NewTimer(time.Duration(duration) * time.Second)

	for _, q := range list {
		fmt.Printf("\nWhat %s equals, sir?\n", q.problem)

		rsChan := make(chan string)

		go readString(rsChan)

		select {
		case <-timer.C:
			fmt.Printf("You gave %v correct answers and %v false\n", total.correct, total.mistaken)
			return
		case input := <-rsChan:
			checkSolution(q.solution, input, &total)
		}
	}

	fmt.Printf("Total score: %v correct and %v false\n", total.correct, total.mistaken)
}

func main() {
	filename := flag.String("csv", "problems", "filename of csv")
	duration := flag.Int("dur", 10, "Duration of timer")
	flag.Parse()

	data := loading(*filename)

	startQuiz(data, *duration)
}
