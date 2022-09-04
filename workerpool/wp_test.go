package workerpool

import (
	"fmt"
	"testing"
)


type MyProcessor struct {
	*ProcessorWithResQue
}

func (m *MyProcessor) Process(task any, qid int) []error {
	
	a := task.(Add)
	res := a.x + a.y
	m.PutResult(res)
	
	return nil
}

func (m *MyProcessor) Handle(errs []error) {
	
}

type Add struct {
	x, y int
}

func TestWP(t *testing.T) {
	prc := &MyProcessor{
		ProcessorWithResQue: WithResQue(DefaultResQueLen),
	}
	wp := New(3, 3, 100, prc, RR)

	wp.Start()

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			wp.AppendTask(Add{i ,j}, 0)
		}
	}

	for v := range prc.GetResults() {
		fmt.Println(v.(int))
	}

	wp.Shut()

}