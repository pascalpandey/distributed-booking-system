# SC4051 Distributed Systems - Facility Booking System Client

This is a server implementation for the SC4051 Distributed Systems Course Project: Design and Implementation of A Distributed Facility Booking System at Nanyang Technological University (NTU).

The server accepts UDP-based protocol requests from the client and is responsible for handling availability queries, making bookings, and managing facility reservations.

## Project Overview

This server implementation is part of a distributed facility booking system that demonstrates concepts such as:

- Client-server Architecture
- Network communication via UDP

## Project Structure

```
sc4051-server/
│
├── main.go                        # Main entry point for the server application
│
├── state/
│   ├── booking_time.go            # Data class for booking time
│   ├── booking.go                 # Data class for each booking
│   ├── facility_state.go          # Data class for each facility's bookings
│   ├── observers.go               # Data class for monitor functionality
│   └── state.go                   # Main data class for state
|
├── server/
│   └── server.go                  # Utility function to init UDP server
│
├── serializer/
│   ├── booking.go                 # Serialize booking messages
│   ├── notify.go                  # Serialize monitor callback messages
│   └── operations.go              # Serialize other operations (except monitor) messages
|
├── handler/
│   └── handler.go                 # Data class to handle deserialized messages
|
├── deserializer/
│   ├── booking.go                 # Deserialize booking messages
│   ├── message.go                 # Deserialize initial client messages
│   └── operations.go              # Deserialize operations messages
|
├── client/
│   └── client.go                  # Data class to send messages back to the client or other servers
|
├── router/
│   ├── main.go                    # Entrypoint of router for reverse proxy setup
│   ├── forward.go                 # Utilities to support forwarding messages to handler servers
│   └── facility_type.go           # Utilities to extract facility type
|
├── cluster/
│   ├── incoming.go                # Handlers to handle cluster consensus messages
│   ├── outgoing.go                # Handlers to send cluster consensus messages
|   ├── state.go                   # State stored inside each server
│   └── utilities.go               # Helper functions for cluster oprations
|
└── README.md                      # This file
```

## Installation Guide

### Prerequisites

- Go 1.23.2 or higher

### Setting Up Dependencies

1. Clone the repository:

   ```bash
   git clone https://github.com/pascalpandey/sc4051-server
   cd sc4051-server
   ```

2. Install Dependencies

   ```bash
   go mod tidy
   ```

3. To run unit tests (implemented for `deserializer`, `serializer`, and `state` packages):

   ```bash
   go test ./serializer ./deserializer ./state
   ```

## Running the Server Application

### Basic Usage

To start the server on port 9000:

```bash
go run main.go
```

### Simulating Unreliable Network Conditions

To simulate unreliable network we have implemented a `drop` flag. When this flag is used every other packet is dropped (i.e. client sends packets 1,2,3,4, the server will process packets 1 and 3 and drop their replies while handling packets 2 and 4 normally).

```bash
go run main.go -drop=true
```

### Supporting At Least Once Semantics

To handle the above conditions the client would retransmit messages on timeout. This however requires the server to handle repeated messages on non indempotent operations. To support this we have implemented a `cache` flag to prevent multiple invocations of non indempotent operations.

```bash
go run main.go -drop=true -cache=true
```

### Reverse Proxy Setup

To improve scalability and reliability, we introduce a new setup. Start 3 servers handling TR, lab, and theatre facilities each. Then route requests into the correct server. This requires running a router which sits in between the client and the partitioned servers.

```bash
# timeout is in seconds
go run ./router/ -port=9000 -timeout="0.5" -labAddr=localhost:8000 -trAddr=localhost:8001 -theatreAddr=localhost:8002
go run main.go -port=8000 -drop=true -cache=true
go run main.go -port=8001 -drop=true -cache=true
go run main.go -port=8002 -drop=true -cache=true
```

### Cluster Setup

In the reverse proxy setup in each partition of the data, when a server goes down the partition won't be able to respond to requests. We have implemented a consensus based cluster setup based on Raft to combat this. It includes leader election, replication, and persistence to disk in order to tolerate floor((n-1)/2) server failures.

Technically we can run separate clusters on each data partition for the best performance and scalability, but it can be a bit messy to demonstrate, so we simulate only a single cluster that holds data from all venue types.

```bash
mkdir states # directory to save states in disk
go run ./router/ -port=9000 -timeout="0.6" -labAddr="localhost:8000,localhost:8001,localhost:8002,localhost:8003,localhost:8004" -trAddr="localhost:8000,localhost:8001,localhost:8002,localhost:8003,localhost:8004" -theatreAddr="localhost:8000,localhost:8001,localhost:8002,localhost:8003,localhost:8004"
go run main.go -port=8000 -cluster="localhost:8000,localhost:8001,localhost:8002,localhost:8003,localhost:8004" -cache=true
go run main.go -port=8001 -cluster="localhost:8000,localhost:8001,localhost:8002,localhost:8003,localhost:8004" -cache=true
go run main.go -port=8002 -cluster="localhost:8000,localhost:8001,localhost:8002,localhost:8003,localhost:8004" -cache=true
go run main.go -port=8003 -cluster="localhost:8000,localhost:8001,localhost:8002,localhost:8003,localhost:8004" -cache=true
go run main.go -port=8004 -cluster="localhost:8000,localhost:8001,localhost:8002,localhost:8003,localhost:8004" -cache=true
```

## Functionality

The server provides these functions:

- **View Bookings**: Display current facility bookings in a table format
- **Query Availability**: Check if a facility is available for a specific time range
- **Make Booking**: Book a facility for a specific time range
- **Offset Booking**: Move a booking to a different time
- **Extend Booking**: Increase the duration of a booking
- **Cancel Booking**: Cancel an existing booking
- **Monitor Venues**: Watch for changes to a facility's booking status

## UDP Communication Protocol

The client will communicate with the server using a UDP-based protocol. Each message follows a specific format:

- **Query Availability**: `{requestId},QUERY,{venueName},{start},{end}`
- **Book Venue**: `{requestId},BOOK,{venueName},{start},{end}`
- **Offset Booking**: `{requestId},OFFSET,{confirmationId},{offset}`
- **Extend Booking**: `{requestId},EXTEND,{confirmationId},{extension}`
- **Cancel Booking**: `{requestId},CANCEL,{confirmationId}`
- **Monitor Venue**: `{requestId},MONITOR,{venueName},{duration}`

Client requests are parsed according to specific patterns handled by the deserializers and responses are formatted and serialized according to rules specified in the serializer package to be sent back to the client.


