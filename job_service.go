package main

type JobService struct {
	ctx *Context
	c   chan Job
}

func CreateJobService(ctx *Context) {
	jobService := JobService{ctx, make(chan Job)}
	for i := 0; i < 5; i++ {
		go jobService.runJobProcessor()
	}

	ctx.JobService = &jobService
}

func (self *JobService) Channel() chan Job {
	return self.c
}

func (self *JobService) runJobProcessor() {
	for {
		job := <-self.Channel()
		job.Execute()
	}
}
