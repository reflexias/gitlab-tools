package pipeline

type (
	Stage struct {
		pipeline *Pipeline
		Name     string
		Jobs     []*Job
	}
)

func NewStage(name string) *Stage {
	return &Stage{
		Name: name,
		Jobs: []*Job{},
	}
}

// BuildJob = Build.CreateJob("Docker Build", "gcr.io/kaniko-project/executor:v1.19.1-debug", "") // Name, Image, Entrypoint
func (this *Stage) Job(name, image, entrypoint string) *Job {
	job := NewJob(name, image, entrypoint)
	this.Jobs = append(this.Jobs, job)
	job.stage = this

	return job
}
