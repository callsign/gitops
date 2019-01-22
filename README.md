# GitOps

CLI tool to deploy applications using GitOps.

## Usage

```text
Usage: gitops <command> <argument(s)>
       Valid commands:
            gitops request-deployment <project-url>
            gitops update-configuration <project-path>
```

This command should be run from the directory of the service to deploy.

## Features

* Clone the GitOps project
* Checkout the environment branch following this mapping:

| Service project branch              | GitOps project branch |
|-------------------------------------|-----------------------|
| ^develop$                           | dev                   |
| ^(release\|hotfix)\\/\\S+$          | staging               |
| ^master$                            | prod                  |
| `<custom>` (see Custom Deployments) | `<custom>`            |

* Update the service version in the application Helm *requirements.yaml* (using *build/packages/version* file)
* Update the the application Helm *requirements.lock*
* Copy the service configuration to the GitOps project (*environments/`<environment>`* to *configurations/`<service>`/`<environment>`*)
* Add a service prefix to the service configuration
* Commit and push the GitOps project changes

## Getting Started

```bash
go get github.com/callsign/gitops/...
```

## Building

```bash
go install github.com/callsign/gitops/...
```

## Testing

```bash
go test github.com/callsign/gitops/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Custom Deployments

To deploy branches not covered by the standard branch to environment mappings, please create a `deployments.yaml` file in the application directory with a content like:

```yaml
deployments:
- branch: feature/foo
  environment: bar
```

## License

GitOps is Open Source software released under the [Apache 2.0 license](https://www.apache.org/licenses/LICENSE-2.0.html)
