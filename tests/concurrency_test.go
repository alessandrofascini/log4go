package tests

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFatalConcurrency(t *testing.T) {
	fmt.Println(os.Getpid())
	go Hypervisor()
}

func Hypervisor() {
	fmt.Println("hypervisor pid: ", os.Getpid())
	fmt.Println("parent pid: ", os.Getppid())
}

func TestCountTo(t *testing.T) {
	for i := range countTo(10) {
		fmt.Println(i)
	}
}

func countTo(max int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < max; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

type AppenderChannel chan struct{}

func TestForSelectLoop(t *testing.T) {
	done := make(chan AppenderChannel)
	select {
	case v := <-done:
		fmt.Println("read from ch:", v)
	default:
		fmt.Println("no value written to ch")
	}
	close(done)
}

func TestChannelOnPanic(t *testing.T) {
	go func() {
		//process, err := os.FindProcess(os.Getpid())
		//if err != nil {
		//	panic(fmt.Errorf("%w", err))
		//}
		//
		<-time.After(5 * time.Second)
	}()
	defer func() {
		e := recover()
		switch v := e.(type) {
		case error:
			pe := errors.Unwrap(v)
			fmt.Println("unwrap", pe)
			ppe := errors.Unwrap(pe)
			fmt.Println("un-unwrap", ppe)
		case string:
			fmt.Println("as string", v)
		default:
			fmt.Println(e)
		}
	}()
	panic(errors.New("new error"))
}
