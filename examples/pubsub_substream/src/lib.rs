use prost::Message as protoMessage;
use prost_types::Any;
use substreams::log_info;
use substreams::pb::sf::substreams::v1::Clock;

use pb::pubsub::v1::{Message, PublishOperation, PublishOperations};

mod pb;

#[substreams::handlers::map]
fn map_clocks(clock: Clock) -> Result<PublishOperations, substreams::errors::Error> {
    let publish_operations = PublishOperations {
        publish_operations: vec![PublishOperation {
            topic_id: "topic".to_string(),
            ordering_key: "key".to_string(),
            message: Some(Message {
                data: clock.number.to_be_bytes().to_vec(),
                data_any: Some(Any {
                    type_url: "type.googleapis.com/sf.substreams.v1.Clock".to_string(),
                    value: clock.encode_to_vec(),
                }),
                attributes: vec![],
            }),
        }],
    };

    log_info!(
        "Publishing: {:?}",
        publish_operations.publish_operations[0]
            .clone()
            .message
            .unwrap()
            .data_any
            .unwrap()
            .type_url
    );
    Ok(publish_operations)
}
