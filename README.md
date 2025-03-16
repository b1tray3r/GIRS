# Gitea Issue - Redmine - Syncer

The purpose of this repository is relatively easy:

If an issue is created in the infra repository, it should be synced to redmine into the devops project.
To have a reference to the synced issue, the service will comment the issue URL to the gitea issue.

As soon as the issue is closed in gitea, the issue should be closed in redmine, by using the commented issue URL.

This syncer only reacts to the opened and closed events.

## Setup

1. Install dependencies

```sh
make install
```

## Build

1. Build the image

```sh
make image
```

2. Run the application via docker

```shell
make run
```

3. Visit `{url}/swagger` for tests and API documentation

## Develop

1. Start the watcher in the project root

```shell
air
```

### Details

#### index.html

This contains the swagger-ui code and requires access to [unpkg](https://unpkg.com/) servers.

The file is requesting the generated `/swagger/openapi3.json` file to build the GUI.

#### .air.toml

Using air the complete build pipeline is executed on the fly.
When building the final application you have to execute all `pre_cmd` entries in the `.air.toml` configuration.

Please note that all `*.gen.go` files are ignored by air.
It is important for air to exclude this files, because generating these file(s) will cause air to loop indefinetly.