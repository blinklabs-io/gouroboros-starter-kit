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

This starter kit demonstrates communication with a Cardano Node using the
Node-to-Client LocalTxMonitor protocol to fetch information about the Node's
mempool contents. It includes a single `main.go` which performs all of the
work.

## Running the Code

Running the code is simple.

```bash
go run .
```

The script will output the contents of the Cardano Node's mempool, then exit.
