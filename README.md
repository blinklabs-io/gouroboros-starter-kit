# gOuroboros starter kit

This starter kit is an introduction tutorial on how to use the
[gOuroboros](https://github.com/blinklabs-io/gouroboros) library to communicate
with a Cardano Node using the Ouroboros family of protocols.

The gOuroboros library is a Cardano library written in [Go](https://go.dev). It
allows for communications with a Cardano Node and performing serialization of
Cardano primitives without relying on external tools and binaries, facilitating
the creation of single binary tools with minimal overhead.

## Dev Environment

For running this starter kit you'll need access to a fully synced instance of a
Cardano Node.

If you do not want to install the required components yourself, you can use the
[Demeter.run](https://demeter.run) platform to create a cloud environment with
access to common Cardano infrastructure. The following command will open this
repo in a private, web-based VSCode IDE with access to a shared Cardano Node.

[![Code in Cardano Workspace](https://demeter.run/code/badge.svg)](https://demeter.run/code?repository=https://github.com/blinklabs-io/gouroboros-starter-kit.git&template=golang)

### Cardano Node Access

Since you're running this starter kit from a _Cardano Workspace_, you already
have access to the required infrastructure, such as the Cardano Node.

The network to which you're connected (mainnet, preview, preprod, etc) is
defined by your project's configuration, selected at the moment of creating
your environment.

To simplify the overall process, _Cardano Workspaces_ come already configured
with specific environmental variables that allows you to connect to the node
without extra step. These are the variables relevant for this particular
tutorial:

- `CARDANO_NODE_SOCKET_PATH`: provides the location of the unix socket within
    the workspace where the cardano node is listening.
- `CARDANO_NODE_MAGIC`: the network magic corresponding to the node that is
    connected to the workspace.

## What is Included

This starter kit demonstrates communication with a Cardano Node using either
the Node-to-Client or Node-to-Node operation modes of the Ouroboros network
protocol and provides examples of multiple Ouroboros mini-protocols.

- BlockFetch
- ChainSync
- LocalTxMonitor
- PeerSharing

### BlockFetch

The BlockFetch mini-protocol allows for fetching specific blocks from a remote
Cardano Node using Node-to-Node communication over the network. The code for
this example is in a single `main.go` file under `./cmd/block-fetch`.

By default, it will fetch the first block of the Babbage Era on the Cardano
mainnet from IOG's backbone servers. To change this, use the following
environment variables:

- `BLOCK_FETCH_ADDRESS`: the address:port pair of a remote Cardano Node to
  retrieve block
- `BLOCK_FETCH_HASH`: the block hash to fetch
- `BLOCK_FETCH_NETWORK`: named Cardano network to use to configure network
  magic automatically
- `BLOCK_FETCH_NETWORK_MAGIC`: magic number used to identify a network
- `BLOCK_FETCH_RETURN_CBOR`: return raw CBOR bytes instead of describing text
- `BLOCK_FETCH_SLOT`: the slot in which the block hash was minted

Default:
```bash
go run ./cmd/block-fetch
```

Return raw CBOR bytes:
```bash
BLOCK_FETCH_RETURN_CBOR=true go run ./cmd/block-fetch
```

### ChainSync

The ChainSync mini-protocol allows for syncronization of the blockchain from a
Cardano Node using either Node-to-Node or Node-to-Client communication. This
protocol is more complex, and is split into smaller pieces. The simplest is in
a single `main.go` file under `./cmd/chain-tip`.

For `chain-tip`, the default configuration will communicate over the local
UNIX socket mounted at `/ipc/node.socket` via Node-to-Client ChainSync.

Running the code:
```bash
go run ./cmd/chain-tip
```

#### chain-sync

The full ChainSync implementation that supports blockchain synchronization with
rollback/rollforward handlers. The code is in a single `main.go` file under
`./cmd/chain-sync`.

This command supports both Node-to-Client (NtC, default) and Node-to-Node (NtN)
protocols. It defaults to NtC using a UNIX socket connection, and automatically
switches to NtN when a TCP address is provided.

To customize the connection, use the following environment variables:

- `CARDANO_NODE_SOCKET_PATH`: Path to the node's UNIX socket (default: `/ipc/node.socket`)
- `CARDANO_NODE_ADDRESS`: Remote node address in host:port format (for TCP/NtN connections)
- `CARDANO_NODE_NETWORK`: Named Cardano network to use to configure network magic automatically (mainnet, preview, preprod, etc.)
- `CARDANO_NODE_MAGIC`: Magic number used to identify a network (overrides network-based magic)

Command-line flags:

- `-start-era`: Era to start chain-sync at (genesis, byron, shelley, allegra, mary, alonzo, babbage, conway)
- `-tip`: Start chain-sync at current chain tip (recommended for fully synced nodes)
- `-bulk`: Use bulk chain-sync mode with NtN
- `-range`: Show start/end block of available range

Examples:

Get current tip and sync forward (recommended for production):
```bash
go run ./cmd/chain-sync -tip
```

The script will continuously output blocks as they're synchronized, displaying
era, slot, block number, and block hash for each block.

### LocalTxMonitor

This starter kit demonstrates communication with a Cardano Node using the
Node-to-Client LocalTxMonitor protocol to fetch information about the Node's
mempool contents. It includes a single `main.go` which performs all of the
work, which is located under `cmd/tx-monitor`.

The default configuration will communicate over the local UNIX socket mounted
at `/ipc/node.socket` via Node-to-Client LocalTxMonitor.

```bash
go run ./cmd/tx-monitor
```

The script will output the contents of the Cardano Node's mempool, then exits.

### PeerSharing

This starter kit demonstrates communication with a Cardano Node using the
Node-to-Node PeerSharing protocol to get a list of connected and discovered
peers from the remote Node. It includes a single `main.go` which performs
all of the work, which is located under `cmd/peer-sharing`.

By default, it will fetch the 10 mainnet peers from from IOG's backbone
 servers. To change this, use the following environment variables:

- `PEER_SHARING_ADDRESS`: the address:port pair of a remote Cardano Node to
  retrieve block
- `PEER_SHARING_NETWORK`: named Cardano network to use to configure network
  magic automatically
- `PEER_SHARING_NETWORK_MAGIC`: magic number used to identify a network
- `PEER_SHARING_PEERS`: number of peers to request, max/default 10

```bash
go run ./cmd/peer-sharing
```

The script will output 10 peer addresses from the Node, then exit.

### Ping

The ping command allows to measure the latency of establishing a connection and performing a ChainSync protocol interaction with a Cardano node. It supports both Node-to-Node (TCP/IP) and UNIX socket connections.The code for this example is in a single main.go file under `./cmd/ping`.

By default, it will connect to a local Cardano node via UNIX socket and measure both connection and protocol latency. To customize the connection, use the following environment variables:

- `CARDANO_NODE_SOCKET_PATH`: Path to the node's UNIX socket
- `CARDANO_NODE_ADDRESS`: Remote node address in host:port format (for TCP connections)
- `CARDANO_NODE_NETWORK`: Named Cardano network to use to configure network magic automatically 
- `CARDANO_NODE_MAGIC`: Magic number used to identify a network (overrides network-based magic)

Default:
```bash
go run ./cmd/ping
```
