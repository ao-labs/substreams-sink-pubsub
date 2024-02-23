# substreams-sink-pubsub

Substreams sink for PubSub helps quickly and easily sync blockchain data using Substreams modules to a PubSub topic.
This repository gives all the keys to run a substreams sink for PubSub that provides high-level data on blocks, for any blockchain supported by StreamingFast.  

## Table of Contents
- [Prerequisites](#prerequisites)
  - [PubSub Creation](#pubsub-creation)
  - [Substreams Creation](#substreams-creation)
- [Installation](#installation)
- [Running](#running)


## Prerequisites

Before sinking any data to a PubSub, make sure to have the following prerequisites:

### PubSub creation

Create a PubSub with a Google cloud projectID associated and a topic on which to publish the data.

### Substreams creation

- Use the `pubsub_substream` provided in the [examples](./examples) directory or create your own substreams.
- Compile the `pubsub_substream` project (or your own substreams):
    ```bash
    cd examples/pubsub_substream
    cargo build --target wasm32-unknown-unknown --release
    ```

Note: If you are creating your own substreams, make sure to create a `map` module with an output type of sf.substreams.sink.pubsub.v1.Publish

## Installation

Install the `substreams-sink-pubsub` binary from source, by running the following command:

```bash
go install ./cmd/substreams-sink-pubsub
```

## Running 

The `substreams-sink-pubsub` binary offers a sink mode of operation.

[//]: # CREATE THIS SECTION


