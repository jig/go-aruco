package aruco

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jig/go-exec"
)

type ExecPython struct {
	pythonCommand  *exec.Cmd
	markersChannel chan Markers
}

func New() (*ExecPython, error) {
	pythonCommand, err := exec.NewCmd("python", "-c", markersPythonCode)
	if err != nil {
		return nil, err
	}
	return &ExecPython{pythonCommand: pythonCommand}, nil
}

func (exec *ExecPython) Run() (int, error) {
	return exec.pythonCommand.Run()
}

func (exec *ExecPython) Dispatch(f func(markers Markers)) {
	for data := range exec.pythonCommand.Output() {
		if data.EOF {
			close(exec.markersChannel)
			return
		} else if !data.IsStderr {
			markers := Markers{}
			if err := json.Unmarshal([]byte(data.Value), &markers); err != nil {
				log.Println(err)
				continue
			}
			f(markers)
		} else {
			fmt.Println(data.Value)
		}
	}
}

func (exec *ExecPython) Output() chan Markers {
	return exec.markersChannel
}

//go:embed markers.py
var markersPythonCode string
