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
	return sha1hmac.Sum(data)
}

func HMac_MD5Bytes(data, key []byte) []byte {
	vhmac := hmac.New(md5.New, key)
	return vhmac.Sum(data)
}

func HMac_SHA256Bytes(data, key []byte) []byte {
	vhmac := hmac.New(sha256.New, key)
	return vhmac.Sum(data)
}

func HMac_SHA512Bytes(data, key []byte) []byte {
	vhmac := hmac.New(sha512.New, key)
	return vhmac.Sum(data)
}

func SHA1Bytes(data []byte) []byte {
	sha1hash := sha1.New()
	return sha1hash.Sum(data)
}

func MD5Bytes(data []byte) []byte {
	md5hash := md5.New()
	return md5hash.Sum(data)
}

func SHA256Bytes(data []byte) []byte {
	sha256Hash := sha256.New()
	return sha256Hash.Sum(data)
}

func SHA512Bytes(data []byte) []byte {
	sha512Hash := sha512.New()
	return sha512Hash.Sum(data)
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
