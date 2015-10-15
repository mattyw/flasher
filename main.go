package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/juju"
	"github.com/juju/juju/juju/osenv"
	"github.com/schoentoon/piglow"
)

var (
	intensity = flag.Int("i", 1, "intensity of the lights")
	envName   = flag.String("env", "amazon", "juju environment to monitor")
	jujuHome  = flag.String("jhome", "/home/mattyw/.juju", "the home of juju")
)

func status(envName string, ch chan bool) error {
	osenv.SetJujuHome(*jujuHome)
	conn, err := juju.NewAPIFromName(envName)
	if err != nil {
		return err
	}
	client := conn.Client()
	status, err := client.Status([]string{})
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", status)
	for _, s := range status.Services {
		for _, u := range s.Units {
			if u.UnitAgent.Status == params.StatusError ||
				u.Workload.Status == params.StatusError ||
				u.AgentState == params.StatusError {
				ch <- false
				return nil
			}
		}
	}
	ch <- true
	return nil
}

func errorFlash() {
	for i := 0; i < 3; i++ {
		err := piglow.Leg(byte(i), byte(*intensity))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
		err = piglow.ShutDown()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func okFlash() {
	for i := 0; i < 3; i++ {
		err := piglow.Leg(byte(i), byte(*intensity))
		if err != nil {
			fmt.Println(err)
		}
	}
	time.Sleep(100 * time.Millisecond)
}

func do(ch chan bool) error {
	if err := piglow.HasPiGlow(); err != nil {
		fmt.Println("no piglow")
		fmt.Println(err)
		return err
	}
	var ok bool
	go func() {
		for {
			ok = <-ch
			time.Sleep(100 * time.Millisecond)
		}
	}()
	for {
		if ok {
			fmt.Println("ok")
			okFlash()
		} else {
			fmt.Println("error")
			errorFlash()
		}
		piglow.ShutDown()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	flag.Parse()
	fmt.Println("Doing")
	ch := make(chan bool)
	go func() {
		for {
			err := status(*envName, ch)
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(2 * time.Second)
		}
	}()
	err := do(ch)
	if err != nil {
		fmt.Println(err)
	}
}
