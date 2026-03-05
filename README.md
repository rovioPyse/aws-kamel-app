# aws-kamel-app

Backend apps repo with two Go Lambdas wired through AWS SAM:

- `latest-state-writer`
- `alarm-processor`

## Prerequisites

- Go (1.23+)
- AWS SAM CLI
- Docker (required for `sam local invoke`)

## Repository Layout

- `services/lambdas/latest-state-writer` - Lambda source
- `services/lambdas/alarm-processor` - Lambda source
- `deploy/template.yaml` - SAM stack wiring

## Run Locally

1. Create sample event files:

```powershell
New-Item -ItemType Directory -Force events | Out-Null
'{}' | Set-Content events/latest-state-writer.json
'{"id":"1","source":"demo","detail-type":"Scheduled Event","detail":{}}' | Set-Content events/alarm-processor.json
```

2. Invoke each Lambda locally:

```powershell
sam local invoke LatestStateWriterFunction --template deploy/template.yaml --event events/latest-state-writer.json
sam local invoke AlarmProcessorFunction --template deploy/template.yaml --event events/alarm-processor.json
```

## Test

Run all Go tests:

```powershell
go test ./...
```

## Build

Build the full SAM application:

```powershell
sam build --template deploy/template.yaml
```

Build only one Lambda:

```powershell
sam build LatestStateWriterFunction --template deploy/template.yaml
sam build AlarmProcessorFunction --template deploy/template.yaml
```

Artifacts are written to `.aws-sam/build/` (or your configured SAM build directory).
