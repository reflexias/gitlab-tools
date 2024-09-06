package pipeline

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type (
	Pipeline struct {
		id               string
		Name             string             `yaml:",omitempty"`
		Includes         []PipelineIncludes `yaml:",omitempty"`
		Variables        map[string]any     `yaml:",omitempty"`
		TriggerVariables map[string]any     `yaml:",omitempty"`
		Stages           []*Stage           `yaml:",omitempty"`
		Default          PipelineDefault    `yaml:",omitempty"`
		Rules            []*JobRule         `yaml:",omitempty"`
		triggerStage     string
	}
	PipelineIncludes struct {
		Repo     string `yaml:",omitempty"`
		Ref      string `yaml:",omitempty"`
		File     string `yaml:",omitempty"`
		Local    string `yaml:"local,omitempty"`
		Project  string `yaml:"project,omitempty"`
		Remote   string `yaml:"remote,omitempty"`
		Template string `yaml:"template,omitempty"`
	}
	PipelineDefault struct {
		AfterScript   []string  `yaml:"after_script,omitempty"`
		BeforeScript  []string  `yaml:"before_script,omitempty"`
		Image         string    `yaml:"image,omitempty"`
		Interruptible bool      `yaml:"interruptible,omitempty"`
		Retry         int       `yaml:"retry,omitempty"`
		Services      []Service `yaml:"services,omitempty"`
		Tags          []string  `yaml:"tags,omitempty"`
		Timeout       string    `yaml:"timeout,omitempty"`
	}
	Service struct {
		Name       string   `yaml:"name"`
		Alias      string   `yaml:"alias,omitempty"`
		Entrypoint []string `yaml:"entrypoint,omitempty"`
		Command    []string `yaml:"command,omitempty"`
	}
)

func NewPipeline(name string) *Pipeline {
	return &Pipeline{
		id:               uuid.NewString(),
		Name:             name,
		Includes:         []PipelineIncludes{},
		Variables:        map[string]any{},
		TriggerVariables: map[string]any{},
		Stages:           []*Stage{},
		Rules:            []*JobRule{},
		triggerStage:     name,
	}
}

func (this *Pipeline) ID() string {
	return this.id
}

func (this *Pipeline) SetID(id string) {
	this.id = id
}

func (this *Pipeline) SetTriggerStage(name string) {
	this.triggerStage = name
}

func (this *Pipeline) Include(repo, ref, file string) {
	this.Includes = append(this.Includes, PipelineIncludes{
		Repo: repo,
		Ref:  ref,
		File: file,
	})
}

// Pipeline.AddVariable("project", Project)
func (this *Pipeline) AddVariable(variable string, value any) {
	this.Variables[variable] = value
}

func (this *Pipeline) AddTriggerVariable(variable string, value any) {
	this.TriggerVariables[variable] = value
}

// Build = Pipeline.CreateStage("Build")
func (this *Pipeline) Stage(format string, a ...any) *Stage {
	name := fmt.Sprintf(format, a...)
	stage := NewStage(name)
	this.Stages = append(this.Stages, stage)
	stage.pipeline = this

	return stage
}

// BuildJob.AddRule("if ...", "always", false) // if, when, allow failure
func (this *Pipeline) AddRule(condition, when string, allowFailure bool) {
	this.Rules = append(this.Rules, &JobRule{
		If:           &condition,
		When:         &when,
		AllowFailure: &allowFailure,
	})
}

func (this *Pipeline) AddIfWhenRule(condition, when string) {
	this.Rules = append(this.Rules, &JobRule{
		If:   &condition,
		When: &when,
	})
}

func (this *Pipeline) AddIfRule(condition string) {
	this.Rules = append(this.Rules, &JobRule{
		If: &condition,
	})
}

func (this *Pipeline) AddWhenRule(condition, when string) {
	this.Rules = append(this.Rules, &JobRule{
		When: &when,
	})
}

func (this *Pipeline) AddExistsWhenRule(exists []string, when string) {
	this.Rules = append(this.Rules, &JobRule{
		Exists: exists,
		When:   &when,
	})
}

func (this *Pipeline) AddChangesWhenRule(changes []string, when string) {
	this.Rules = append(this.Rules, &JobRule{
		Changes: changes,
		When:    &when,
	})
}

func (this *Pipeline) Tags(tags ...string) {
	this.Default.Tags = append(this.Default.Tags, tags...)
}

func (this *Pipeline) Render() (out string) {
	out += "#################################\n"
	out += "# " + this.Name + " (" + this.id + ")\n"
	out += "#################################\n"
	out += "\n"

	out += "# Default\n"
	out += Marshal("default", this.Default)
	out += "\n"

	if len(this.Includes) > 0 {
		out += "# Includes\n"
		out += Marshal("include", this.Includes)
		out += "\n"
	}

	if len(this.Variables) > 0 {
		out += "# Variables\n"
		out += Marshal("variables", this.Variables)
		out += "\n"
	}
	stages := []string{}
	jobsNames := map[string]bool{}

	for _, stage := range this.Stages {
		stages = append(stages, stage.Name)

		for _, job := range stage.Jobs {
			job.Stage = stage.Name

			if _, ok := jobsNames[job.Name]; ok {
				log.Fatal("Duplicate job name: " + job.Name)
			}
			jobsNames[job.Name] = true
		}
	}
	out += "# Stages\n"
	out += Marshal("stages", stages)
	out += "\n"

	out += "#################################\n"
	out += "# Jobs\n"
	for _, stage := range this.Stages {
		log.Println("Rendering jobs for stage: " + stage.Name)
		out += "# Stage: " + stage.Name + "\n"

		for _, job := range stage.Jobs {
			job.Stage = stage.Name

			log.Println("Rendering job: " + job.Name)

			out += Marshal(job.Name, job)
			out += "\n"
		}
	}

	return
}

func Marshal(key string, o any) string {
	data := map[string]any{}
	data[key] = o

	out, err := yaml.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
