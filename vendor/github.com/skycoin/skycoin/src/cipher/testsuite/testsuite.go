/*
Package testsuite is the cipher testdata testsuite
*/
package testsuite

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/base58"
	"github.com/skycoin/skycoin/src/cipher/bip32"
	secp256k1 "github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

// InputTestDataJSON contains hashes to be signed
type InputTestDataJSON struct {
	Hashes []string `json:"hashes"`
}

// KeysTestDataJSON contains address, public key, secret key and list of signatures
type KeysTestDataJSON struct {
	Address        string   `json:"address"`
	BitcoinAddress string   `json:"bitcoin_address"`
	Secret         string   `json:"secret"`
	Public         string   `json:"public"`
	Signatures     []string `json:"signatures,omitempty"`
}

// SeedTestDataJSON contains data generated by Seed
type SeedTestDataJSON struct {
	Seed string             `json:"seed"`
	Keys []KeysTestDataJSON `json:"keys"`
}

// Bip32KeysTestDataJSON contains address, public key, secret key and list of signatures
type Bip32KeysTestDataJSON struct {
	Path        string `json:"path"`
	XPriv       string `json:"xpriv"`
	XPub        string `json:"xpub"`
	Identifier  string `json:"identifier"`
	Depth       byte   `json:"depth"`
	ChildNumber uint32 `json:"child_number"`

	KeysTestDataJSON
}

// Bip32SeedTestDataJSON contains data generated by Seed
type Bip32SeedTestDataJSON struct {
	Seed         string                  `json:"seed"`
	BasePath     string                  `json:"base_path"`
	ChildNumbers []uint32                `json:"child_numbers"`
	Keys         []Bip32KeysTestDataJSON `json:"keys"`
}

// InputTestData contains hashes to be signed
type InputTestData struct {
	Hashes []cipher.SHA256
}

// ToJSON converts InputTestData to InputTestDataJSON
func (d *InputTestData) ToJSON() *InputTestDataJSON {
	hashes := make([]string, len(d.Hashes))
	for i, h := range d.Hashes {
		hashes[i] = h.Hex()
	}

	return &InputTestDataJSON{
		Hashes: hashes,
	}
}

// InputTestDataFromJSON converts InputTestDataJSON to InputTestData
func InputTestDataFromJSON(d *InputTestDataJSON) (*InputTestData, error) {
	hashes := make([]cipher.SHA256, len(d.Hashes))
	for i, h := range d.Hashes {
		var err error
		hashes[i], err = cipher.SHA256FromHex(h)
		if err != nil {
			return nil, err
		}
	}

	return &InputTestData{
		Hashes: hashes,
	}, nil
}

// KeysTestData contains address, public key, secret key and list of signatures
type KeysTestData struct {
	Address        cipher.Address
	BitcoinAddress cipher.BitcoinAddress
	Secret         cipher.SecKey
	Public         cipher.PubKey
	Signatures     []cipher.Sig
}

// ToJSON converts KeysTestData to KeysTestDataJSON
func (k *KeysTestData) ToJSON() *KeysTestDataJSON {
	sigs := make([]string, len(k.Signatures))
	for i, s := range k.Signatures {
		sigs[i] = s.Hex()
	}

	return &KeysTestDataJSON{
		Address:        k.Address.String(),
		BitcoinAddress: k.BitcoinAddress.String(),
		Secret:         k.Secret.Hex(),
		Public:         k.Public.Hex(),
		Signatures:     sigs,
	}
}

// KeysTestDataFromJSON converts KeysTestDataJSON to KeysTestData
func KeysTestDataFromJSON(d *KeysTestDataJSON) (*KeysTestData, error) {
	addr, err := cipher.DecodeBase58Address(d.Address)
	if err != nil {
		return nil, err
	}

	btcAddr, err := cipher.DecodeBase58BitcoinAddress(d.BitcoinAddress)
	if err != nil {
		return nil, err
	}

	s, err := cipher.SecKeyFromHex(d.Secret)
	if err != nil {
		return nil, err
	}

	p, err := cipher.PubKeyFromHex(d.Public)
	if err != nil {
		return nil, err
	}

	var sigs []cipher.Sig
	if d.Signatures != nil {
		sigs = make([]cipher.Sig, len(d.Signatures))
		for i, s := range d.Signatures {
			var err error
			sigs[i], err = cipher.SigFromHex(s)
			if err != nil {
				return nil, err
			}
		}
	}

	return &KeysTestData{
		Address:        addr,
		BitcoinAddress: btcAddr,
		Secret:         s,
		Public:         p,
		Signatures:     sigs,
	}, nil
}

// SeedTestData contains data generated by Seed
type SeedTestData struct {
	Seed []byte
	Keys []KeysTestData
}

// ToJSON converts SeedTestData to SeedTestDataJSON
func (s *SeedTestData) ToJSON() *SeedTestDataJSON {
	keys := make([]KeysTestDataJSON, len(s.Keys))
	for i, k := range s.Keys {
		kj := k.ToJSON()
		keys[i] = *kj
	}

	return &SeedTestDataJSON{
		Seed: base64.StdEncoding.EncodeToString(s.Seed),
		Keys: keys,
	}
}

// SeedTestDataFromJSON converts SeedTestDataJSON to SeedTestData
func SeedTestDataFromJSON(d *SeedTestDataJSON) (*SeedTestData, error) {
	seed, err := base64.StdEncoding.DecodeString(d.Seed)
	if err != nil {
		return nil, err
	}

	keys := make([]KeysTestData, len(d.Keys))
	for i, kj := range d.Keys {
		k, err := KeysTestDataFromJSON(&kj)
		if err != nil {
			return nil, err
		}
		keys[i] = *k
	}

	return &SeedTestData{
		Seed: seed,
		Keys: keys,
	}, nil
}

// ValidateSeedData validates the provided SeedTestData against the current cipher library.
// inputData is required if SeedTestData contains signatures
func ValidateSeedData(seedData *SeedTestData, inputData *InputTestData) error {
	keys := cipher.MustGenerateDeterministicKeyPairs(seedData.Seed, len(seedData.Keys))
	if len(seedData.Keys) != len(keys) {
		return errors.New("cipher.GenerateDeterministicKeyPairs generated an unexpected number of keys")
	}

	for i, s := range keys {
		if err := validateKeyTestData(inputData, s, seedData.Keys[i]); err != nil {
			return err
		}
	}

	return nil
}

func validateKeyTestData(inputData *InputTestData, s cipher.SecKey, data KeysTestData) error {
	if s == (cipher.SecKey{}) {
		return errors.New("secret key is null")
	}
	if data.Secret != s {
		return errors.New("generated secret key does not match provided secret key")
	}

	p := cipher.MustPubKeyFromSecKey(s)
	if p == (cipher.PubKey{}) {
		return errors.New("public key is null")
	}
	if data.Public != p {
		return errors.New("derived public key does not match provided public key")
	}

	addr1 := cipher.AddressFromPubKey(p)
	if addr1 == (cipher.Address{}) {
		return errors.New("address is null")
	}
	if data.Address != addr1 {
		return errors.New("derived address does not match provided address")
	}

	addr2 := cipher.MustAddressFromSecKey(s)
	if addr1 != addr2 {
		return errors.New("cipher.AddressFromPubKey and cipher.AddressFromSecKey generated different addresses")
	}

	btcAddr1 := cipher.BitcoinAddressFromPubKey(p)
	if btcAddr1 == (cipher.BitcoinAddress{}) {
		return errors.New("bitcoin address is null")
	}
	if data.BitcoinAddress != btcAddr1 {
		return errors.New("derived bitcoin address does not match provided bitcoin address")
	}

	btcAddr2 := cipher.MustBitcoinAddressFromSecKey(s)
	if btcAddr1 != btcAddr2 {
		return errors.New("cipher.BitcoinAddressFromPubKey and cipher.BitcoinAddressFromSecKey generated different addresses")
	}

	validSec := secp256k1.VerifySeckey(s[:])
	if validSec != 1 {
		return errors.New("secp256k1.VerifySeckey failed")
	}

	validPub := secp256k1.VerifyPubkey(p[:])
	if validPub != 1 {
		return errors.New("secp256k1.VerifyPubkey failed")
	}

	if inputData == nil && len(data.Signatures) != 0 {
		return errors.New("seed data contains signatures but input data was not provided")
	}

	if inputData != nil {
		if len(data.Signatures) != len(inputData.Hashes) {
			return errors.New("Number of signatures in seed data does not match number of hashes in input data")
		}

		for j, h := range inputData.Hashes {
			sig := data.Signatures[j]
			if sig == (cipher.Sig{}) {
				return errors.New("provided signature is null")
			}

			err := cipher.VerifyPubKeySignedHash(p, sig, h)
			if err != nil {
				return fmt.Errorf("cipher.VerifyPubKeySignedHash failed: %v", err)
			}

			err = cipher.VerifyAddressSignedHash(addr1, sig, h)
			if err != nil {
				return fmt.Errorf("cipher.VerifyAddressSignedHash failed: %v", err)
			}

			err = cipher.VerifySignatureRecoverPubKey(sig, h)
			if err != nil {
				return fmt.Errorf("cipher.VerifySignatureRecoverPubKey failed: %v", err)
			}

			p2, err := cipher.PubKeyFromSig(sig, h)
			if err != nil {
				return fmt.Errorf("cipher.PubKeyFromSig failed: %v", err)
			}

			if p != p2 {
				return errors.New("public key derived from signature does not match public key derived from secret")
			}

			sig2 := cipher.MustSignHash(h, s)
			if sig2 == (cipher.Sig{}) {
				return errors.New("created signature is null")
			}

			// NOTE: signatures are not deterministic, they use a nonce,
			// so we don't compare the generated sig to the provided sig
		}
	}

	return nil
}

// Bip32KeysTestData contains address, public key, secret key and list of signatures
type Bip32KeysTestData struct {
	Path  string
	XPriv *bip32.PrivateKey
	KeysTestData
}

// ToJSON converts Bip32KeysTestData to Bip32KeysTestDataJSON
func (k *Bip32KeysTestData) ToJSON() *Bip32KeysTestDataJSON {
	sigs := make([]string, len(k.Signatures))
	for i, s := range k.Signatures {
		sigs[i] = s.Hex()
	}

	return &Bip32KeysTestDataJSON{
		Path:        k.Path,
		XPriv:       k.XPriv.String(),
		XPub:        k.XPriv.PublicKey().String(),
		Identifier:  hex.EncodeToString(k.XPriv.Identifier()),
		Depth:       k.XPriv.Depth,
		ChildNumber: k.XPriv.ChildNumber(),

		KeysTestDataJSON: KeysTestDataJSON{
			Address:        k.Address.String(),
			BitcoinAddress: k.BitcoinAddress.String(),
			Secret:         k.Secret.Hex(),
			Public:         k.Public.Hex(),
			Signatures:     sigs,
		},
	}
}

// Bip32KeysTestDataFromJSON converts Bip32KeysTestDataJSON to Bip32KeysTestData
func Bip32KeysTestDataFromJSON(d *Bip32KeysTestDataJSON) (*Bip32KeysTestData, error) {
	addr, err := cipher.DecodeBase58Address(d.Address)
	if err != nil {
		return nil, err
	}

	btcAddr, err := cipher.DecodeBase58BitcoinAddress(d.BitcoinAddress)
	if err != nil {
		return nil, err
	}

	s, err := cipher.SecKeyFromHex(d.Secret)
	if err != nil {
		return nil, err
	}

	p, err := cipher.PubKeyFromHex(d.Public)
	if err != nil {
		return nil, err
	}

	var sigs []cipher.Sig
	if d.Signatures != nil {
		sigs = make([]cipher.Sig, len(d.Signatures))
		for i, s := range d.Signatures {
			var err error
			sigs[i], err = cipher.SigFromHex(s)
			if err != nil {
				return nil, err
			}
		}
	}

	xPrivBytes, err := base58.Decode(d.XPriv)
	if err != nil {
		return nil, err
	}

	xPriv, err := bip32.DeserializePrivateKey(xPrivBytes)
	if err != nil {
		return nil, err
	}

	identifier, err := hex.DecodeString(d.Identifier)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(xPriv.Identifier(), identifier) {
		return nil, errors.New("xpriv identifier does not match identifier")
	}

	if xPriv.Depth != d.Depth {
		return nil, errors.New("xpriv depth does not match depth")
	}

	if xPriv.PublicKey().String() != d.XPub {
		return nil, errors.New("xpub derived from xpriv does not match xpub")
	}

	if xPriv.ChildNumber() != d.ChildNumber {
		return nil, errors.New("xpriv child number does not match child number")
	}

	return &Bip32KeysTestData{
		Path:  d.Path,
		XPriv: xPriv,

		KeysTestData: KeysTestData{
			Address:        addr,
			BitcoinAddress: btcAddr,
			Secret:         s,
			Public:         p,
			Signatures:     sigs,
		},
	}, nil
}

// Bip32SeedTestData contains data generated by Seed
type Bip32SeedTestData struct {
	Seed         []byte
	BasePath     string
	ChildNumbers []uint32
	Keys         []Bip32KeysTestData
}

// ToJSON converts Bip32SeedTestData to Bip32SeedTestDataJSON
func (s *Bip32SeedTestData) ToJSON() *Bip32SeedTestDataJSON {
	keys := make([]Bip32KeysTestDataJSON, len(s.Keys))
	for i, k := range s.Keys {
		kj := k.ToJSON()
		keys[i] = *kj
	}

	return &Bip32SeedTestDataJSON{
		Seed:         base64.StdEncoding.EncodeToString(s.Seed),
		BasePath:     s.BasePath,
		ChildNumbers: s.ChildNumbers,
		Keys:         keys,
	}
}

// Bip32SeedTestDataFromJSON converts Bip32SeedTestDataJSON to Bip32SeedTestData
func Bip32SeedTestDataFromJSON(d *Bip32SeedTestDataJSON) (*Bip32SeedTestData, error) {
	seed, err := base64.StdEncoding.DecodeString(d.Seed)
	if err != nil {
		return nil, err
	}

	keys := make([]Bip32KeysTestData, len(d.Keys))
	for i, kj := range d.Keys {
		k, err := Bip32KeysTestDataFromJSON(&kj)
		if err != nil {
			return nil, err
		}
		keys[i] = *k
	}

	return &Bip32SeedTestData{
		Seed:         seed,
		BasePath:     d.BasePath,
		ChildNumbers: d.ChildNumbers,
		Keys:         keys,
	}, nil
}

// ValidateBip32SeedData validates the provided Bip32SeedTestData against the current cipher library.
// inputData is required if Bip32SeedTestData contains signatures
func ValidateBip32SeedData(seedData *Bip32SeedTestData, inputData *InputTestData) error {
	mk, err := bip32.NewPrivateKeyFromPath(seedData.Seed, seedData.BasePath)
	if err != nil {
		return err
	}

	if len(seedData.ChildNumbers) != len(seedData.Keys) {
		return errors.New("len(seedData.ChildNumbers) must equal len(seedData.Keys)")
	}

	for i, n := range seedData.ChildNumbers {
		k, err := mk.NewPrivateChildKey(n)
		if err != nil {
			return err
		}

		if err := validateBip32KeyTestData(inputData, seedData.BasePath, seedData.Seed, k, n, seedData.Keys[i]); err != nil {
			return err
		}
	}

	return nil
}

func validateBip32KeyTestData(inputData *InputTestData, basePath string, seed []byte, s *bip32.PrivateKey, childNumber uint32, data Bip32KeysTestData) error {
	path := fmt.Sprintf("%s/%d", basePath, childNumber)
	pathXPriv, err := bip32.NewPrivateKeyFromPath(seed, path)
	if err != nil {
		return err
	}

	if s.String() != pathXPriv.String() {
		return errors.New("xpriv generated with NewPrivateChildKey differs from xpriv generated with NewPrivateKeyFromPath")
	}

	pubKey := cipher.MustNewPubKey(s.PublicKey().Key)
	secKey := cipher.MustNewSecKey(s.Key)

	if cipher.MustPubKeyFromSecKey(secKey) != pubKey {
		return errors.New("pubkey derived from bip32 key does not match pubkey derived from bip32 key converted to cipher key")
	}

	return validateKeyTestData(inputData, secKey, data.KeysTestData)
}
