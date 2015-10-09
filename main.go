package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/schoentoon/piglow"
)

var (
	intensity = flag.Int("i", 1, "intensity of the lights")
)

func main() {
	flag.Parse()
	fmt.Println("Doing")

	if err := piglow.HasPiGlow(); err != nil {
		fmt.Println("no piglow")
		fmt.Println(err)
		return
	}
	for {
		for i := 0; i < 3; i++ {
			err := piglow.Leg(byte(i), byte(*intensity))
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(500 * time.Millisecond)
			err = piglow.ShutDown()
			if err != nil {
				fmt.Println(err)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
