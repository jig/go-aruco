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

	msgs := run(t, cmd, 0)

	assert.Equal(t, 1, len(msgs))
	assert.True(t, msgs[0].EOF)
}

func TestPrintln(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./testcode2/println.go")
	if err != nil {
		t.Fatal(err)
	}

	msgs := run(t, cmd, 0)

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

	msgs := run(t, cmd, 1)

	assert.Equal(t, len(msgs), 2)
	assert.Equal(t, msgs[0].Value, "exit status 1")
	assert.True(t, msgs[0].IsStderr)
	assert.True(t, msgs[1].EOF)
}

func TestStderr(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./testcode4/stderr.go")
	if err != nil {
		t.Fatal(err)
	}

	msgs := run(t, cmd, 0)

	assert.Equal(t, len(msgs), 3)
	assert.Equal(t, msgs[0].Value, "Hello Stderr!")
	assert.True(t, msgs[0].IsStderr)
	assert.Equal(t, msgs[1].Value, "Hello errors!")
	assert.True(t, msgs[1].IsStderr)
	assert.True(t, msgs[2].EOF)
}

func TestStdin(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./testcode5/stdin.go")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		cmd.Write("example\n")
	}()

	msgs := run(t, cmd, 0)

	assert.Equal(t, 1, len(msgs))
	assert.True(t, msgs[0].EOF)
}

func TestStdinFail(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./testcode5/stdin.go")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		cmd.Write("bogus\n")
	}()

	msgs := run(t, cmd, 1)

	assert.Equal(t, len(msgs), 2)
	assert.Equal(t, msgs[0].Value, "exit status 1")
	assert.True(t, msgs[0].IsStderr)
	assert.True(t, msgs[1].EOF)
}

func run(t *testing.T, cmd *Cmd, exitCode int) []Data {
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
	if exitCode != exitCode {
		t.Fatal(1)
	}
	wg.Wait()
	return msgs
}
