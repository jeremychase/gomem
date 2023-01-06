package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
)

var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var iterations = flag.Uint("iterations", 1, "number of iterations to run")

var storage = map[int]int{}

type data struct {
	Key   int
	Value int
}

func main() {
	rand.Seed(42)
	parseFlags()

	for *iterations > 0 {
		*iterations--

		err := leak()
		if err != nil {
			log.Fatal(err)
		}

		if *iterations == 0 {
			// Last run
		}
	}

}

func leak() error {
	t := data{rand.Int(), rand.Int()}

	// Marshal/Unmarshal into new value for spurious allocations
	var r data
	{
		js, err := json.Marshal(t)
		if err != nil {
			log.Println("marshal failure")
			return err
		}

		if err := json.Unmarshal(js, &r); err != nil {
			log.Println("unmarshal failure")
			return err
		}
	}

	// 'leak'
	storage[r.Key] = r.Value

	return nil
}

func parseFlags() {
	flag.Parse()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {

			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
	}
}
