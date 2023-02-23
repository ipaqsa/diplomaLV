package service

func NewQueue() *Queue {
	return &Queue{
		data: make([]Task, 5),
	}
}

func (queue *Queue) Push(t Task) {
	queue.mtx.Lock()
	defer queue.mtx.Unlock()
	queue.data = append(queue.data, t)
}

func (queue *Queue) Pop() Task {
	queue.mtx.Lock()
	defer queue.mtx.Unlock()
	data := queue.data[0]
	queue.data = queue.data[1:]
	return data
}

func (queue *Queue) Len() int {
	queue.mtx.Lock()
	defer queue.mtx.Unlock()
	return len(queue.data)
}

func (queue *Queue) IsEmpty() bool {
	queue.mtx.Lock()
	defer queue.mtx.Unlock()
	return len(queue.data) == 0
}
