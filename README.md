# OpenStadia

OpenStadia is an open-source project that serves as an alternative to Google Stadia. It allows users to remotely connect
to a powerful computer using the WebRTC protocol.

## Deployment Options

There are two deployment options available:

1. Local Deployment: This option does not require a connection to remote servers, but it can only be used within a local
   network or with a static public IP address (router configuration may be required).
2. Hub Server Deployment: This option involves using a hub server to establish connections between users. The open
   version of the hub server can be accessed at https://github.com/openstadia/openstadia-hub.

## Supported OS

OpenStadia currently supports the following operating system:

- Linux (Tested on Ubuntu 22.04)

## Configuration

Configuration settings for OpenStadia are specified in the `openstadia.yaml` file, which should be located in the same
directory as the executable file. The following parameters can be configured:

- hub: The URL address of the hub server. This parameter is not used in local deployments. Please note that there is
  currently no public hub available.
- local: Config for local deployment
- apps: A list of applications that users can launch from their remote devices. Currently, all applications open
  on a virtual display to avoid disrupting the system's operation.

Example configuration

```yaml
hub:
  addr: "hub.openstadia.com"
  token: "my-awesome-secret"

local:
  addr: "0.0.0.0:8080"

apps:
  - name: "my-awesome-app"
    command: [ "/home/user/my-awesome-app" ]
    width: 1920
    height: 1080

```

## Getting Started

To get started with OpenStadia, follow these steps:

1. Clone the OpenStadia repository.
2. Configure the openstadia.yaml file with the desired settings.
3. Choose the deployment option based on your requirements:
    - For local deployment, ensure you have a powerful computer with a static public IP address or configure your router
      accordingly.
    - For hub server deployment, set up the hub server using the instructions provided
      at https://github.com/openstadia/openstadia-hub.
4. Build or run the OpenStadia application.
5. Connect to OpenStadia using a supported web browser or client application.
6. Enjoy remote access to a powerful computer for gaming or other resource-intensive tasks.

## Build Instructions

To build and run OpenStadia using the Go programming language, please follow these instructions:

1. Install Go on your system by following the official documentation: [Installing Go](https://golang.org/doc/install).
2. Clone the OpenStadia repository using the following command:

```shell
git clone https://github.com/openstadia/openstadia.git
```

3. Navigate to the project directory:

```shell
cd openstadia
```

4. Install the project dependencies using the following command:

```shell
go get
```

5. Build the OpenStadia application using the following command:

```shell
go build
```

6. Run the OpenStadia application:

```shell
./openstadia
```

7. Connect to OpenStadia using a supported web browser or client application. http://127.0.0.1:8080

8. Enjoy remote access to a powerful computer for gaming or other resource-intensive tasks.

## Contributing

Contributions to OpenStadia are welcome! If you encounter any issues or have ideas for improvements, please submit them
via the GitHub issue tracker. You can also contribute by submitting pull requests with bug fixes or new features.

Please refer to the CONTRIBUTING.md file in the OpenStadia repository for more information on how to contribute.

## License

OpenStadia is released under the [MIT License](https://opensource.org/licenses/MIT). Please see the LICENSE file for
more details.