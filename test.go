package main

import (
	"log"
	"os"

	"github.com/reflexias/gitlab-tools/pkg/getfile"
	"github.com/reflexias/gitlab-tools/pkg/pipeline"
	"gopkg.in/yaml.v3"
)

type (
	Data struct {
		Things       []string
		Environments []string
	}
)

func main() {
	// Get a file
	rawData, err := getfile.GitFile("", os.Getenv("GITLAB_TOKEN"), "booleansnailfish/vital-ci-example", "data.yaml", "")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the file
	data := &Data{}
	err = yaml.Unmarshal(rawData, data)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new workflow
	workflow := pipeline.NewWorkflow()
	workflow.AddVariable("foo", "bar")

	// Create the build pipeline
	buildPipeline := workflow.CreatePipeline("build")
	buildPipeline.AddTriggerVariable("PARENT_PIPELINE_ID", "$CI_PIPELINE_ID")
	buildPipeline.AddIfWhenRule("$SKIP_BUILD == 'true'", "never")
	buildPipeline.AddIfWhenRule("$CI_COMMIT_TAG != null", "never")

	buildStage := buildPipeline.Stage("build")

	for _, thing := range data.Things {
		log.Println("Adding build job for:", thing)
		buildJob := buildStage.Job("Build "+thing, "ubuntu:latest", "")
		buildJob.AddCommand("env")

		publishJob := buildStage.Job("Publish "+thing, "ubuntu:latest", "")
		publishJob.AddCommand("env")
		publishJob.DependsOnJob(buildJob)

		releaseJob := buildStage.Job("Release "+thing, "ubuntu:latest", "")
		releaseJob.AddCommand("env")
		releaseJob.DependsOnJob(publishJob)
	}

	// Create Compliance Pipeline
	compliancePipeline := workflow.CreatePipeline("compliance")
	compliancePipeline.AddTriggerVariable("PARENT_PIPELINE_ID", "$CI_PIPELINE_ID")
	compliancePipeline.SetTriggerStage("build")
	complianceStage := compliancePipeline.Stage("compliance")

	unitTest := complianceStage.Job("Run Unit Test", "ubuntu:latest", "")
	unitTest.AddCommand("env")

	codeCoverage := complianceStage.Job("Run Code Coverage", "ubuntu:latest", "")
	codeCoverage.AddCommand("env")

	// Create the deploy pipeline
	deployPipeline := workflow.CreatePipeline("deploy")
	deployPipeline.AddTriggerVariable("PARENT_PIPELINE_ID", "$CI_PIPELINE_ID")
	deployPipeline.AddIfWhenRule("$CI_PIPELINE_SOURCE == 'merge_request_event'", "never")
	deployPipeline.AddIfWhenRule("$CI_COMMIT_BRANCH != $CI_DEFAULT_BRANCH", "never")

	for _, environment := range data.Environments {
		log.Println("Adding deploy Stage for:", environment)
		deployStage := deployPipeline.Stage("deploy-" + environment)

		planJob := deployStage.Job("Plan "+environment, "ubuntu:latest", "")
		planJob.AddCommand("env")

		deployJob := deployStage.Job("Deploy "+environment, "ubuntu:latest", "")
		deployJob.AddCommand("env")
		deployJob.Need("Plan " + environment)
		deployJob.SetEnvironment(environment, "deploy", "", "")

		smokeJob := deployStage.Job("Smoke "+environment, "ubuntu:latest", "")
		smokeJob.AddCommand("env")
		smokeJob.Need("Deploy " + environment)
	}

	// Output the main pipeline
	log.Println("Writing: pipeline.yml")
	err = os.WriteFile("output/pipeline.yml", []byte(workflow.Render()), 0660)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(".gitlab-ci.yml", []byte(workflow.Render()), 0660)
	if err != nil {
		log.Fatal(err)
	}

	// Output the Child Pipelines
	for _, pipeline := range workflow.Pipelines {
		log.Println("Writing:", pipeline.Name+".yml")
		err := os.WriteFile("output/"+pipeline.Name+".yml", []byte(pipeline.Render()), 0660)
		if err != nil {
			log.Fatal(err)
		}
	}
}
