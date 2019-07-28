package walletsManager

import (
	"github.com/fibercrypto/FiberCryptoWallet/src/coin/skycoin/blockchain/api"
	"github.com/fibercrypto/FiberCryptoWallet/src/core"
	"github.com/fibercrypto/FiberCryptoWallet/src/models/wallets"
	qtcore "github.com/therecipe/qt/core"
)

func init() {
	WalletManager_QmlRegisterType2("WalletsManager", 1, 0, "WalletManager")
}

type WalletManager struct {
	qtcore.QObject
	WalletEnv     core.WalletEnv
	SeedGenerator core.SeedGenerator

	_ func()                                                                       `constructor:"init"`
	_ func(seed string, label string, password string, scanN int) *wallets.QWallet `slot:"createEncryptedWallet"`
	_ func(seed string, label string, scanN int) *wallets.QWallet                  `slot:"createUnencryptedWallet"`
	_ func(entropy int) string                                                     `slot:"getNewSeed"`
	_ func(seed string) int                                                        `slot:"verifySeed"`
	_ func(id string, n int, password string)                                      `slot:"newWalletAddress"`
	_ func(id string, password string)                                             `slot:"encryptWallet"`
	_ func(id string, password string)                                             `slot:"decryptWallet"`
	_ func() []*wallets.QWallet                                                    `slot:"getWallets"`
	_ func(id string) []*wallets.QAddress                                          `slot:"getAddresses"`
}

func (walletM *WalletManager) init() {
	walletM.ConnectCreateEncryptedWallet(walletM.createEncryptedWallet)
	walletM.ConnectCreateUnencryptedWallet(walletM.createUnencryptedWallet)
	walletM.ConnectGetNewSeed(walletM.getNewSeed)
	walletM.ConnectVerifySeed(walletM.verifySeed)
	walletM.ConnectNewWalletAddress(walletM.newWalletAddress)
	walletM.ConnectEncryptWallet(walletM.encryptWallet)
	walletM.ConnectDecryptWallet(walletM.decryptWallet)
	walletM.ConnectGetWallets(walletM.getWallets)
	walletM.ConnectGetAddresses(walletM.getAddresses)

	walletM.WalletEnv = new(api.WalletNode)
	walletM.SeedGenerator = new(api.SeedService)

}

func (walletM *WalletManager) createEncryptedWallet(seed, label, password string, scanN int) *wallets.QWallet {
	pwd := func(message string) (string, error) {
		return password, nil
	}
	wlt, err := walletM.WalletEnv.GetWalletSet().CreateWallet(label, seed, true, pwd, scanN)
	if err != nil {
		return nil
	}

	return fromWalletToQWallet(wlt, true)

}

func (walletM *WalletManager) createUnencryptedWallet(seed, label string, scanN int) *wallets.QWallet {
	pwd := func(message string) (string, error) {
		return "", nil
	}

	wlt, err := walletM.WalletEnv.GetWalletSet().CreateWallet(label, seed, false, pwd, scanN)
	if err != nil {
		return nil
	}
	return fromWalletToQWallet(wlt, false)

}

func (walletM *WalletManager) getNewSeed(entropy int) string {
	seed, err := walletM.SeedGenerator.GenerateMnemonic(entropy)
	if err != nil {
		return ""
	}
	return seed
}

func (walletM *WalletManager) verifySeed(seed string) int {
	ok, err := walletM.SeedGenerator.VerifyMnemonic(seed)
	if err != nil {
		return 0
	}
	if ok {
		return 1
	}
	return 0

}

func (walletM *WalletManager) encryptWallet(id, password string) {
	pwd := func(message string) (string, error) {
		return password, nil
	}
	walletM.WalletEnv.GetStorage().Encrypt(id, pwd)
}

func (walletM *WalletManager) decryptWallet(id, password string) {
	pwd := func(message string) (string, error) {
		return password, nil
	}
	walletM.WalletEnv.GetStorage().Decrypt(id, pwd)
}

func (walletM *WalletManager) newWalletAddress(id string, n int, password string) {
	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(id)
	pwd := func(message string) (string, error) {
		return password, nil
	}
	wltEntrieslen := 0
	it, err := wlt.GetLoadedAddresses()
	if err != nil {
		return
	}
	for it.Next() {
		wltEntrieslen++
	}
	wlt.GenAddresses(core.AccountAddress, uint32(wltEntrieslen), uint32(n), pwd)
}

func (walletM *WalletManager) getWallets() []*wallets.QWallet {
	qwallets := make([]*wallets.QWallet, 0)
	it := walletM.WalletEnv.GetWalletSet().ListWallets()
	for it.Next() {
		encrypted, err := walletM.WalletEnv.GetStorage().IsEncrypted(it.Value().GetId())
		if err != nil {
			continue
		}
		if encrypted {
			qwallets = append(qwallets, fromWalletToQWallet(it.Value(), true))
		} else {
			qwallets = append(qwallets, fromWalletToQWallet(it.Value(), false))
		}

	}
	return qwallets
}

func (walletM *WalletManager) getAddresses(Id string) []*wallets.QAddress {
	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(Id)
	qaddresses := make([]*wallets.QAddress, 0)
	it, err := wlt.GetLoadedAddresses()
	if err != nil {
		return nil
	}
	for it.Next() {
		addr := it.Value()
		qaddress := wallets.NewQAddress(nil)
		qaddress.SetAddress(addr.String())
		sky, err := addr.GetCryptoAccount().GetBalance("Sky")
		if err != nil {
			continue
		}
		qaddress.SetAddressSky(sky)
		coinH, err := addr.GetCryptoAccount().GetBalance("CoinHour")
		if err != nil {
			continue
		}
		qaddress.SetAddressCoinHours(coinH)
		qaddresses = append(qaddresses, qaddress)

	}
	return qaddresses
}

func fromWalletToQWallet(wlt core.Wallet, isEncrypted bool) *wallets.QWallet {
	qwallet := wallets.NewQWallet(nil)
	qwallet.SetName(wlt.GetLabel())
	qwallet.SetFileName(wlt.GetId())
	qwallet.SetEncryptionEnabled(0)
	if isEncrypted {
		qwallet.SetEncryptionEnabled(1)
	}
	bl, err := wlt.GetCryptoAccount().GetBalance("Sky")
	if err != nil {
		bl = 0
	}
	qwallet.SetSky(bl)
	bl, err = wlt.GetCryptoAccount().GetBalance("CoinHour")
	if err != nil {
		bl = 0
	}
	qwallet.SetCoinHours(bl)

	return qwallet
}
