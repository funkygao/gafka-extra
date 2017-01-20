package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	fileName = flag.String("fname", "/tmp/big.1gb", "file name for testing.")
)

func main() {
	flag.Parse()
	go func() {
		f, err := os.OpenFile(*fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		for i := 0; i < 1000; i++ {
			//time.Sleep(time.Millisecond * 50)
			f.WriteString(strings.Repeat("X", 10<<10) + "\n")
		}
	}()

	go func() {
		f, err := os.OpenFile(*fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		for i := 0; i < 1000; i++ {
			//time.Sleep(time.Millisecond * 50)
			f.WriteString(strings.Repeat("Y", 12<<10) + "\n")
		}
	}()

	time.Sleep(time.Second * 2)

	b, err := ioutil.ReadFile(*fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		for _, xOrY := range l {
			if xOrY != rune(l[0]) {
				println(l)
				break
			}
		}

	}

}
