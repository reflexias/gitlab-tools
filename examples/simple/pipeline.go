package main

import (
	"log"
	"os"

	"github.com/reflexias/gitlab-tools/pkg/pipeline"
)

var environments = []string{"dev", "staging", "prod"}

func main() {
	// Create the deploy pipeline
	deployPipeline := pipeline.NewPipeline("deploy")

	for _, environment := range environments {
		log.Println("Adding deploy Stage for:", environment)
		deployStage := deployPipeline.Stage("deploy-%s", environment)

		planJob := deployStage.Job("Plan %s", environment)
		planJob.SetImage("ubuntu:latest")
		planJob.AddCommand("env")
		planJob.SetEnvironment(environment, "prepare", "", "")

		deployJob := deployStage.Job("Deploy %s", environment)
		deployJob.SetImage("ubuntu:latest")
		deployJob.AddCommand("env")
		deployJob.NeedsJob(planJob)
		deployJob.SetEnvironment(environment, "start", "", "")
		deployJob.AddVaultSecret("DB_PASSWORD", "kv-v2", "ops", environment+"/db", "password")

		smokeJob := deployStage.Job("Smoke %s", environment)
		smokeJob.SetImage("ubuntu:latest")
		smokeJob.AddCommand("env")
		smokeJob.NeedsJob(deployJob)
		smokeJob.SetEnvironment(environment, "verify", "", "")
	}

	err := os.WriteFile("output/deploy.yml", []byte(deployPipeline.Render()), 0660)
	if err != nil {
		log.Fatal(err)
	}
}
