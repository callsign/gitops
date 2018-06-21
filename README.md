GitOps
======

CLI tool to deploy applications using GitOps.

Usage
-----

```
Usage: gitops <command> <argument(s)>
       Valid commands:
            gitops request-deployment <project-url>
            gitops update-configuration <environment> <project-path>
```

This command should be run from the directory of the service to deploy.

Features
--------

* Clone the GitOps project
* Checkout the environment branch following this mapping:

| Service project branch     | GitOps project branch |
|----------------------------|-----------------------|
| ^develop$                  | verify                |
| ^(release\|hotfix)\\/\\S+$ | staging               |
| ^master$                   | production            |

* Update the service version in the application Helm *requirements.yaml* (using *build/packages/version* file)
* Update the the application Helm *requirements.lock*
* Copy the service configuration to the GitOps project (*environments/`<environment>`* to *configurations/`<service>`*)
* Add a service prefix to the service configuration
* Commit and push the GitOps project changes

Getting Started
---------------
```
go get github.com/callsign/gitops/...
```

Building
--------
```
go install github.com/callsign/gitops/...
```

Testing
-------
```
go test github.com/callsign/gitops/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
