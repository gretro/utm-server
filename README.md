# UTM Server
This server exposes the UTM command line API as REST endpoints. It is useful to allow remote machines (or Docker containers) to create virtual machines on the host for development purposes.

This is useful when you need UTM to act as a IaaS test provider.

> This server has been tested to work against UTM 4.4.5

## Useful links

* [UTM](https://mac.getutm.app/)
* [UTM Command Line](https://docs.getutm.app/scripting/scripting/#command-line-interface)


## Requirements
* Go 1.21
* macOS
* UTM 4.4.5+
* Make

## Launching the server
The server launches with sensitive defaults on port 8080. 

Depending on your architecture, launch the server on the host where UTM is located by running the appropriate executable.

  * arm64: `utm_server_arm64`
  * amd64: `utm_server_amd64`

You can change the configuration via some environment variables.

| Name        | Description                                                | Default                 |
|-------------|------------------------------------------------------------|-------------------------|
| `UTM_PATH`  | Specifies the path where to find  `UTM.app`                | `/Applications/UTM.app` |
| `HTTP_HOST` | Specifies on which address the server listens for requests | `0.0.0.0`               |
| `HTTP_PORT` | Specifies on which port the server listens for requests    | 8080                    |

## Build project
You can build the project for both arm4 and amd64 using the following make command:
```bash
make build
```
