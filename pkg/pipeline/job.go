package pipeline

import "fmt"

type (
	Job struct {
		Stage        string
		stage        *Stage
		Name         string              `yaml:"-"`
		Image        *JobImage           `yaml:",omitempty"`
		Variables    map[string]any      `yaml:",omitempty"`
		Secrets      map[string]*Secret  `yaml:",omitempty"`
		IDTokens     map[string]*IDToken `yaml:"id_tokens,omitempty"`
		Dependencies []string            `yaml:",omitempty"`
		Needs        []string            `yaml:",omitempty"`
		Extends      []string            `yaml:",omitempty"`
		Script       []string            `yaml:",omitempty"`
		Artifacts    *Artifacts          `yaml:",omitempty"`
		PullPolicy   *string             `json:"pull_policy,omitempty" yaml:"pull_policy,omitempty"`
		When         string              `yaml:",omitempty"`
		Trigger      JobTrigger          `yaml:",omitempty"`
		Inherit      JobInherit          `yaml:",omitempty"`
		Cache        []*JobCache         `yaml:",omitempty"`
		Environment  Environment         `yaml:",omitempty"`
		Rules        []*JobRule          `yaml:",omitempty"`
		BeforeScript []string            `yaml:"before_script,omitempty"`
		AfterScript  []string            `yaml:"after_script,omitempty"`
		AllowFailure bool                `yaml:"allow_failure,omitempty"`
		Retry        int                 `yaml:"retry,omitempty"`
		Services     []*Service          `yaml:"services,omitempty"`
		Tags         []string            `yaml:"tags,omitempty"`
		Timeout      string              `yaml:"timeout,omitempty"`
	}
	Artifacts struct {
		Paths []string `yaml:",omitempty"`
	}
	JobRule struct {
		Exists       []string          `yaml:",omitempty"`
		Changes      []string          `yaml:",omitempty"`
		When         string            `yaml:",omitempty"`
		Variables    map[string]string `yaml:"variables,omitempty"`
		If           *string           `yaml:",omitempty"`
		AllowFailure *bool             `yaml:"allow_failure,omitempty"`
	}
	JobImage struct {
		Name       string `yaml:",omitempty"`
		Entrypoint string `yaml:",omitempty"`
	}
	JobTrigger struct {
		Strategy string              `yaml:",omitempty"`
		Include  []JobTriggerInclude `yaml:",omitempty"`
	}
	JobTriggerInclude struct {
		Artifact string `yaml:",omitempty"`
		Job      string `yaml:",omitempty"`
		Remote   string `yaml:",omitempty"`
	}
	JobInherit struct {
		Variables bool `yaml:",omitempty"`
	}
	JobCache struct {
		Key   string   `yaml:",omitempty"`
		Paths []string `yaml:",omitempty"`
	}
	Secret struct {
		Vault VaultSecret `yaml:",omitempty"`
	}
	VaultSecret struct {
		Engine SecretEngine `yaml:",omitempty"`
		Path   string       `yaml:",omitempty"`
		Field  string       `yaml:",omitempty"`
	}
	SecretEngine struct {
		Path string `yaml:",omitempty"`
		Name string `yaml:",omitempty"`
	}
	IDToken struct {
		Aud []string `yaml:",omitempty"`
	}
	Environment struct {
		Name   string `yaml:",omitempty"`
		Url    string `yaml:",omitempty"`
		Action string `yaml:",omitempty"`
		Tier   string `yaml:"deployment_tier,omitempty"`
	}
)

func NewJob(format string, a ...any) *Job {
	name := fmt.Sprintf(format, a...)
	job := &Job{
		Name:      name,
		Variables: map[string]any{},
		Secrets:   map[string]*Secret{},
		Extends:   []string{},
		Script:    []string{},
		Rules:     []*JobRule{},
	}

	return job
}

func (this *Job) Extend(format string, a ...any) {
	name := fmt.Sprintf(format, a...)
	this.Extends = append(this.Extends, name)
}

func (this *Job) SetImage(format string, a ...any) {
	name := fmt.Sprintf(format, a...)
	if this.Image == nil {
		this.Image = &JobImage{}
	}
	this.Image.Name = name
}

func (this *Job) SetEntrypoint(entrypoint string) {
	if this.Image == nil {
		this.Image = &JobImage{}
	}
	this.Image.Entrypoint = entrypoint
}

func (this *Job) Need(format string, a ...any) {
	name := fmt.Sprintf(format, a...)
	this.Needs = append(this.Needs, name)
}

func (this *Job) NeedsJob(j *Job) {
	this.Needs = append(this.Needs, j.Name)
}

func (this *Job) Dependency(format string, a ...any) {
	name := fmt.Sprintf(format, a...)
	this.Dependencies = append(this.Dependencies, name)
}
func (this *Job) DependsOnJob(j *Job) {
	this.Dependencies = append(this.Dependencies, j.Name)
}

func (this *Job) AddVariable(variable string, value any) {
	this.Variables[variable] = value
}

func (this *Job) AddCommand(command string) {
	this.Script = append(this.Script, command)
}

// Add a Vault Secret CICD Variable, engine, engine-path, secret path, field
// https://docs.gitlab.com/ee/ci/yaml/#secretsvault
func (this *Job) AddVaultSecret(Variable, engine, enginePath, secretPath, field string) {
	this.Secrets[Variable] = &Secret{
		Vault: VaultSecret{
			Engine: SecretEngine{
				Name: engine,
				Path: enginePath,
			},
			Path:  secretPath,
			Field: field,
		},
	}
}

func (this *Job) AddIDToken(name, aud string) {
	if this.IDTokens == nil {
		this.IDTokens = map[string]*IDToken{}
	}
	this.IDTokens[name] = &IDToken{
		Aud: []string{aud},
	}
}

func (this *Job) SetWhen(when string) {
	this.When = when
}

func (this *Job) SetEnvironment(name, action, url, tier string) {
	this.Environment = Environment{
		Name:   name,
		Action: action,
		Url:    url,
		Tier:   tier,
	}
}

// BuildJob.AddRule("if ...", "always", false) // if, when, allow failure
func (this *Job) AddRule(condition, when string, allowFailure bool) {
	this.Rules = append(this.Rules, &JobRule{
		If:           &condition,
		When:         &when,
		AllowFailure: &allowFailure,
	})
}

func (this *Job) AddIfWhenRule(condition, when string) {
	this.Rules = append(this.Rules, &JobRule{
		If:   &condition,
		When: &when,
	})
}

func (this *Job) AddIfRule(condition string) {
	this.Rules = append(this.Rules, &JobRule{
		If: &condition,
	})
}

func (this *Job) AddWhenRule(condition, when string) {
	this.Rules = append(this.Rules, &JobRule{
		When: &when,
	})
}

func (this *Job) AddExistsWhenRule(exists []string, when string) {
	this.Rules = append(this.Rules, &JobRule{
		Exists: exists,
		When:   &when,
	})
}

func (this *Job) AddChangesWhenRule(changes []string, when string) {
	this.Rules = append(this.Rules, &JobRule{
		Changes: changes,
		When:    &when,
	})
}

func (this *Job) AddArtifact(format string, a ...any) {
	file := fmt.Sprintf(format, a...)
	if this.Artifacts == nil {
		this.Artifacts = &Artifacts{
			Paths: []string{},
		}
	}

	this.Artifacts.Paths = append(this.Artifacts.Paths, file)
}

func (this *Job) AddCache(key string, paths ...string) {
	if len(paths) > 0 {
		this.Cache = append(this.Cache, &JobCache{
			Key:   key,
			Paths: paths,
		})
	}
}
