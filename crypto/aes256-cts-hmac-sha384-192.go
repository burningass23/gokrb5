package crypto

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/sha512"
	"hash"

	"gopkg.in/burningass23/gokrb5.v7/crypto/common"
	"gopkg.in/burningass23/gokrb5.v7/crypto/rfc8009"
	"gopkg.in/burningass23/gokrb5.v7/iana/chksumtype"
	"gopkg.in/burningass23/gokrb5.v7/iana/etypeID"
)

// RFC https://tools.ietf.org/html/rfc8009

// Aes256CtsHmacSha384192 implements Kerberos encryption type aes256-cts-hmac-sha384-192
type Aes256CtsHmacSha384192 struct {
}

// GetETypeID returns the EType ID number.
func (e Aes256CtsHmacSha384192) GetETypeID() int32 {
	return etypeID.AES256_CTS_HMAC_SHA384_192
}

// GetHashID returns the checksum type ID number.
func (e Aes256CtsHmacSha384192) GetHashID() int32 {
	return chksumtype.HMAC_SHA384_192_AES256
}

// GetKeyByteSize returns the number of bytes for key of this etype.
func (e Aes256CtsHmacSha384192) GetKeyByteSize() int {
	return 192 / 8
}

// GetKeySeedBitLength returns the number of bits for the seed for key generation.
func (e Aes256CtsHmacSha384192) GetKeySeedBitLength() int {
	return e.GetKeyByteSize() * 8
}

// GetHashFunc returns the hash function for this etype.
func (e Aes256CtsHmacSha384192) GetHashFunc() func() hash.Hash {
	return sha512.New384
}

// GetMessageBlockByteSize returns the block size for the etype's messages.
func (e Aes256CtsHmacSha384192) GetMessageBlockByteSize() int {
	return 1
}

// GetDefaultStringToKeyParams returns the default key derivation parameters in string form.
func (e Aes256CtsHmacSha384192) GetDefaultStringToKeyParams() string {
	return "00008000"
}

// GetConfounderByteSize returns the byte count for confounder to be used during cryptographic operations.
func (e Aes256CtsHmacSha384192) GetConfounderByteSize() int {
	return aes.BlockSize
}

// GetHMACBitLength returns the bit count size of the integrity hash.
func (e Aes256CtsHmacSha384192) GetHMACBitLength() int {
	return 192
}

// GetCypherBlockBitLength returns the bit count size of the cypher block.
func (e Aes256CtsHmacSha384192) GetCypherBlockBitLength() int {
	return aes.BlockSize * 8
}

// StringToKey returns a key derived from the string provided.
func (e Aes256CtsHmacSha384192) StringToKey(secret string, salt string, s2kparams string) ([]byte, error) {
	saltp := rfc8009.GetSaltP(salt, "aes256-cts-hmac-sha384-192")
	return rfc8009.StringToKey(secret, saltp, s2kparams, e)
}

// RandomToKey returns a key from the bytes provided.
func (e Aes256CtsHmacSha384192) RandomToKey(b []byte) []byte {
	return rfc8009.RandomToKey(b)
}

// EncryptData encrypts the data provided.
func (e Aes256CtsHmacSha384192) EncryptData(key, data []byte) ([]byte, []byte, error) {
	return rfc8009.EncryptData(key, data, e)
}

// EncryptMessage encrypts the message provided and concatenates it with the integrity hash to create an encrypted message.
func (e Aes256CtsHmacSha384192) EncryptMessage(key, message []byte, usage uint32) ([]byte, []byte, error) {
	return rfc8009.EncryptMessage(key, message, usage, e)
}

// DecryptData decrypts the data provided.
func (e Aes256CtsHmacSha384192) DecryptData(key, data []byte) ([]byte, error) {
	return rfc8009.DecryptData(key, data, e)
}

// DecryptMessage decrypts the message provided and verifies the integrity of the message.
func (e Aes256CtsHmacSha384192) DecryptMessage(key, ciphertext []byte, usage uint32) ([]byte, error) {
	return rfc8009.DecryptMessage(key, ciphertext, usage, e)
}

// DeriveKey derives a key from the protocol key based on the usage value.
func (e Aes256CtsHmacSha384192) DeriveKey(protocolKey, usage []byte) ([]byte, error) {
	return rfc8009.DeriveKey(protocolKey, usage, e), nil
}

// DeriveRandom generates data needed for key generation.
func (e Aes256CtsHmacSha384192) DeriveRandom(protocolKey, usage []byte) ([]byte, error) {
	return rfc8009.DeriveRandom(protocolKey, usage, e)
}

// VerifyIntegrity checks the integrity of the ciphertext message.
// As the hash is calculated over the iv concatenated with the AES cipher output not the plaintext the pt value to this
// interface method is not use. Pass any []byte.
func (e Aes256CtsHmacSha384192) VerifyIntegrity(protocolKey, ct, pt []byte, usage uint32) bool {
	// We don't need ib just there for the interface
	return rfc8009.VerifyIntegrity(protocolKey, ct, usage, e)
}

// GetChecksumHash returns a keyed checksum hash of the bytes provided.
func (e Aes256CtsHmacSha384192) GetChecksumHash(protocolKey, data []byte, usage uint32) ([]byte, error) {
	return common.GetHash(data, protocolKey, common.GetUsageKc(usage), e)
}

// VerifyChecksum compares the checksum of the message bytes is the same as the checksum provided.
func (e Aes256CtsHmacSha384192) VerifyChecksum(protocolKey, data, chksum []byte, usage uint32) bool {
	c, err := e.GetChecksumHash(protocolKey, data, usage)
	if err != nil {
		return false
	}
	return hmac.Equal(chksum, c)
}
