package main

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"strconv"
	"unsafe"
)

var defaultLetters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString returns a random string with a fixed length
func randomString(n int) []byte {

	letters := defaultLetters

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return b
}

// DS2b returns an 8-byte big endian representation of Digit string
// v ("123456") -> uint64(123456) -> 8-byte big endian.
func DS2b(v string) []byte {
	i, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return []byte("")
	}
	return I2b(i)
}

// DS2i returns uint64 of Digit string
// v ("123456") -> uint64(123456).
func DS2i(v string) uint64 {
	i, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return uint64(0)
	}
	return i
}

// I2b returns an 8-byte big endian representation of v
// v uint64(123456) -> 8-byte big endian.
func I2b(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// B2i return an int64 of v
// v (8-byte big endian) -> uint64(123456).
func B2i(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

// B2ds return a Digit string of v
// v (8-byte big endian) -> uint64(123456) -> "123456".
func B2ds(v []byte) string {
	return strconv.FormatUint(binary.BigEndian.Uint64(v), 10)
}

// B2s converts byte slice to a string without memory allocation.
// []byte("abc") -> "abc" s
func B2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2b converts string to a byte slice without memory allocation.
// "abc" -> []byte("abc")
func S2b(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
