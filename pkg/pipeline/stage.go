package pipeline

import "fmt"

type (
	Stage struct {
		pipeline *Pipeline
		Name     string
		Jobs     []*Job
	}
)

func NewStage(format string, a ...any) *Stage {
	name := fmt.Sprintf(format, a...)
	return &Stage{
		Name: name,
		Jobs: []*Job{},
	}
}

// BuildJob = Build.CreateJob("Docker Build", "gcr.io/kaniko-project/executor:v1.19.1-debug", "") // Name, Image, Entrypoint
func (this *Stage) Job(format string, a ...any) *Job {
	name := fmt.Sprintf(format, a...)

	job := NewJob(name)
	this.Jobs = append(this.Jobs, job)
	job.stage = this

	return job
}
