# OpenStadia

OpenStadia is an open-source project that serves as an alternative to Google Stadia. It allows users to remotely connect
to a powerful computer using the WebRTC protocol.

## Social

- Discord: [OpenStadia](https://discord.gg/tJGeKTEdgj)
- Slack: [OpenStadia](https://openstadia.slack.com/)
- Website: [www.openstadia.com](https://www.openstadia.com/)
- VK: [OpenStadia](https://vk.com/openstadia)

## Deployment Options

There are two deployment options available:

1. Local Deployment: This option does not require a connection to remote servers, but it can only be used within a local
   network or with a static public IP address (router configuration may be required).
2. Hub Server Deployment: This option involves using a hub server to establish connections between users. The open
   version of the hub server can be accessed at https://github.com/openstadia/openstadia-hub.

## Supported OS

OpenStadia currently supports the following operating system:

- Linux (Tested on Ubuntu 22.04)
- Windows (WIP)

## Support Table

|                 | Linux           | Windows         | macOS           | External Device |
|-----------------|-----------------|-----------------|-----------------|-----------------|
| Mouse           | :green_circle:  | :yellow_circle: | :purple_circle: | :yellow_circle: |
| Keyboard        | :yellow_circle: | :yellow_circle: | :purple_circle: | :yellow_circle: |
| Gamepad         | :green_circle:  | :yellow_circle: | :purple_circle: | :yellow_circle: |
| Virtual Display | :green_circle:  | :black_circle:  | :red_circle:    | :black_circle:  |
| Container       | :purple_circle: | :black_circle:  | :red_circle:    | :black_circle:  |

- :green_circle: - supported
- :purple_circle: - support is planned
- :red_circle: - support is NOT planned
- :yellow_circle: - work in progress
- :black_circle: - support is NOT possible

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

### Common

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

### Linux

1. Install build requirements

```shell
sudo apt-get install \
  libx11-dev \
  libxext-dev \
  libvpx-dev \
  libx11-dev \
  xorg-dev \
  libxtst-dev
```

2. Build the OpenStadia application using the following command:

```shell
go build
```

### Windows

1. Install MSYS2 on your system by following the official documentation: [Installing MSYS2](https://www.msys2.org/).
2. Open MSYS2 UCRT64 console
3. Install build requirements

```shell
pacman -S mingw-w64-ucrt-x86_64-toolchain
pacman -S mingw-w64-ucrt-x86_64-libvpx
pacman -S mingw-w64-ucrt-x86_64-libpng
```

4. Set environment variables and navigate to the project directory

```shell
export PATH=$PATH:/c/Go/bin
cd /c/path-to-openstadia/openstadia/
```

5. Build the OpenStadia application using the following command:

```shell
go build -v -x
```

### Launching

1. Run the OpenStadia application:

```shell
./openstadia
```

2. Connect to OpenStadia using a supported web browser or client application. http://127.0.0.1:8080

3. Enjoy remote access to a powerful computer for gaming or other resource-intensive tasks.

## Contributing

Contributions to OpenStadia are welcome! If you encounter any issues or have ideas for improvements, please submit them
via the GitHub issue tracker. You can also contribute by submitting pull requests with bug fixes or new features.

Please refer to the CONTRIBUTING.md file in the OpenStadia repository for more information on how to contribute.

## License

OpenStadia is released under the [MIT License](https://opensource.org/licenses/MIT). Please see the LICENSE file for
more details.