package raccoon

import (
	"sync"

	"fmt"

	log "github.com/Sirupsen/logrus"
)

type Dispatcher interface {
	Dispatch(js []Job)
}

type WorkerPoolDispatcher struct {
	Workers int
	sync.WaitGroup
}

type dispatchingJob struct {
	job  Job
	host Host
}

func (w *WorkerPoolDispatcher) Dispatch(jobs []Job) {
	log.WithFields(log.Fields{
		"package": "dispatcher",
	}).Info(fmt.Sprintf("%d workers pool selected", w.Workers))

	freeWorkersPool := make(chan chan dispatchingJob, w.Workers)
	dispatchingJobCh := make(chan dispatchingJob)

	//Launch workers
	for i := 0; i < w.Workers; i++ {
		ch := make(chan dispatchingJob)
		go w.worker(ch, freeWorkersPool)
		freeWorkersPool <- ch
	}

	go w.dispatcher(dispatchingJobCh, freeWorkersPool)

	for _, job := range jobs {
		w.Add(len(job.Cluster.Hosts))
		for _, host := range job.Cluster.Hosts {
			dispatchingJobCh <- dispatchingJob{job, host}
		}
	}

	w.Wait()

	//Close channels
	close(freeWorkersPool)
	close(dispatchingJobCh)
}

func (w *WorkerPoolDispatcher) dispatcher(dispatchingJobCh chan dispatchingJob, freeWorkersPool chan chan dispatchingJob) {
	for {
		dispatchingJob := <-dispatchingJobCh
		worker := <-freeWorkersPool
		worker <- dispatchingJob
	}
}

func (w *WorkerPoolDispatcher) worker(dispatchingJobCh chan dispatchingJob, freeWorkersPool chan chan dispatchingJob) {
	for {
		dispatchingJob := <-dispatchingJobCh
		fmt.Printf("Receiving host %s on ssh port %d\n", dispatchingJob.host.IP, dispatchingJob.host.SSH_port)

		dispatcher := SequentialDispatcher{}
		dispatcher.executeJobOnHost(dispatchingJob.job, dispatchingJob.host)
		w.Done()

		freeWorkersPool <- dispatchingJobCh
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
		}).Error("Error initializing node: " + err.Error())
		return
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
		}).Error("Error initializing node: " + err.Error())
		return
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
