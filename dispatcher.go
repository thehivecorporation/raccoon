package raccoon

import (
	"sync"

	log "github.com/Sirupsen/logrus"
)

type Dispatcher interface {
	Dispatch(js []Job)
}

type WorkerPoolDispatcher struct {
	Workers int
}

func (w *WorkerPoolDispatcher) Dispatch(jobs []Job) {
	freeWorkersPool := make(chan chan Job, w.Workers)

	jobsCh := make(chan Job)

	go w.dispatcher(jobsCh, freeWorkersPool)

	for _, job := range jobs {
		jobsCh <- job
	}

	//Close channels
	close(freeWorkersPool)
	close(jobsCh)
}

func (w *WorkerPoolDispatcher) dispatcher(jobsCh chan Job, freeWorkersPool chan chan Job) {
	for {
		job := <-jobsCh
		log.WithFields(log.Fields{
			"cluster":        job.Cluster.Name,
			"infrastructure": job.Task.Title,
			"maintainer":     job.Task.Maintainer,
			"package":        "dispatcher",
		}).Info("Launching Raccoon...")

		worker := <-freeWorkersPool
		worker <- job
	}
}

func (w *WorkerPoolDispatcher) worker(jobsCh chan Job, freeWorkersPool chan chan Job) {
	for {
		job := <-jobsCh

		dispatcher := SequentialDispatcher{}
		dispatcher.Dispatch([]Job{job})

		freeWorkersPool <- jobsCh
	}
}

type SimpleDispatcher struct {
	sync.WaitGroup
}

func (s *SimpleDispatcher) Dispatch(jobs []Job) {
	for _, job := range jobs {
		for _, node := range job.Cluster.Hosts {
			s.Add(1)
			go s.executeJobOnHost(job, node)
		}
	}

	s.Wait()
}

func (s *SimpleDispatcher) executeJobOnHost(j Job, h Host) {
	err := h.InitializeNode()
	if err != nil {
		log.WithFields(log.Fields{
			"host":    h.IP,
			"package": "dispatcher",
		}).Warn("Error initializing node: " + err.Error())
	}

	for _, instruction := range j.Task.Commands {
		instruction.Execute(h)
	}

	err = h.CloseNode()
	if err != nil {
		log.WithFields(log.Fields{
			"host":    h.IP,
			"package": "dispatcher",
		}).Warn("Error closing session: " + err.Error())
	}

	s.Done()
}

type SequentialDispatcher struct{}

func (s *SequentialDispatcher) Dispatch(jobs []Job) {
	for _, job := range jobs {
		log.WithFields(log.Fields{
			"cluster":        job.Cluster.Name,
			"infrastructure": job.Task.Title,
			"maintainer":     job.Task.Maintainer,
			"package":        "dispatcher",
		}).Info("Launching Raccoon...")

		for _, node := range job.Cluster.Hosts {
			s.executeJobOnHost(job, node)
		}
	}
}

func (s *SequentialDispatcher) executeJobOnHost(j Job, h Host) {
	err := h.InitializeNode()
	if err != nil {
		log.WithFields(log.Fields{
			"host":    h.IP,
			"package": "dispatcher",
		}).Warn("Error initializing node: " + err.Error())
	}

	for _, instruction := range j.Task.Commands {
		instruction.Execute(h)
	}

	err = h.CloseNode()
	if err != nil {
		log.WithFields(log.Fields{
			"host":    h.IP,
			"package": "dispatcher",
		}).Warn("Error closing session: " + err.Error())
	}
}
