package Utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash/crc32"
	"hash/crc64"
)

func HMac_SHA1Bytes(data, key []byte) []byte {
	sha1hmac := hmac.New(sha1.New, key)
	sha1hmac.Write(data)
	return sha1hmac.Sum(nil)
}

func HMac_MD5Bytes(data, key []byte) []byte {
	vhmac := hmac.New(md5.New, key)
	vhmac.Write(data)
	return vhmac.Sum(nil)
}

func HMac_SHA256Bytes(data, key []byte) []byte {
	vhmac := hmac.New(sha256.New, key)
	vhmac.Write(data)
	return vhmac.Sum(nil)
}

func HMac_SHA512Bytes(data, key []byte) []byte {
	vhmac := hmac.New(sha512.New, key)
	vhmac.Write(data)
	return vhmac.Sum(nil)
}

func SHA1Bytes(data []byte) []byte {
	sha1hash := sha1.New()
	sha1hash.Write(data)
	return sha1hash.Sum(nil)
}

func MD5Bytes(data []byte) []byte {
	md5hash := md5.New()
	md5hash.Write(data)
	return md5hash.Sum(nil)
}

func SHA256Bytes(data []byte) []byte {
	sha256Hash := sha256.New()
	sha256Hash.Write(data)
	return sha256Hash.Sum(nil)
}

func SHA512Bytes(data []byte) []byte {
	sha512Hash := sha512.New()
	sha512Hash.Write(data)
	return sha512Hash.Sum(nil)
}

func Crc32IEEE(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func Crc64ISO(data []byte) uint64 {
	return crc64.Checksum(data, crc64.MakeTable(crc64.ISO))
}

func Crc64ECMA(data []byte) uint64 {
	return crc64.Checksum(data, crc64.MakeTable(crc64.ECMA))
}

func SHA1(data []byte) string {
	return hex.EncodeToString(SHA1Bytes(data))
}

func MD5(data []byte) string {
	return hex.EncodeToString(MD5Bytes(data))
}

func SHA256(data []byte) string {
	return hex.EncodeToString(SHA256Bytes(data))
}

func SHA512(data []byte) string {
	return hex.EncodeToString(SHA512Bytes(data))
}

func HMac_SHA1(data, key []byte) string {
	return hex.EncodeToString(HMac_SHA1Bytes(data, key))
}

func HMac_MD5(data, key []byte) string {
	return hex.EncodeToString(HMac_MD5Bytes(data, key))
}

func HMac_SHA256(data, key []byte) string {
	return hex.EncodeToString(HMac_SHA256Bytes(data, key))
}

func HMac_SHA512(data, key []byte) string {
	return hex.EncodeToString(HMac_SHA512Bytes(data, key))
}
