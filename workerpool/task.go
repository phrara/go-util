package workerpool

type Tasker interface {
	Process(task any) error
	Handle(err error)
}

type Task struct {
}

func (t *Task) Process(task any) error {
	//TODO implement me
	panic("implement me")
}

func (t *Task) Handle(err error) {
	//TODO implement me
	panic("implement me")
}

