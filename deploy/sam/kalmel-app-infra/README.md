# kalmel-app-infra

AWS SAM stack for Kamel Lambdas:

- `LatestStateWriterFunction`
- `AlarmProcessorFunction`

Both functions use custom Go build steps (`BuildMethod: makefile`) from:

- `services/lambdas/latest-state-writer`
- `services/lambdas/alarm-processor`

## Prerequisites

- AWS CLI configured for the target account/region
- AWS SAM CLI installed
- Go installed (used by Lambda Makefiles during `sam build`)
- Docker installed (required for SAM local emulation)

## Files

```text
deploy/sam/kalmel-app-infra/
  Makefile
  README.md
  samconfig.toml
  template.yaml
```

## Stack parameter

This template requires one parameter:

- `Environment`: `dev` or `prod`

## Required CloudFormation exports

The stack imports shared infrastructure values. These exports must exist for the selected environment:

- `KamelLambdaSG-${Environment}`
- `KamelPrivateSubnet1-${Environment}`
- `KamelPrivateSubnet2-${Environment}`
- `KamelLatestTable-${Environment}`
- `KamelAlarmsTable-${Environment}`

## Build

From `deploy/sam/kalmel-app-infra`:

```powershell
make
```

or:

```powershell
sam build
```

## Validate

```powershell
sam validate --lint --parameter-overrides Environment=dev
```

Use `Environment=prod` when validating production configuration.

## Deploy

First-time guided deploy (recommended):

```powershell
sam deploy --guided --parameter-overrides Environment=dev
```

Non-guided deploy (after `samconfig.toml` is initialized):

```powershell
sam deploy --parameter-overrides Environment=dev
```

For production:

```powershell
sam deploy --parameter-overrides Environment=prod
```

Notes:

- `samconfig.toml` currently sets `stack_name = "kalmel-app-infra"`.
- `confirm_changeset = true` is enabled, so deploy asks for confirmation before execution.
- `capabilities = "CAPABILITY_IAM"` is already configured.

## Local invoke

This stack does not define API events. Invoke functions directly with event files:

```powershell
sam local invoke LatestStateWriterFunction --event events/latest-state-writer.json
sam local invoke AlarmProcessorFunction --event events/alarm-processor.json
```

If your handlers need AWS resources (DynamoDB/VPC dependencies), local invoke may require additional mocking or local endpoints.
