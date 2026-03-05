# Folder-Structure.md

This document defines a general-purpose repository structure for a backend "apps repo" that contains:

- AWS SAM application stack (for Lambdas and app wiring)
- Multiple Lambda services
- One or more ECS/Fargate services (Dockerized)
- Shared libraries used across services

It is intentionally not specific to any single product. Use it as a standard template for similar cloud backends.

## Goals of this structure

### 1) Make "what runs" obvious

A repo should clearly show:

- what deploys as a Lambda
- what deploys as an ECS container
- what is shared library code
- what is deployment wiring

This reduces onboarding time and makes CI/CD automation straightforward.

### 2) Enforce modularity by default

Services should not casually depend on other services' implementation details.

We want:

- shared code -> `libs/`
- service-private code -> each service's `internal/`

This prevents tight coupling and spaghetti dependencies as the repo grows.

### 3) Enable independent builds

Each deployable unit has its own entrypoint, so you can build/test:

- one Lambda
- one ECS worker

without touching others.

### 4) Keep the repo scalable

The same layout works for:

- 2 services today
- 25 services later

with minimal refactoring.

## Standard repository layout

```text
apps-repo/
  README.md
  Folder-Structure.md               # this document
  .gitignore
  Makefile                          # or Taskfile.yml (optional)
  go.mod
  go.sum

  deploy/                           # Deployment wiring (can support multiple deployment methods/apps)
    sam/                            # SAM-based apps (current)
      <app-name-1>/
        template.yaml               # functions, schedules, permissions, wiring
        samconfig.toml              # env profiles (dev/prod)
        params/                     # env params (optional)
          dev.json
          prod.json
      <app-name-2>/
        template.yaml
        samconfig.toml
    cdk/                            # Optional future deployment method
      <app-name-1>/...
    terraform/                      # Optional future deployment method
      <app-name-1>/...

  libs/                             # Shared libraries (reusable, stable APIs)
    <lib-name-1>/
      *.go
    <lib-name-2>/
      *.go
    <lib-name-3>/
      *.go

  services/                         # Deployable services (no cross-imports)
    lambdas/
      <lambda-service-name-1>/
        cmd/
          handler/
            main.go                 # executable entrypoint (package main)
        internal/                   # PRIVATE to this service (compiler-enforced)
          ...
        config/
          config.go                 # config loading for this service (optional)
        test/
          *_test.go                 # unit tests (optional)

      <lambda-service-name-2>/
        cmd/handler/main.go
        internal/...
        config/...
        test/...

    ecs/
      <ecs-service-name-1>/         # one ECS app / worker / job
        cmd/
          worker/
            main.go                 # executable entrypoint (package main)
        internal/                   # PRIVATE to this service (keep algorithms/logic here)
          ...
        deploy/
          Dockerfile                # container build instructions
        config/
          config.go                 # config loading (optional)
        test/
          *_test.go                 # unit tests (optional)

      <ecs-service-name-2>/
        cmd/worker/main.go
        internal/...
        deploy/Dockerfile

  scripts/                          # CLI wrappers for dev/CI (optional)
    build.sh
    test.sh
    deploy-dev.sh
    deploy-prod.sh

  .github/
    workflows/
      ci.yml
      deploy-dev.yml
      deploy-prod.yml
```

## Why each directory exists

### deploy/ - deployment wiring, not service code

Purpose: Define app-level resources and wiring across one or more deployable apps:

- Lambda functions, triggers, permissions
- EventBridge schedules
- ECS task definitions + schedules
- IAM roles/policies
- (Any app-level resources that are not "core infra")

Why it matters:
Keeping deployment wiring separate:

- prevents mixing infra definitions into service code
- makes it easy for CI to run deployment pipelines consistently (SAM now, other methods later)
- allows multiple services to be deployed as one app stack or multiple app stacks

Recommended organization:

- group by deployment method first (`deploy/sam`, `deploy/cdk`, `deploy/terraform`, etc.)
- then group by app name under each method (`deploy/sam/<app-name>/...`)
- keep all deployment-only artifacts here (templates, params, env/profiles, stack configs)
- for a single-app repo, `deploy/template.yaml` is acceptable initially; move to method/app subfolders when more apps/methods are introduced

Rule: `deploy/` should reference service folders via paths (e.g., `CodeUri`) but must not contain service logic.

### services/ - deployable compute units

Purpose: Each folder under `services/` corresponds to something you deploy and operate independently:

- a Lambda function service
- an ECS container service/job

Why it matters:
This is the most important boundary in the repo. It answers:

- "What are our running services?"
- "Where is the entrypoint?"
- "How do I build/test just this one?"

Rule: `services/` is for things that run. It should not contain shared libs.

### libs/ - shared libraries only

Purpose: Code that is meant to be reused across multiple services:

- domain types (schemas, validation helpers)
- AWS SDK helpers/wrappers
- observability (logging/tracing)
- shared clients/utilities

Why it matters:
This prevents duplication and makes changes safer. If multiple services need a thing, it belongs here.

Rule: `libs/` must not import anything under `services/`.

## Why we use cmd/ and internal/

### cmd/ - explicit entrypoints

In Go, a build target is identified by a folder containing:

- `package main`
- `func main()`

Putting entrypoints under `cmd/<name>/main.go` makes it obvious:

- which directories compile into executables
- which directories are library code

This is a Go community convention that scales extremely well.

Rule: Every service must have exactly one entrypoint under `cmd/`.

### internal/ - real enforcement, not just a convention

Go treats `internal/` specially:

Packages inside `internal/` can only be imported by code within the parent directory tree.

That means each service can keep implementation details private, and other services can't depend on them.

This is critical for a scalable monorepo because it provides a compiler-enforced boundary.

Rule: Anything service-specific goes into that service's `internal/`.
Anything reusable goes into `libs/`.

## Naming standards

### Service names

Use purpose-based names, not tech-based names:

- Yes: `latest-state-writer`, `alarm-processor`, `hourly-processor`
- No: `lambda1`, `go-service`, `worker2`

This helps operational clarity in AWS, logs, dashboards, and on-call.

### Library names

Use noun-based names that describe what the library provides:

`telemetry`, `awsx`, `observability`, `config`, `httpx`

## Build outputs and where they go

### Lambdas (SAM)

- Build artifacts go to: `.aws-sam/build/` by default
- You can choose a common build directory: `sam build --build-dir .build/sam`

Rule: build artifacts are never committed.

### ECS (Docker)

- The artifact is the Docker image
- No build directory is required unless you explicitly output binaries

Rule: ECS services must have a `deploy/Dockerfile` in their own folder.

## Import rules (must-follow)

Services can import only:

- `libs/...`
- their own `internal/...`

Services must not import:

- another service folder under `services/...`

`libs/` must never import `services/`

These three rules preserve modularity and keep the repo maintainable.

## How to add a new service (the standardized workflow)

### Add a new Lambda service

1. Create folder: `services/lambdas/<name>/`
2. Add entrypoint: `cmd/handler/main.go`
3. Put logic in: `internal/`
4. Add wiring in the target app under `deploy/` (for example `deploy/sam/<app-name>/template.yaml`)

### Add a new ECS service

1. Create folder: `services/ecs/<name>/`
2. Add entrypoint: `cmd/worker/main.go`
3. Put logic in: `internal/`
4. Add Dockerfile: `deploy/Dockerfile`
5. Add wiring in the target app under `deploy/` (for example `deploy/sam/<app-name>/template.yaml`) for task defs/schedules

## Why we keep this format consistent

Consistency gives you three long-term benefits:

### 1) It's automation-friendly

Scripts and CI can assume:

- where entrypoints exist
- where Dockerfiles exist
- where SAM template lives
- which folders are shared libs

### 2) It reduces cognitive load

Developers don't need to "rediscover" structure every time.

### 3) It prevents architectural drift

A standard structure makes it harder to accidentally:

- duplicate shared logic in services
- couple services to each other
- hide deployable code in random folders

## Optional additions (if/when needed)

- `Taskfile.yml` (Task) instead of `Makefile`
- `docs/` for ADRs and runbooks
- `tools/` for linters, codegen utilities
- CI checks to enforce no cross-service imports
- additional deployment methods under `deploy/` (CDK/Terraform/other)

## Summary

This folder structure is designed to be:

- Clear: what runs vs what's shared vs what's deployment wiring
- Modular: shared libs only in `libs/`, private code in each service's `internal/`
- Scalable: handles many services without becoming chaotic
- CI-friendly: predictable build/deploy patterns with support for multiple deployment methods over time

Use this structure for any backend apps repo that contains serverless + containerized workloads.
