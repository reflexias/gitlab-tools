package pipeline

type (
	Stage struct {
		Name      string
		Variables map[string]any
		Jobs      map[string]*Job
	}
)

func NewStage(name string) *Stage {
	return &Stage{
		Name:      name,
		Variables: map[string]any{},
		Jobs:      map[string]*Job{},
	}
}

// BuildJob = Build.CreateJob("Docker Build", "gcr.io/kaniko-project/executor:v1.19.1-debug", "") // Name, Image, Entrypoint
func (this *Stage) Job(name, image, entrypoint string) *Job {
	job := NewJob(name, image, entrypoint)
	this.Jobs[name] = job

	return job
}
