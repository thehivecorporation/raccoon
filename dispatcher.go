package raccoon

import (
	"sync"

	"fmt"

	log "github.com/Sirupsen/logrus"
)

//Dispatcher is the interface to implement new dispatching strategies in the app
type Dispatcher interface {
	Dispatch(js []Job)
}

//WorkerPoolDispatcher is a dispatcher that acts as a worker pool. It won't
//maintain more connections to hosts atthe same time than the specified number
//in the Workers field
type WorkerPoolDispatcher struct {
	Workers int
	sync.WaitGroup
}

type dispatchingJob struct {
	job  Job
	host Host
}

//Dispatch will create the specified workers on the workers field, launching a
//goroutine with each. Then it will create a dispatcher routing in a different
//goroutine and finally it will send the jobs to the launched dispatcher for
//distribution.
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
		fmt.Printf("Receiving host %s on ssh port %d\n", dispatchingJob.host.IP, dispatchingJob.host.SSHPort)

		dispatcher := SequentialDispatcher{}
		dispatcher.executeJobOnHost(dispatchingJob.job, dispatchingJob.host)
		w.Done()

		freeWorkersPool <- dispatchingJobCh
	}
}

//SimpleDispatcher is the dispatching strategy to use when the number of hosts
//isn't very high. It will open a new SSH connection to each host at the same
//time since the beginning so, if you are having performance issues, try the
//workers pool or the sequential dispatcher to see if they get solved.
//SimpleDispatcher is the default dispatching strategy
type SimpleDispatcher struct {
	sync.WaitGroup
}

//Dispatch will create a new goroutine for each host in the job and launch the
//tasks that were specified on it.
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

//SequentialDispatcher is a dispatching strategy that doesn't use any concurrent
//approach. It will execute each command on each host sequentially and waits for
//each to finish before starting the next.
type SequentialDispatcher struct{}

//Dispatch ranges over the hosts to execute the tasks on them sequentially without any kind of concurrency.
func (s *SequentialDispatcher) Dispatch(jobs []Job) {
	for _, job := range jobs {
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
