# spout

`spout` is a lightweight log viewer for cloud service logs.

![spout](https://user-images.githubusercontent.com/605953/184556440-262d8a1b-cb14-47d8-b9db-357c5bfba568.gif)

## Features

- Quick filtering by [jq](https://stedolan.github.io/jq/)
- Invoke from command line with pre-defined options
- Support multiple log sources
  - Google Cloud Logging and JSON formatted local file for now

## Installation

### From releases

Download a binary for your platform from https://github.com/m-mizutani/spout/releases

### Build

`go install` can not be used to install because of requiring `npm`. You can build `spout` by following steps.

```sh
$ git clone https://github.com/m-mizutani/spout.git
$ cd spout
$ cd frontend && npm install && npm run export && cd ..
$ go build .
```

## Usage

### Google Cloud Logging

Example of command line:

```bash
$ spout gcp --project <your-project-id> -d 10m --filter "resource.type=k8s_container"
```

Options:

```
   --addr value, -a value       Server address for browser mode (default: "127.0.0.1:3280") [$SPOUT_ADDR]
   --base-time value, -t value  Base time [$SPOUT_BASE_TIME]
   --duration value, -d value   Duration, e.g. 10m, 30s (default: "10m") [$SPOUT_DURATION]
   --filter value, -f value     Google Cloud Logging filter  (accepts multiple inputs) [$SPOUT_GCP_FILTER]
   --limit value                Limit of fetching log (default: 1000) [$SPOUT_GCP_LIMIT]
   --mode value, -m value       Run mode [console|browser] (default: "browser") [$SPOUT_MODE]
   --project value, -p value    Google Cloud Project ID [$SPOUT_GCP_PROJECT]
   --range value, -r value      Range type [before|after|around] (default: "before") [$SPOUT_RANGE]
```

### Preset options

Save `.spout.toml` like following at **current working directory**.

```toml
[stg]
command = "gcp"
options = [
    "--limit", "100",
    "--project", "your-staging-service",
    "--filter", 'resource.type=k8s_container labels."k8s-pod/services_ubie_app/app"="your_app"',
]

[prd]
command = "gcp"
options = [
    "--limit", "100",
    "--project", "your-production-service",
    "--filter", 'resource.type=k8s_container labels."k8s-pod/services_ubie_app/app"="your_app"',
]
```

Then, you can call predefined options as following:

```bash
$ spout call prd # means 'gcp --limit 100 --project your-production-service ...'
$ spout call prd --filter "user_20352904853" # Append '--filter "user_20352904853"' to existing prd options
```

## License

Apache version 2.0
