/*
 *
 * FILE: testrun_sum_std.go
 * LATEST: 10:19 05 May 2025
 * DESC: sum values from iterative Kyber batch jobs.
 * AUTHOR: Levi Neuwirth <ln@levineuwirth.org>
 *
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var count float64
var testSums map[string]float64
var lastTest string
var gen_a []float64
var indcpa_keypair []float64
var indcpa_enc []float64
var keypair_derand []float64
var keypair []float64
var encaps []float64
var decaps []float64

func main() {
	if len(os.Args) < 1 {
		log.Fatal("Usage: ./testrun_sum_std <path to slurm.OUT file>")
	}

	outRaw, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully opened slurm STDOUT")
	}
	defer outRaw.Close()
	initTestSums()
	count = 0
	lastTest = "none"
	scanner := bufio.NewScanner(outRaw)
	for scanner.Scan() {
		localLine := scanner.Text()
		// Check if we've hit a new test iteration
		if strings.Contains(localLine, "Loop spin:") {
			count += 1
			continue
			// Otherwise, we might have data from a previously indicated test.
		} else if strings.Contains(localLine, "average:") {
			// We split the line and add to the appropriate testSums index.
			line := localLine[9:]
			var numberStr strings.Builder
			for _, ch := range line {
				if (ch >= '0' && ch <= '9') || ch == '.' {
					numberStr.WriteRune(ch)
				} else {
					break
				}
			}

			add, err := strconv.ParseFloat(numberStr.String(), 64)
			if err != nil {
				log.Printf("Failed to parse number from line %q: %v", localLine, err)
				continue
			}
			testSums[lastTest] += add
			// And now for the stddev:
			switch lastTest {
			case "gen_a:":
				gen_a = append(gen_a, add)
			case "indcpa_keypair:":
				indcpa_keypair = append(indcpa_keypair, add)
			case "indcpa_enc:":
				indcpa_enc = append(indcpa_enc, add)
			case "kyber_keypair_derand:":
				keypair_derand = append(keypair_derand, add)
			case "kyber_keypair:":
				keypair = append(keypair, add)
			case "kyber_encaps:":
				encaps = append(encaps, add)
			case "kyber_decaps:":
				decaps = append(decaps, add)
			default:
				continue
			}
			continue
			// We aren't concerned with the medians here.
		} else if strings.Contains(localLine, "median:") {
			continue
		}

		// Here, figure out what the test was for the next data.
		trimmed := strings.TrimSpace(localLine)
		if strings.HasSuffix(trimmed, ":") && !strings.Contains(trimmed, "average") && !strings.Contains(trimmed, "median") {
			lastTest = trimmed
			continue
		}

	}

	// Now we take the averages and stddevs.
	fmt.Printf("gen_a avg: %f\ngen_a stddev: %f\n", testSums["gen_a:"]/count, calcStddev("gen_a:", gen_a))
	fmt.Printf("indcpa keypair avg: %f\nindcpa_keypair stddev: %f\n", testSums["indcpa_keypair:"]/count, calcStddev("indcpa_keypair:", indcpa_keypair))
	fmt.Printf("indcpa enc avg: %f\nindcpa_enc stddev: %f\n", testSums["indcpa_enc:"]/count, calcStddev("indcpa_enc:", indcpa_enc))
	fmt.Printf("keypair_derand avg: %f\nkeypair_derand stddev:: %f\n", testSums["kyber_keypair_derand:"]/count, calcStddev("kyber_keypair_derand:", keypair_derand))
	fmt.Printf("keypair avg: %f\nkeypair stddev:: %f\n", testSums["kyber_keypair:"]/count, calcStddev("kyber_keypair:", keypair))
	fmt.Printf("encaps avg: %f\nencaps stddev:: %f\n", testSums["kyber_encaps:"]/count, calcStddev("kyber_encaps:", encaps))
	fmt.Printf("decaps avg: %f\ndecaps stddev:: %f\n", testSums["kyber_decaps:"]/count, calcStddev("kyber_decaps:", decaps))
}

func initTestSums() {
	testSums = make(map[string]float64)
	testSums["gen_a:"] = 0
	testSums["indcpa_keypair:"] = 0
	testSums["indcpa_enc:"] = 0
	testSums["kyber_keypair_derand:"] = 0
	testSums["kyber_keypair:"] = 0
	testSums["kyber_encaps:"] = 0
	testSums["kyber_decaps:"] = 0
}

func calcStddev(test string, inputs []float64) (result float64) {
	mean := float64(testSums[test] / float64(len(inputs)))
	var variance float64
	for _, value := range inputs {
		variance += (value - mean) * (value - mean)
	}

	return math.Sqrt(variance / float64(len(inputs)))
}
