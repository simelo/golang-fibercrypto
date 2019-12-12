package data

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fibercrypto/fibercryptowallet/src/core"
	"github.com/fibercrypto/fibercryptowallet/src/util/logging"
	"github.com/skycoin/skycoin/src/cipher/bip39"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"strconv"
)

const (
	// NoSecurity  No security
	NoSecurity = iota
	// ObfuscationSecurity data obfuscation security
	ObfuscationSecurity
	// PasswordSecurity password security
	PasswordSecurity
)

const (
	Hash         = "hash"
	Entropy      = "entropy"
	SecurityType = "secType"
)

var (
	// Errors
	errParseContact        = errors.New("inserted contact cannot be parse")
	errInvalidContact      = errors.New("you try to inserted a invalid contact")
	errInvalidSecType      = errors.New("invalid security type")
	errAddrsBookHasNotInit = errors.New("address book not has init")
)

// addrsBook implement AddressBook interface for boltdb database.
type addrsBook struct {
	storage core.Storage
	key     []byte
}

var logDb = logging.MustGetLogger("AddressBook Data")

// NewAddressBook create a new instance of AddessBook and open the database of the given route.
// If database is open return bolt.errDatabaseOpen.
func NewAddressBook(storage core.Storage) core.AddressBook {
	return &addrsBook{
		storage: storage,
		key:     nil,
	}
}

// Init initialize an address book. Pass secType(security type) and password if is PasswordSecurity.
func (addrsBook *addrsBook) Init(secType int, password string) error {
	logDb.Info("initialize AddressBook")
	if !addrsBook.IsOpen() {
		return bolt.ErrDatabaseNotOpen
	}

	if addrsBook.HasInit() {
		return fmt.Errorf("address book has init")
	}

	var hash, entropy []byte
	var err error

	switch secType {
	case NoSecurity, ObfuscationSecurity:
		if err := addrsBook.insertConfig(secType, hash, entropy); err != nil {
			return err
		}
		break
	case PasswordSecurity:
		addrsBook.key = []byte(password)
		if entropy, err = addrsBook.genEntropy(); err != nil {
			return err
		}

		hash, err = bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return err
		}

		if err := addrsBook.insertConfig(secType, hash, entropy); err != nil {
			return err
		}
		break
	default:
		return errInvalidSecType
	}

	return nil
}

// Authenticate authentic a user in the Address Book. ( Only SecType : PasswordSecurity )
func (addrsBook *addrsBook) Authenticate(password string) error {
	logDb.Info("authenticate AddressBook")
	if !addrsBook.IsOpen() {
		logDb.Error(bolt.ErrDatabaseNotOpen)
		return bolt.ErrDatabaseNotOpen
	}

	if !addrsBook.HasInit() {
		logDb.Error(bolt.ErrDatabaseNotOpen)
		return errAddrsBookHasNotInit
	}

	secType, err := addrsBook.GetSecType()
	if err != nil {
		logDb.Error(err)
		return err
	}

	if secType != PasswordSecurity {
		return nil
	}

	addrsBook.key = []byte(password)
	if err := addrsBook.verifyHash(); err != nil {
		logDb.Error(err)
		return err
	}

	return nil
}

// InsertContact insert a contact into the address book.
// If any of its address exist return error.
func (addrsBook *addrsBook) InsertContact(contact core.Contact) (uint64, error) {
	if !contact.IsValid() {
		return 0, errInvalidContact
	}

	contactsList, err := addrsBook.ListContact()
	if err != nil && err != errBucketEmpty {
		return 0, err
	}
	for _, v := range contact.GetAddresses() {
		if err := addrsBook.addressExists(v, contactsList); err != nil {
			return 0, err
		}
	}
	if err := addrsBook.nameExists(contact, contactsList); err != nil {
		return 0, err
	}

	encryptedData, err := addrsBook.encryptContact(contact.(*Contact))
	if err != nil {
		return 0, err
	}

	// Commit transaction before exit.
	return addrsBook.GetStorage().InsertValue(encryptedData)
}

// GetContact get a contact by ID.
func (addrsBook *addrsBook) GetContact(id uint64) (core.Contact, error) {
	encryptData, err := addrsBook.GetStorage().GetValue(id)
	if err != nil {
		logDb.Error(err)
		return nil, err
	}
	if _, ok := encryptData.([]byte); !ok {
		logDb.Error(errValueNoMatch(encryptData.([]byte), []byte{}))
		return nil, errValueNoMatch(encryptData.([]byte), []byte{})
	}
	contact, err := addrsBook.decryptContact(encryptData.([]byte))
	if err != nil {
		return nil, err
	}
	contact.SetID(id)
	return contact, nil
}

// ListContact list all contact in the address book.
func (addrsBook *addrsBook) ListContact() ([]core.Contact, error) {
	var contactsList []core.Contact
	encryptContactList, err := addrsBook.GetStorage().ListValues()
	if err != nil {
		logDb.Error(err)
		return nil, err
	}
	for id, encryptContact := range encryptContactList {
		if _, ok := encryptContact.([]byte); !ok {
			return nil, errValueNoMatch(encryptContact, []byte{})
		}
		contact, err := addrsBook.decryptContact(encryptContact.([]byte))
		if err != nil {
			logDb.Error(err)
			return nil, err
		}
		contact.SetID(id)
		contactsList = append(contactsList, contact)
	}
	return contactsList, nil
}

// DeleteContact delete a contact from the address book by its ID.
func (addrsBook *addrsBook) DeleteContact(id uint64) error {
	logDb.Info("Removing a contact from AddressBook")
	return addrsBook.GetStorage().DeleteValue(id)
}

// UpdateContact update a contact in the address book by its ID.
func (addrsBook *addrsBook) UpdateContact(id uint64, newContact core.Contact) error {
	logDb.Infof("Updating contact with id:%d", id)
	if !newContact.IsValid() {
		return errInvalidContact
	}

	var contactsList []core.Contact
	var err error
	if contactsList, err = addrsBook.ListContact(); err != nil {
		return err
	}
	for e := range contactsList {
		if contactsList[e].GetID() == id {
			contactsList[e] = nil
			break
		}
	}

	for _, ncAddrs := range newContact.GetAddresses() {
		if err := addrsBook.addressExists(ncAddrs, contactsList); err != nil {
			return err
		}
	}
	if err := addrsBook.nameExists(newContact, contactsList); err != nil {
		return err
	}

	if _, ok := newContact.(*Contact); !ok {
		return errParseContact
	}

	encryptedData, err := addrsBook.encryptContact(newContact.(*Contact))
	if err != nil {
		return err
	}
	return addrsBook.GetStorage().UpdateValue(id, encryptedData)
}

// GetPath return database path
func (addrsBook *addrsBook) GetPath() string {
	return addrsBook.GetStorage().Path()
}

// Close shuts down the database.
func (addrsBook *addrsBook) Close() error {
	if err := addrsBook.GetStorage().Close(); err != nil {
		return err
	}
	return nil
}

// HasInit verify if database has been initialize
func (addrsBook *addrsBook) HasInit() bool {
	if addrsBook.storage.GetConfig() != nil {
		return true
	}
	return false
}

// IsOpen verify if database is open
func (addrsBook *addrsBook) IsOpen() bool {
	if addrsBook.storage.Path() != "" {
		return true
	}
	return false
}

func (addrsBook *addrsBook) GetStorage() core.Storage {
	return addrsBook.storage
}

func (addrsBook *addrsBook) GetSecType() (int, error) {
	logDb.Info("Getting security type.")
	return strconv.Atoi(addrsBook.GetStorage().GetConfig()[SecurityType])
}

// genEntropy generate annil Entropy by a mnemonic. If mnemonic is nil,
// it generate a random.
func (addrsBook *addrsBook) genEntropy() ([]byte, error) {
	mn, err := bip39.NewDefaultMnemonic()
	if err != nil {
		return nil, err
	}
	e, err := bip39.EntropyFromMnemonic(mn)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (addrsBook *addrsBook) getEntropyFromConfig() []byte {
	logDb.Info("Getting entropy.")
	return []byte(addrsBook.GetStorage().GetConfig()[Entropy])

}

// getHashFromConfig get hash from config bucket.
func (addrsBook *addrsBook) getHashFromConfig() []byte {
	logDb.Info("Getting hash.")
	return []byte(addrsBook.GetStorage().GetConfig()[Hash])
}

// encryptContact encrypt a contact using a password with AES-GCM.
func (addrsBook *addrsBook) encryptContact(c *Contact) ([]byte, error) {
	secType, err := addrsBook.GetSecType()
	if err != nil {
		return nil, err
	}
	switch secType {
	case NoSecurity:
		return c.MarshalBinary()
	case ObfuscationSecurity:
		data, err := c.MarshalBinary()
		if err != nil {
			return nil, err
		}
		return []byte(base64.StdEncoding.EncodeToString(data)), nil
	case PasswordSecurity:
		return addrsBook.encryptAESGCM(c)
	}

	return nil, fmt.Errorf("invalid security type")
}

// Decrypt a cipher message using a password with AES-GCM and return a Contact.
func (addrsBook *addrsBook) decryptContact(cipherMsg []byte) (core.Contact, error) {
	secType, err := addrsBook.GetSecType()
	if err != nil {
		return nil, err
	}
	switch secType {
	case NoSecurity:
		c := Contact{}
		if err := c.UnmarshalBinary(cipherMsg); err != nil {
			return nil, err
		}

		return &c, nil

	case ObfuscationSecurity:
		c := Contact{}
		data, err := base64.StdEncoding.DecodeString(string(cipherMsg))
		if err != nil {
			return nil, err
		}
		if err := c.UnmarshalBinary(data); err != nil {
			return nil, err
		}
		return &c, nil
	case PasswordSecurity:
		return addrsBook.decryptAESGCM(cipherMsg)
	}

	return nil, errInvalidSecType
}

//
func (addrsBook *addrsBook) verifyHash() error {
	hash := addrsBook.getHashFromConfig()
	return bcrypt.CompareHashAndPassword(hash, addrsBook.key)
}

// addressExists search an address in the list of contacts into the AddressBook.
// If find the address return error, else return nil.
func (addrsBook *addrsBook) addressExists(address core.StringAddress, contacts []core.Contact) error {
	for _, v := range contacts {
		c, ok := v.(*Contact)
		if ok {
			for _, addrs := range c.Address {
				if bytes.Compare(addrs.GetValue(), address.GetValue()) == 0 &&
					bytes.Compare(addrs.GetCoinType(), address.GetCoinType()) == 0 {
					return fmt.Errorf("Address with value: %s  and Cointype: %s alredy exist",
						address.GetValue(), address.GetCoinType())
				}
			}
		}
	}

	return nil
}

// nameExists search an name in the list of contacts into the AddressBook.
// If find the address return error, else return nil.
func (addrsBook *addrsBook) nameExists(contact core.Contact, contacts []core.Contact) error {
	for _, c := range contacts {
		if dataContact, ok := c.(*Contact); ok {
			if bytes.Compare(contact.(*Contact).Name, dataContact.Name) == 0 {
				return fmt.Errorf(" Contact with name: %s alredy exist", contact.(*Contact).Name)
			}
		}
	}

	return nil
}

func (addrsBook *addrsBook) insertConfig(secType int, hash, entropy []byte) error {
	if err := addrsBook.GetStorage().InsertConfig(
		map[string]string{
			SecurityType: strconv.Itoa(secType),
			Hash:         string(hash),
			Entropy:      string(entropy)}); err != nil {
		logDb.Error(err)
		return err
	}
	return nil
}

func (addrsBook *addrsBook) encryptAESGCM(c *Contact) ([]byte, error) {
	block, err := aes.NewCipher(
		pbkdf2.Key(addrsBook.getEntropyFromConfig(), addrsBook.key, 4096, 32, sha512.New))
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	bc, err := c.MarshalBinary()
	if err != nil {
		return nil, err
	}

	cipherText := aesGCM.Seal(nonce, nonce, bc, nil)
	return cipherText, nil
}

func (addrsBook *addrsBook) decryptAESGCM(cipherMsg []byte) (core.Contact, error) {
	block, err := aes.NewCipher(pbkdf2.Key(
		addrsBook.getEntropyFromConfig(), addrsBook.key, 4096, 32, sha512.New))
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	var c Contact
	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := cipherMsg[:nonceSize], cipherMsg[nonceSize:]

	data, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	if err := c.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return &c, nil
}
