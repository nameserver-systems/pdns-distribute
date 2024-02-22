package worker

import (
	"runtime"
	"strconv"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/modeljob"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/powerdns"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

var Workerqueue chan *modeljob.PowerDNSAPIJob //nolint:gochecknoglobals

func StartWorker(conf *config.ServiceConfiguration) {
	workercount := conf.APIWorker

	Workerqueue = make(chan *modeljob.PowerDNSAPIJob, workercount)

	for num := 0; num < workercount; num++ {
		go worker(Workerqueue)
	}
}

func EnqueJob(j *modeljob.PowerDNSAPIJob) {
	Workerqueue <- j
	logger.DebugLog("[Enqueue Job] Job successfully enqueued. Jobs in workerqueue: " + strconv.Itoa(len(Workerqueue)))
	logger.DebugLog("[Enqueue Job] Go Routine Count: " + strconv.Itoa(runtime.NumGoroutine()))
	logger.DebugLog("[Enqueue Job] New Job: " + string(j.Msg.Data()))
}

func CloseWorkerQueue() {
	if Workerqueue != nil {
		close(Workerqueue)
	}
}

func worker(jobchan <-chan *modeljob.PowerDNSAPIJob) {
	for job := range jobchan {
		switch job.Jobtype {
		case modeljob.AddZone:
			powerdns.AddZone(job.Msg, job.Conf)

		case modeljob.ChangeZone:
			powerdns.ChangeZone(job.Msg, job.Conf)

		case modeljob.DeleteZone:
			powerdns.DeleteZone(job.Msg, job.Conf)
		}
	}
}
