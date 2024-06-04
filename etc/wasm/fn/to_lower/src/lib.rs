#![allow(dead_code)]

pub use camel_wasm_sdk::Message;

// *****************************************************************************
//
// Functions
//
// ******************************************************************************

#[cfg_attr(all(target_arch = "wasm32"), export_name = "process")]
#[no_mangle]
pub extern fn process() -> u64 {
    let msg = Message{};
    let val = msg.content();
    let res = String::from_utf8(val).unwrap().to_lowercase().as_bytes().to_vec();

    msg.set_content(res);

    return 0
}
