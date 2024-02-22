// @generated
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Publish {
    #[prost(message, repeated, tag="1")]
    pub messages: ::prost::alloc::vec::Vec<Message>,
}
impl ::prost::Name for Publish {
const NAME: &'static str = "Publish";
const PACKAGE: &'static str = "pubsub.v1";
fn full_name() -> ::prost::alloc::string::String {
                ::prost::alloc::format!("pubsub.v1.{}", Self::NAME)
            }}
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Message {
    #[prost(bytes="vec", tag="1")]
    pub data: ::prost::alloc::vec::Vec<u8>,
    #[prost(message, repeated, tag="2")]
    pub attributes: ::prost::alloc::vec::Vec<Attribute>,
    #[prost(string, tag="3")]
    pub ordering_key: ::prost::alloc::string::String,
}
impl ::prost::Name for Message {
const NAME: &'static str = "Message";
const PACKAGE: &'static str = "pubsub.v1";
fn full_name() -> ::prost::alloc::string::String {
                ::prost::alloc::format!("pubsub.v1.{}", Self::NAME)
            }}
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Attribute {
    #[prost(string, tag="1")]
    pub key: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub value: ::prost::alloc::string::String,
}
impl ::prost::Name for Attribute {
const NAME: &'static str = "Attribute";
const PACKAGE: &'static str = "pubsub.v1";
fn full_name() -> ::prost::alloc::string::String {
                ::prost::alloc::format!("pubsub.v1.{}", Self::NAME)
            }}
// @@protoc_insertion_point(module)
