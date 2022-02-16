package pipeline

type (
	Job struct {
		Stage        string
		Name         string             `yaml:"-"`
		Image        JobImage           `yaml:",omitempty"`
		Variables    map[string]any     `yaml:",omitempty"`
		Secrets      map[string]Secret  `yaml:",omitempty"`
		IDTokens     map[string]IDToken `yaml:"id_tokens,omitempty"`
		Dependencies []string           `yaml:",omitempty"`
		Needs        []string           `yaml:",omitempty"`
		Extends      []string           `yaml:",omitempty"`
		Script       []string           `yaml:",omitempty"`
		Artifacts    []string           `yaml:",omitempty"`
		PullPolicy   *string            `json:"pull_policy,omitempty" yaml:"pull_policy,omitempty"`
		When         string             `yaml:",omitempty"`
		Trigger      JobTrigger         `yaml:",omitempty"`
		Inherit      JobInherit         `yaml:",omitempty"`
		Cache        []JobCache         `yaml:",omitempty"`
		Environment  JobEnvironment     `yaml:",omitempty"`
		Rules        []JobRules         `yaml:",omitempty"`
	}
	JobRules struct {
		Exists       []string `yaml:",omitempty"`
		Changes      []string `yaml:",omitempty"`
		When         *string  `yaml:",omitempty"`
		If           *string  `yaml:",omitempty"`
		AllowFailure *bool    `yaml:"allow_failure,omitempty"`
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
		Engine SecretEngine `yaml:",omitempty"`
		Path   string       `yaml:",omitempty"`
		Field  string       `yaml:",omitempty"`
	}
	SecretEngine struct {
		Path string `yaml:",omitempty"`
		Name string `yaml:",omitempty"`
	}
	IDToken struct {
		Aud string `yaml:",omitempty"`
	}
	JobEnvironment struct {
		Name   string `yaml:",omitempty"`
		Url    string `yaml:",omitempty"`
		Action string `yaml:",omitempty"`
		Tier   string `yaml:"deployment_tier,omitempty"`
	}
)

func NewJob(name, image, entrypoint string) *Job {
	job := &Job{
		Name:      name,
		Variables: map[string]any{},
		Secrets:   map[string]Secret{},
		Extends:   []string{},
		Script:    []string{},
		Artifacts: []string{},
		Rules:     []JobRules{},
	}

	if image != "" {
		job.Image = JobImage{
			Name:       image,
			Entrypoint: entrypoint,
		}
	}
	return job
}

func (this *Job) Extend(name string) {
	this.Extends = append(this.Extends, name)
}

func (this *Job) Need(name string) {
	this.Needs = append(this.Needs, name)
}

func (this *Job) NeedsJob(j *Job) {
	this.Needs = append(this.Needs, j.Name)
}

func (this *Job) Dependency(name string) {
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
	this.Secrets[Variable] = Secret{
		Engine: SecretEngine{
			Name: engine,
			Path: enginePath,
		},
		Path:  secretPath,
		Field: field,
	}
}

func (this *Job) AddIDToken(name, aud string) {
	if this.IDTokens == nil {
		this.IDTokens = map[string]IDToken{}
	}
	this.IDTokens[name] = IDToken{
		Aud: aud,
	}
}

func (this *Job) SetWhen(when string) {
	this.When = when
}

func (this *Job) SetEnvironment(name, action, url, tier string) {
	this.Environment = JobEnvironment{
		Name:   name,
		Action: action,
		Url:    url,
		Tier:   tier,
	}
}

// BuildJob.AddRule("if ...", "always", false) // if, when, allow failure
func (this *Job) AddRule(condition, when string, allowFailure bool) {
	this.Rules = append(this.Rules, JobRules{
		If:           &condition,
		When:         &when,
		AllowFailure: &allowFailure,
	})
}

func (this *Job) AddIfWhenRule(condition, when string) {
	this.Rules = append(this.Rules, JobRules{
		If:   &condition,
		When: &when,
	})
}

func (this *Job) AddIfRule(condition string) {
	this.Rules = append(this.Rules, JobRules{
		If: &condition,
	})
}

func (this *Job) AddWhenRule(condition, when string) {
	this.Rules = append(this.Rules, JobRules{
		When: &when,
	})
}

func (this *Job) AddExistsWhenRule(exists []string, when string) {
	this.Rules = append(this.Rules, JobRules{
		Exists: exists,
		When:   &when,
	})
}

func (this *Job) AddChangesWhenRule(changes []string, when string) {
	this.Rules = append(this.Rules, JobRules{
		Changes: changes,
		When:    &when,
	})
}

func (this *Job) AddArtifact(file string) {
	this.Artifacts = append(this.Artifacts, file)
}

func (this *Job) AddCache(key string, paths ...string) {
	if len(paths) > 0 {
		this.Cache = append(this.Cache, JobCache{
			Key:   key,
			Paths: paths,
		})
	}
}
