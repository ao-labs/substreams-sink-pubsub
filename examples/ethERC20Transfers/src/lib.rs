use pb::eth::token::transfers::v1::Transfers;
use pb::sf::substreams::sink::pubsub::v1::{Attribute, Message, Publish};

mod pb;
#[substreams::handlers::map]
fn map_eth_transfers(transfers: Transfers) -> Result<Publish, substreams::errors::Error> {
    let messages = transfers
        .transfers
        .iter()
        .map(|transfer| {
            let from = transfer.from.clone();
            let to = transfer.to.clone();
            Message {
                data: serde_json::to_string(transfer).unwrap().into_bytes(),
                attributes: vec![
                    Attribute {
                        key: "from".to_string(),
                        value: format!{"0x{from}"},
                    },
                    Attribute {
                        key: "to".to_string(),
                        value: format!{"0x{to}"},
                    },
                ],
            }
        })
        .collect();

    let publish = Publish {
        messages,
        };

    Ok(publish)
}
