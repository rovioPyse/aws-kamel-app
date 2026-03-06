# kalmel-app-infra

SAM application for Kamel Lambda infrastructure wiring.

## Scope

This stack defines and deploys:

- `LatestStateWriterFunction`
- `AlarmProcessorFunction`

Both functions are built from service-local Makefiles and use the custom runtime pattern (`provided.al2023`, `bootstrap`).

## Source wiring

- `LatestStateWriterFunction` -> `services/lambdas/latest-state-writer`
- `AlarmProcessorFunction` -> `services/lambdas/alarm-processor`

Wiring is defined in `template.yaml` using `CodeUri` relative paths.

## Stack inputs

Required parameter:

- `Environment` (`dev` or `prod`)

## External dependencies

This stack expects existing CloudFormation exports per environment:

- `KamelLambdaSG-${Environment}`
- `KamelPrivateSubnet1-${Environment}`
- `KamelPrivateSubnet2-${Environment}`
- `KamelLatestTable-${Environment}`
- `KamelAlarmsTable-${Environment}`

## Files in this folder

- `template.yaml` - function and policy wiring
- `samconfig.toml` - SAM default configuration
- `Makefile` - local SAM build wrapper
- `README.md` - this document

## Notes

- The stack currently focuses on Lambda resources and does not define API Gateway routes.
- Build output (`.aws-sam/`) is generated artifact content, not source-of-truth.
