# Dynamic Linux Tunneling

This software creates Gretap Tunnels over FOU for Dynamic client endpoints.  
It also works behind NAT444 (CGN-LSN).


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
MODE=server				# Current server mode
CLIENT_NAME=aaf			# Client name
REMOTE_ADDR=203.0.113.56 # Client connection IP
REMOTE_PORT=32284 		# Client connection UDP port
MTU=1448 				# MTU value that is defined Client list
INTERFACE=dyn1 			# Interface name for this client
FOU_PORT=65201 			# Local fou port for listening incoming packets
```
