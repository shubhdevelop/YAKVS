package store

import (
	"unsafe"
)

const (
	OBJ_STRING = 0
	OBJ_LIST   = 1
	OBJ_SET    = 2
	OBJ_ZSET   = 3
	OBJ_HASH   = 4
	// more later
)

const (
	OBJ_ENCODING_RAW    = 0
	OBJ_ENCODING_INT    = 1
	OBJ_ENCODING_HT     = 2
	OBJ_ENCODING_ZIPMAP = 3
	// ... etc
)

type kvObj struct {
	/* 
	In Go, we can't use bitfields like C, so we need to pack these manually
	The first byte contains: type (4 bits) and encoding (4 bits)
	*/
	typeAndEncoding uint8
	/*
	 lru is 24 bits (3 bytes) in C, we'll use uint32 and manually handle it
	 only lower 24 bits used
	*/
	lru uint32 
	refcount int32
	ptr      unsafe.Pointer
}

// Helper methods to get/set the bitfield values
func (r *kvObj) getType() uint8 {
	return r.typeAndEncoding >> 4
}

func (r *kvObj) setType(t uint8) {
	r.typeAndEncoding = (r.typeAndEncoding & 0x0F) | (t << 4)
}

func (r *kvObj) getEncoding() uint8 {
	return r.typeAndEncoding & 0x0F
}

func (r *kvObj) setEncoding(e uint8) {
	r.typeAndEncoding = (r.typeAndEncoding & 0xF0) | (e & 0x0F)
}

func (r *kvObj) getLRU() uint32 {
	return r.lru & 0xFFFFFF // mask to 24 bits
}

func (r *kvObj) setLRU(lru uint32) {
	r.lru = lru & 0xFFFFFF // ensure only 24 bits
}

func createIntObj(value int) *kvObj {
	obj := &kvObj{
		refcount: 1,
		lru: 0,
	}
	
	obj.setType(OBJ_STRING)
	obj.setEncoding(OBJ_ENCODING_INT)
	obj.ptr = unsafe.Pointer(&value)
	return obj
}

func createStringObj(value string) *kvObj {
	obj := &kvObj{
		refcount: 1,
		lru: 0,
	}
	
	obj.setType(OBJ_STRING)
	obj.setEncoding(OBJ_ENCODING_RAW)
	obj.ptr = unsafe.Pointer(&value)
	return obj
}
