package pipeline

import (
	"log"
	"sort"

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
		triggerStage     string
		Rules            []JobRules `yaml:",omitempty"`
	}
	PipelineIncludes struct {
		Repo string `yaml:",omitempty"`
		Ref  string `yaml:",omitempty"`
		File string `yaml:",omitempty"`
	}
	PipelineDefault struct {
		Tags []string `yaml:",omitempty"`
	}
)

func NewPipeline(name string) *Pipeline {
	return &Pipeline{
		Name:             name,
		Includes:         []PipelineIncludes{},
		Variables:        map[string]any{},
		TriggerVariables: map[string]any{},
		Stages:           []*Stage{},
		Rules:            []JobRules{},
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
func (this *Pipeline) Stage(name string) *Stage {
	stage := NewStage(name)
	this.Stages = append(this.Stages, stage)

	return stage
}

// BuildJob.AddRule("if ...", "always", false) // if, when, allow failure
func (this *Pipeline) AddRule(condition, when string, allowFailure bool) {
	this.Rules = append(this.Rules, JobRules{
		If:           &condition,
		When:         &when,
		AllowFailure: &allowFailure,
	})
}

func (this *Pipeline) AddIfWhenRule(condition, when string) {
	this.Rules = append(this.Rules, JobRules{
		If:   &condition,
		When: &when,
	})
}

func (this *Pipeline) AddIfRule(condition string) {
	this.Rules = append(this.Rules, JobRules{
		If: &condition,
	})
}

func (this *Pipeline) AddWhenRule(condition, when string) {
	this.Rules = append(this.Rules, JobRules{
		When: &when,
	})
}

func (this *Pipeline) AddExistsWhenRule(exists []string, when string) {
	this.Rules = append(this.Rules, JobRules{
		Exists: exists,
		When:   &when,
	})
}

func (this *Pipeline) AddChangesWhenRule(changes []string, when string) {
	this.Rules = append(this.Rules, JobRules{
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

	if len(this.Default.Tags) > 0 {
		out += "# Default\n"
		out += Marshal("default", this.Default)
		out += "\n"
	}

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
	jobs := []*Job{}

	for _, stage := range this.Stages {
		stages = append(stages, stage.Name)

		for _, job := range stage.Jobs {
			job.Stage = stage.Name

			jobs = append(jobs, job)
		}
	}
	out += "# Stages\n"
	out += Marshal("stages", stages)
	out += "\n"

	out += "#################################\n"
	out += "# Jobs\n"
	for _, stage := range this.Stages {
		out += "# Stage: " + stage.Name + "\n"

		keys := make([]string, 0, len(stage.Jobs))

		for k := range stage.Jobs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			job := stage.Jobs[k]
			job.Stage = stage.Name

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
