package workerpool

import "sync"

type Strategy int

const (
	// RR select que by using Round Robin
	RR Strategy = 0
	// SRC select que by using outer source
	SRC Strategy = 1
)

type WorkerPool struct {
	MaxQueNum int
	AvgWorker int
	MaxQueLen int
	taskQueue []chan any
	processor    Processor
	qid       int
	Mod       Strategy
	mux sync.Mutex
}

func New(maxQueNum, avgWorker, maxQueLen int, processor Processor, mod Strategy) *WorkerPool {
	return &WorkerPool{
		MaxQueNum: maxQueNum,
		AvgWorker: avgWorker,
		MaxQueLen: maxQueLen,
		taskQueue: make([]chan any, maxQueNum),
		processor:    processor,
		qid:       0,
		Mod:       mod,
		mux: sync.Mutex{},
	}
}

func (w *WorkerPool) work(qid int) {
	for task := range w.taskQueue[qid] {
		// process the task
		errs := w.processor.Process(task, qid)
		// count the number of processing
		w.processor.addPCnt()
		if errs != nil {
			// handle the errors
			w.processor.Handle(errs)
		}
	}
}

func (w *WorkerPool) Start() {
	for i := 0; i < w.MaxQueNum; i++ {
		w.taskQueue[i] = make(chan any, w.MaxQueLen)
		for j := 0; j < w.AvgWorker; j++ {
			go w.work(i)
		}
	}
}

// AppendTask used Round Robin or Source
func (w *WorkerPool) AppendTask(task any, src int) {
	switch w.Mod {
	case RR:
		w.mux.Lock()
		w.qid = w.qid % w.MaxQueNum
		w.taskQueue[w.qid] <- task
		//fmt.Printf("task que %d recv a task\n", w.qid)
		w.qid++
		w.mux.Unlock()
	case SRC:
		qid := src % w.MaxQueNum
		w.taskQueue[qid] <- task
		//fmt.Printf("task que %d recv a task\n", w.qid)
	default:
		return
	}
	w.processor.addCnt()
}

func (w *WorkerPool) Shut() {
	for i := 0; i < w.MaxQueNum; i++ {
		close(w.taskQueue[i])
	}
}
