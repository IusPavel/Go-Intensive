// Program reads from standard input N numbers and prints metrics.
// N is first argument. Metrics are 'mean', 'median', 'mode', 'standard deviation'.
package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	sum := checkArguments()
	fmt.Printf("Enter %d numbers beetween -100000 and 100000\n", sum)
	seq := readSequence(sum)
	metrics := calculating(seq)
	printPreResult(metrics)
}

// readSequence reads 'sum' numbers from standard input and stores in int-slice.
func readSequence(sum int) []int {
	seq := make([]int, 0, sum)
	reader := bufio.NewScanner(os.Stdin)
	for sum > 0 {
		fmt.Printf("Elements left: %d. Enter next number: ", sum)
		reader.Scan()
		a := reader.Text()
		if len(a) == 0 {
			log.Print("Empty string. Not accepted\n")
			continue
		}
		ai, err := strconv.ParseInt(a, 10, 32)
		if err != nil {
			log.Printf("invalid argument %s. Enter again\n", a)
			continue
		}
		if ai < -100000 {
			fmt.Printf("The number %d is less than -100000. Enter again\n", ai)
			continue
		}
		if ai > 100000 {
			fmt.Printf("The number %d is more than 100000. Enter again\n", ai)
			continue
		}
		seq = append(seq, int(ai))
		sum--
	}
	return seq
}

// calculating sorts sequence, computes values and saves result in map[string]float64.
func calculating(seq []int) map[string]float64 {
	sort.Ints(seq)
	m := make(map[string]float64, 4)
	m["mean"] = func() float64 {
		sum := 0
		for _, v := range seq {
			sum += v
		}
		return float64(sum) / float64(len(seq))
	}()
	m["median"] = func() float64 { // 1 2 -> 1.5 // 1 2 3 -> 2
		if len(seq)%2 != 0 {
			return float64(seq[len(seq)/2])
		}
		return float64(seq[len(seq)/2-1]+seq[len(seq)/2]) / 2
	}()
	m["mode"] = func() float64 {
		if len(seq) == 1 {
			return 1.0
		}
		oldMode := seq[0]
		oldLen := 1
		newMode := seq[0]
		newLen := 1
		for _, v := range seq[1:] {
			if oldMode == v {
				oldLen++
			} else if newMode == v {
				newLen++
				if newLen > oldLen {
					oldMode = newMode
					oldLen = newLen
				}
			} else {
				newMode = v
				newLen = 1
			}
		}
		return float64(oldMode)
	}()
	m["sd"] = func() float64 {
		sd := 0.0
		for _, v := range seq {
			sd += math.Pow(float64(v)-m["mean"], 2)
		}
		sd = math.Sqrt(sd / float64(len(seq)-1))
		return sd
	}()
	return m
}

// printPreResult decides what will be printed and prints it.
func printPreResult(metrics map[string]float64) {
	argv := os.Args
	if len(argv) > 2 {
		argv = argv[2:]
		for _, v := range argv {
			fmt.Printf("%s: %.4f\n", v, metrics[v])
		}
	} else {
		for k, v := range metrics {
			fmt.Printf("%s: %.4f\n", k, v)
		}
	}

}

// checkArguments checks entered arguments. 1st - len of sequence, others are specified metrics.
func checkArguments() int {
	argv := os.Args
	if len(argv) < 2 {
		log.Fatalln("Specify quantity of numbers")
	}
	n, err := strconv.ParseInt(argv[1], 10, 64)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	if n < 1 {
		log.Fatalln("Quantity of numbers can't be less than 1")
	}
	if len(argv) > 2 {
		for _, v := range argv[2:] {
			if strings.ToLower(v) != "mean" &&
				strings.ToLower(v) != "median" &&
				strings.ToLower(v) != "mode" &&
				strings.ToLower(v) != "sd" {
				log.Fatalln("Incorrect argument:", v, "\nCorrected are: \"mean\", \"median\", \"mode\", \"sd\"")
			}
		}
	}
	return int(n)
}
