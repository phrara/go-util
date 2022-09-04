package workerpool

import "sync"



const (
	DefaultResQueLen = 100
)

type Processor interface {
	Process(task any, qid int) []error
	Handle(errs []error)
	addCnt()
	Count() int
	addPCnt()
}

type DefaultProcessor struct {
	count int
	pCnt int
	mu sync.Mutex
}

func Default() *DefaultProcessor {
	return &DefaultProcessor{
		count: 0,
		pCnt: 0,
		mu: sync.Mutex{},
	}
}

func (d *DefaultProcessor) Process(task any, qid int) []error {
	return nil
}

func (d *DefaultProcessor) Handle(errs []error) {
	
}
func (d *DefaultProcessor) addCnt() {
	d.count++
}

func (d *DefaultProcessor) Count() int {
	return d.count
}

func (d *DefaultProcessor) addPCnt() {
	d.mu.Lock()
	d.pCnt++
	d.mu.Unlock()
}

type ProcessorWithResQue struct {
	*DefaultProcessor

	resQue chan any
}

func WithResQue(resQueLen int) *ProcessorWithResQue {
	return &ProcessorWithResQue{
		resQue: make(chan any, resQueLen),
		DefaultProcessor: Default(),
	}
}

func (p *ProcessorWithResQue) PutResult(res any) {
	if p.resQue != nil {
		p.resQue <- res
	} else {
		p.resQue = make(chan any, DefaultResQueLen)
		p.resQue <- res
	}
}

func (p *ProcessorWithResQue) GetResults() chan any {
	if p.resQue != nil {
		return p.resQue
	} else {
		p.resQue = make(chan any, DefaultResQueLen)
		return p.resQue
	}
}

// count the number of processing
// when pcCnt == count, it means that all tasks have been processed
// so, the resQue could be closed 
func (d *ProcessorWithResQue) addPCnt() {
	d.mu.Lock()
	d.pCnt++
	if d.count == d.pCnt {
		close(d.resQue)
	}
	d.mu.Unlock()
}
