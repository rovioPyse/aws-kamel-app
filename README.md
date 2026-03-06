# aws-kamel-app

Backend repository for Kamel serverless services.

## What is in this repo

- Two Go Lambda services, each as an independent Go module:
  - `services/lambdas/latest-state-writer`
  - `services/lambdas/alarm-processor`
- One SAM application that wires and deploys these Lambdas:
  - `deploy/sam/kalmel-app-infra`

## Current architecture

- Each Lambda has its own `go.mod` and `go.sum`.
- There is no root Go module (`go.mod` removed intentionally).
- SAM uses function-level makefile builds (`BuildMethod: makefile`).
- Lambda source wiring is done via `CodeUri` from the SAM template.

## Key paths

- SAM template: `deploy/sam/kalmel-app-infra/template.yaml`
- SAM config: `deploy/sam/kalmel-app-infra/samconfig.toml`
- SAM app docs: `deploy/sam/kalmel-app-infra/README.md`
- Repo structure reference: `Folder-Structure.md`

## Environment model

The SAM stack is environment-parameterized with:

- `Environment=dev`
- `Environment=prod`

The stack imports shared network/table resources via CloudFormation exports for the selected environment.
