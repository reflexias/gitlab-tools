package pipeline

import "github.com/google/uuid"

type (
	Workflow struct {
		ID               string
		Pipelines        []*Pipeline
		Default          PipelineDefault `yaml:",omitempty"`
		Variables        map[string]any  `yaml:",omitempty"`
		GenerateImage    string
		GenerateCommands []string
	}
)

func NewWorkflow() *Workflow {
	return &Workflow{
		Pipelines: []*Pipeline{},
		Variables: map[string]any{},
		GenerateCommands: []string{
			"go run test.go",
		},
		GenerateImage: "ubuntu:latest",
	}
}

func (this *Workflow) SetGenerateImage(name string) {
	this.GenerateImage = name
}

func (this *Workflow) SetGenerateCommands(commands []string) {
	this.GenerateCommands = commands
}

func (this *Workflow) Tags(tags ...string) {
	this.Default.Tags = append(this.Default.Tags, tags...)
}

func (this *Workflow) CreatePipeline(name string) *Pipeline {
	pipeline := NewPipeline(name)
	this.Pipelines = append(this.Pipelines, pipeline)

	return pipeline
}

func (this *Workflow) AddVariable(variable string, value any) {
	this.Variables[variable] = value
}

func (this *Workflow) Render() (out string) {
	id := uuid.NewString()
	this.ID = id
	out += "##################################################################\n"
	out += "# Dynamic Job ID: " + id + "\n"
	out += "##################################################################\n"
	out += "\n"

	if len(this.Default.Tags) > 0 {
		out += "# Default\n"
		out += Marshal("default", this.Default)
		out += "\n"
	}

	this.AddVariable("DYNAMIC_JOB_ID", id)

	if len(this.Variables) > 0 {
		out += "# Variables\n"
		out += Marshal("variables", this.Variables)
		out += "\n"
	}

	artifacts := []string{}
	stages := []string{"generate"}
	stageMap := map[string]bool{}

	// Add Pipeline Stages
	for _, pipeline := range this.Pipelines {
		artifacts = append(artifacts, "output/"+pipeline.Name+".yml")
		if _, ok := stageMap[pipeline.triggerStage]; !ok {
			stages = append(stages, pipeline.triggerStage)
			stageMap[pipeline.triggerStage] = true
		}
	}

	out += "# Stages\n"
	out += Marshal("stages", stages)
	out += "\n"

	def := &Job{
		Stage: "generate",
		Image: &JobImage{
			Name: this.GenerateImage,
		},
		Script:    this.GenerateCommands,
		Variables: map[string]any{},
		Artifacts: &Artifacts{Paths: artifacts},
	}

	out += "# Generate Jobs Here!\n"
	out += Marshal("generate", def)
	out += "\n"

	// Add Pipeline Jobs
	for _, pipeline := range this.Pipelines {
		def := &Job{
			Stage: pipeline.triggerStage,
			Inherit: JobInherit{
				Variables: true,
			},
			Variables: pipeline.TriggerVariables,
			Trigger: JobTrigger{
				Strategy: "depend",
				Include: []JobTriggerInclude{
					{
						Artifact: "output/" + pipeline.Name + ".yml",
						Job:      "generate",
					},
				},
			},
		}

		name := "Trigger " + pipeline.Name
		out += "# Trigger " + pipeline.Name + "\n"
		out += Marshal(name, def)
		out += "\n"
	}

	return
}
