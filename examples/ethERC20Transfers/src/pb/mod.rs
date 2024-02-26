pub mod eth {
    pub mod token {
        pub mod transfers {
            // @@protoc_insertion_point(attribute:eth.token.transfers.v1)
            pub mod v1 {
                include!("eth.token.transfers.v1.rs");
                // @@protoc_insertion_point(eth.token.transfers.v1)
            }
        }
    }
}

pub mod sf {
    pub mod substreams {
        pub mod sink {
            pub mod pubsub {
                // @@protoc_insertion_point(attribute:sf.substreams.sink.pubsub.v1)
                pub mod v1 {
                    include!("sf.substreams.sink.pubsub.v1.rs");
                    // @@protoc_insertion_point(sf.substreams.sink.pubsub.v1)
                }
            }
        }
        pub mod rpc {
            // @@protoc_insertion_point(attribute:sf.substreams.rpc.v2)
            pub mod v2 {
                include!("sf.substreams.rpc.v2.rs");
                // @@protoc_insertion_point(sf.substreams.rpc.v2)
            }
        }
        // @@protoc_insertion_point(attribute:sf.substreams.v1)
        pub mod v1 {
            include!("sf.substreams.v1.rs");
            // @@protoc_insertion_point(sf.substreams.v1)
        }
    }
}
pub mod substreams {
    pub mod sink {
        pub mod files {
            // @@protoc_insertion_point(attribute:substreams.sink.files.v1)
            pub mod v1 {
                include!("substreams.sink.files.v1.rs");
                // @@protoc_insertion_point(substreams.sink.files.v1)
            }
        }
    }
}

