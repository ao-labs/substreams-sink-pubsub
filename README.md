# substreams-sink-pubsub

Substreams sink for PubSub helps quickly and easily sync blockchain data using Substreams modules to a PubSub topic.
This repository gives all the keys to run a substreams sink for PubSub that provides high-level data on blocks, for any blockchain supported by StreamingFast.  



## Running 

Run the `substreams-sink-pubsub` binary with the following command:

-Mentionned the projectID, the substreams endpoint and the topic on which to publish the data

```bash
substreams-sink-pubsub mainnet.eth.streamingfast.io:443 <project_id> <topic_name> <block_range> <substreams_yaml_path> --module <module_name>
```

### Example 

- Build the substreams 
