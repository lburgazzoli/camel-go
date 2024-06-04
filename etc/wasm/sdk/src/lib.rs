#![allow(dead_code)]

use std::mem;

// ************************************
//
// Functions expected by the host
//
// ************************************

#[cfg_attr(all(target_arch = "wasm32"), export_name = "alloc")]
#[no_mangle]
pub extern "C" fn alloc(size: u32) -> *mut u8 {
    let mut buf = Vec::with_capacity(size as usize);
    let ptr = buf.as_mut_ptr();

    // tell Rust not to clean this up
    mem::forget(buf);

    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "dealloc")]
#[no_mangle]
pub unsafe extern "C" fn dealloc(ptr: &mut u8, len: i32) {
    // Retakes the pointer which allows its memory to be freed.
    let _ = Vec::from_raw_parts(ptr, 0, len as usize);
}

// ************************************
//
// Exported host functions
//
// ************************************

extern "C" {
	fn message_get_id() -> u64;
	fn message_get_time() -> u64;

	fn message_set_type(ptr: *const u8, len: i32);
    fn message_get_type() -> u64;
	fn message_set_source(ptr: *const u8, len: i32);
    fn message_get_source() -> u64;
	fn message_set_subject(ptr: *const u8, len: i32);
    fn message_get_subject() -> u64;
	fn message_set_content_schema(ptr: *const u8, len: i32);
    fn message_get_content_schema() -> u64;
	fn message_set_content_type(ptr: *const u8, len: i32);
    fn message_get_content_type() -> u64;

	fn message_set_content(ptr: *const u8, len: i32);
	fn message_get_content() -> u64;

	fn message_set_header(key_ptr: *const u8, lkey_len: i32, val_ptr: *const u8, val_len: i32);
	fn message_get_header(ptr: *const u8, len: i32) -> u64;
	fn message_remove_header(ptr: *const u8, len: i32);

	fn message_set_attribute(key_ptr: *const u8, lkey_len: i32, val_ptr: *const u8, val_len: i32);
	fn message_get_attribute(ptr: *const u8, len: i32) -> u64;
	fn message_remove_attribute(ptr: *const u8, len: i32);
}

// ************************************
//
// SDK
//
// ************************************

pub struct Message {
}

impl Message {
	pub fn id(&self) -> Vec<u8> {
	    let ptr_and_len = unsafe {
            message_get_id()
        };

        let in_ptr = (ptr_and_len >> 32) as *mut u8;
        let in_len = (ptr_and_len as u32) as usize;

        return unsafe {
            Vec::from_raw_parts(in_ptr, in_len, in_len)
        };
	}

	pub fn content(&self) -> Vec<u8> {
	    let ptr_and_len = unsafe {
            message_get_content()
        };

        let in_ptr = (ptr_and_len >> 32) as *mut u8;
        let in_len = (ptr_and_len as u32) as usize;

        return unsafe {
            Vec::from_raw_parts(in_ptr, in_len, in_len)
        };
	}

    pub fn set_content(&self, v: Vec<u8>) {
        let out_len = v.len();
        let out_ptr = v.as_ptr();

        unsafe {
            message_set_content(out_ptr, out_len as i32);
        };
    }
}