package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
)

var iterations = flag.Uint("iterations", 1, "number of iterations to run")
var mprofPrefix = flag.String("mprofPrefix", "output", "prefix for mprof files")

var storage = map[int]int{}

func main() {
	rand.Seed(42)
	mid, end := parseFlags()
	defer mid.Close()
	defer end.Close()

	i := uint(0)
	for i < *iterations {
		i++

		err := leak()
		if err != nil {
			log.Fatal(err)
		}

		// Write mid point mprof
		if i*2 == *iterations {
			fmt.Printf("writing %s at %d\n", mid.Name(), i)
			err := pprof.WriteHeapProfile(mid)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Write final mprof
		if i == *iterations {
			fmt.Printf("writing %s at %d\n", end.Name(), i)
			err := pprof.WriteHeapProfile(end)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

// leak puts random values into a map while doing unnecessary work.
func leak() error {
	// data exists for marshalling purposes
	type data struct {
		Key   int
		Value int
	}

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

func parseFlags() (*os.File, *os.File) {
	flag.Parse()

	if *iterations%2 != 0 {
		log.Println("odd iteration count; adding one to allow mid point")
		*iterations++
	}

	mid, err := os.Create(*mprofPrefix + "-mid.mprof")
	if err != nil {
		log.Fatal(err)
	}
	end, err := os.Create(*mprofPrefix + "-end.mprof")
	if err != nil {
		log.Fatal(err)
	}

	return mid, end
}
