package main

type Job struct {
	op func()
}

type WorkPool struct {
	numWorkers int
	jobs       chan Job
	jobCounter SafeCounter
	quit       chan bool
}

func NewWorkPool(numWorkers int) WorkPool {
	return WorkPool{
		jobs:       make(chan Job, numWorkers),
		jobCounter: NewSafeCounter(numWorkers),
		quit:       make(chan bool),
		numWorkers: numWorkers,
	}
}

func startWorker(pool *WorkPool) {
	for job := range pool.jobs {
		select {
		case <-pool.quit:
			return
		default:
			job.op()
			pool.jobCounter.Add()
			// fmt.Println("Finished Job: ", jobNum)
		}
	}
}

func (pool *WorkPool) Start() {
	pool.jobCounter.Start()
	for i := 0; i < pool.numWorkers; i++ {
		go startWorker(pool)
	}

}

func (pool *WorkPool) Stop() {
	pool.quit <- true
	pool.jobCounter.Stop()
}

func (pool *WorkPool) AddJob(job Job) {
	pool.jobs <- job
}
