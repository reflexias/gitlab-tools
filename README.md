# Gitlab CI Pipeline tools

[![Go Reference](https://pkg.go.dev/badge/github.com/reflexias/gitlab-tools.svg)](https://pkg.go.dev/github.com/reflexias/gitlab-tools)

Generate gitlab-ci pipelines using Golang!  Simple interfaces and the rendered output is clean and ordered the way you defined it.

# Example

Clone this repo then:

```bash
mkdir -p output/
go run examples/simple/pipeline.go
cat output/deploy.yml
```