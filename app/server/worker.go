package server

import (
	"sync"
)

type Task func()

type Worker struct {
	id         int
	taskQueue  chan Task
	workerPool chan chan Task
	stop       chan bool
}

func NewWorker(id int, workerPool chan chan Task) *Worker {
	return &Worker{
		id:         id,
		taskQueue:  make(chan Task),
		workerPool: workerPool,
		stop:       make(chan bool),
	}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			w.workerPool <- w.taskQueue
			select {
			case task := <-w.taskQueue:
				task()
			case <-w.stop:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.stop <- true
	}()
}

type ThreadPool struct {
	workerPool chan chan Task
	workers    []*Worker
	wg         sync.WaitGroup
}

func NewThreadPool(numWorkers int) *ThreadPool {
	workerPool := make(chan chan Task, numWorkers)
	workers := make([]*Worker, numWorkers)

	for i := range numWorkers {
		workers[i] = NewWorker(i, workerPool)
	}

	return &ThreadPool{
		workerPool: workerPool,
		workers:    workers,
		wg:         sync.WaitGroup{},
	}
}

func (tp *ThreadPool) Start() {
	for _, worker := range tp.workers {
		tp.wg.Add(1)
		worker.Start(&tp.wg)
	}
}

func (tp *ThreadPool) Submit(task Task) {
	workerTaskQueue := <-tp.workerPool
	workerTaskQueue <- task
}

func (tp *ThreadPool) Stop() {
	for _, worker := range tp.workers {
		worker.Stop()
	}
	close(tp.workerPool)
	tp.wg.Wait()
}
