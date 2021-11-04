# Dynamic Linux Tunneling

This software creates Gretap Tunnels over FOU for Dynamic client endpoints.  
It also works behind NAT444 (CGN-LSN). You can use low end devices for connecting two sites with a lightweight network tunnel mechanism that is supported by the Linux kernel.

Before starting using this software, please be sure your kernel supports FOU encap.  
If your kernel has a FOU support but not enabled, you can enable it with `modprobe fou` command.

## Server Side

System is requires 2 ports for testing, identifying and network traffic.
First default port 65200 UDP is used for detecting client IP and port pair for tunneling, 65200 TCP is used for when the tunnel interfaces are ready and up, TCP keepalive packets checks the tunnel health.
Second port used for handling incoming tunneling packets. By default FOU port is 65201.

### Server configuration

The application is configured with environment variables.

```bash
LOG_LEVEL=2     # Debug level 1, Info level 2, Error level 3
LOG_FILE=-      # Save output to file, default is - for stdout
IP=[::]         # Listen address for first port
PORT=65200      # Client test and identification. It listens to both TCP and UDP.
FOU_PORT=65201  # Incoming tunneling packet destination.
SCRIPT_FILE=    # Execute when the new client established
CLIENT_LIST=    # Client configurations file
```

For client verification and interface MTU configuration, client configuration for the server is stored in a json file.

```json
[
    {"ClientName":"client1","clientKey":"KUnqdrF54YrHxDQK", "MTU":1460}
]
```

With defining `SCRIPT_FILE` environment variable, when the new client connection is established, the system will execute a defined file with below environment variables.

```bash
MODE=server              # Current server mode
CLIENT_NAME=client1      # Client name
REMOTE_ADDR=203.0.113.56 # Client connection IP
REMOTE_PORT=32284        # Client connection UDP port
MTU=1448                 # MTU value that is defined Client list
INTERFACE=dyn1           # Interface name for this client
FOU_PORT=65201           # Local fou port for listening incoming packets
```

Example script is available at `/scripts/server.sh`.

## Client Side

To automate configure tunnel endpoints at client side, you can execute this software as client mode. You can switch client mode with setting the first argument as client or set the environment variable `MODE=client` .

### Client Configuration

You can again configure the application with below environment variables.

```bash
LOG_LEVEL=2   # Debug level 1, Info level 2, Error level 3
LOG_FILE=-    # Save output to file, default is - for stdout
SERVER_LIST   # Server list json file path
SCRIPT_FILE   # Execute a file when the connection established
MODE=client   #  Set program mode as client
```

Server list file contains client name and remote server infos such as IP and PORT and also each connection has a key for validating the client.

```json
{
    "ClientName": "client1",
    "servers": [
        {
            "remoteAddr": "203.0.113.80",
            "remotePort": 65200,
            "clientKey": "client1keyForHas",
            "MTU": 1460
        }
    ]
}
```

Below arguments are passed to the script file when the system establishes the connection.

```bash
MODE=client                                     # Current application mode
REMOTE_ADDR=203.0.113.80                        # Server IP address
REMOTE_PORT=65200                               # Server Managment port
MTU=1460                                        # Tunnel MTU
REMOTE_LOCAL_IPV6=fe80::4c6:1dff:fea8:fd03/64   # Link Local IPv6 addr for remote server endpoint
WHOAMI_IP=203.0.113.56                          # Client IP address detected by server
WHOAMI_PORT=32068                               # Client Incoming port detected by server
INTERFACE=dyn1                                  # Created interface for this connection
FOU_PORT=48398                                  # Client Incoming port before nat. 
REMOTE_FOU_PORT=65201                           # Server Inbound Port
```

You can find example script at `/scripts/client.sh`.
