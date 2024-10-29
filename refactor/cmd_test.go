package cmd

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSleep1s(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./testcode1/sleep1s.go")
	if err != nil {
		t.Fatal(err)
	}

	errored := true
	go func() {
		for data := range cmd.Output() {
			if !data.IsStderr {
				if data.EOF {
					errored = false
					return
				}
			}
		}
	}()
	exitCode, err := cmd.Run()
	if err != nil {
		t.Fatal(1)
	}
	if exitCode != 0 {
		t.Fatal(1)
	}
	if errored {
		t.Fatal(1)
	}
}

func TestPrintln(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./testcode2/println.go")
	if err != nil {
		t.Fatal(err)
	}

	msgs := []Data{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for msg := range cmd.Output() {
			msgs = append(msgs, msg)
		}
		wg.Done()
	}()
	exitCode, err := cmd.Run()
	if err != nil {
		t.Fatal(1)
	}
	if exitCode != 0 {
		t.Fatal(1)
	}
	wg.Wait()

	assert.Equal(t, len(msgs), 3)
	assert.Equal(t, msgs[0].Value, "Hello World!")
	assert.False(t, msgs[0].IsStderr)
	assert.Equal(t, msgs[1].Value, "Hello Moon!")
	assert.False(t, msgs[1].IsStderr)
	assert.True(t, msgs[2].EOF)
}

func TestExit1(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./testcode3/exit1.go")
	if err != nil {
		t.Fatal(err)
	}

	errored := false
	go func() {
		for range cmd.Output() {
			errored = true
		}
	}()
	exitCode, err := cmd.Run()
	if err != nil {
		t.Fatal(1)
	}
	if exitCode != 1 {
		t.Fatal(1)
	}
	if errored {
		t.Fatal(1)
	}
}

// func TestStderr(t *testing.T) {
// 	cmd, err := NewCmd("go", "run", "./testcode4/stderr.go")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	errored := false
// 	go func() {
// 		for range cmd.Output() {
// 			errored = true
// 		}
// 		if <-cmd.Output() != "Hello Stderr!" {
// 			errored = true
// 			return
// 		}
// 		if <-cmd.Stderr() == "Hello errors!" {
// 			errored = false
// 		}
// 	}()
// 	exitCode, err := cmd.Run()
// 	if err != nil {
// 		t.Fatal(1)
// 	}
// 	if exitCode != 0 {
// 		t.Fatal(1)
// 	}
// 	if errored {
// 		t.Fatal(1)
// 	}
// }

func TestStdin(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./testcode5/stdin.go")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		cmd.Write("example\n")
	}()
	exitCode, err := cmd.Run()
	if err != nil {
		t.Fatal(1)
	}
	if exitCode != 0 {
		t.Fatal(1)
	}
}

func TestStdinFail(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./testcode5/stdin.go")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		cmd.Write("bogus\n")
	}()
	exitCode, err := cmd.Run()
	if err != nil {
		t.Fatal(1)
	}
	if exitCode != 1 {
		t.Fatal(1)
	}
}
