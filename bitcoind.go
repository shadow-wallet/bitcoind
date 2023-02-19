package bitcoind

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Bitcoind struct {
	client *rpcClient
}

func New(addr, user, pswd string) (*Bitcoind, error) {
	cli := &rpcClient{
		addr:       addr,
		user:       user,
		pswd:       pswd,
		httpClient: http.DefaultClient,
	}
	return &Bitcoind{cli}, nil
}

// LoadWallet loads a wallet.
func (b *Bitcoind) LoadWallet(account string) (err error) {
	_, err = b.client.call("", "loadwallet", []any{account})
	return
}

func (b *Bitcoind) GetBalance(account string, minconf uint64) (balance float64, err error) {
	r, err := b.client.call(account, "getbalance", []any{"*", minconf})
	if err = handleError(err, &r); err != nil {
		return
	}
	balance, err = strconv.ParseFloat(string(r.Result), 64)
	return
}

// ImportPrivKey Adds a private key (as returned by dumpprivkey) to your wallet.
// This may take a while, as a rescan is done, looking for existing transactions.
// Optional [rescan] parameter added in 0.8.0.
// Note: There's no need to import public key, as in ECDSA (unlike RSA) this
// can be computed from private key.
func (b *Bitcoind) ImportPrivKey(privKey, account string, rescan bool) error {
	r, err := b.client.call(account, "importprivkey", []any{privKey, account, rescan})
	return handleError(err, &r)
}

// GetNewAddress return a new address for account [account].
func (b *Bitcoind) GetNewAddress(account string) (addr string, err error) {
	r, err := b.client.call(account, "getnewaddress", []any{account})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &addr)
	return
}

// CreateWallet creates a new wallet.
func (b *Bitcoind) CreateWallet(account string) (err error) {
	_, err = b.client.call("", "createwallet", []any{account, false, false, "", false, false})
	return
}

// GetPeerInfo returns data about each connected node
func (b *Bitcoind) GetPeerInfo() (peerInfo []Peer, err error) {
	r, err := b.client.call("", "getpeerinfo", nil)
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &peerInfo)
	return
}

// EncryptWallet encrypts the wallet with <passphrase>.
func (b *Bitcoind) EncryptWallet(account, passphrase string) error {
	r, err := b.client.call(account, "encryptwallet", []any{passphrase})
	return handleError(err, &r)
}

// SendToAddress send an amount to a given address
func (b *Bitcoind) SendToAddress(fromAccount, toAddress string, amount float64, comment, commentTo string, subfee bool) (txID string, err error) {
	r, err := b.client.call(fromAccount, "sendtoaddress", []any{toAddress, amount, comment, commentTo, subfee})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &txID)
	return
}

// DumpPrivKey return private key as string associated to public <address>
func (b *Bitcoind) DumpPrivKey(account, address string) (privKey string, err error) {
	r, err := b.client.call(account, "dumpprivkey", []any{address})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &privKey)
	return
}

func (b *Bitcoind) ListDescriptors(account string, bb bool) error {
	r, err := b.client.call(account, "listdescriptors", []any{bb})
	if err = handleError(err, &r); err != nil {
		return err
	}
	return nil
}

// GetWalletInfo - Returns an object containing various wallet state info.
// https://bitcoincore.org/en/doc/0.16.0/rpc/wallet/getwalletinfo/
func (b *Bitcoind) GetWalletInfo(account string) (i WalletInfo, err error) {
	r, err := b.client.call(account, "getwalletinfo", nil)
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &i)
	return
}

// ValidateAddressResponse represents a response to "validateaddress" call
type ValidateAddressResponse struct {
	IsValid      bool   `json:"isvalid"`
	Address      string `json:"address"`
	IsMine       bool   `json:"ismine"`
	IsScript     bool   `json:"isscript"`
	PubKey       string `json:"pubkey"`
	IsCompressed bool   `json:"iscompressed"`
	Account      string `json:"account"`
}

// ValidateAddress return information about <bitcoinaddress>.
func (b *Bitcoind) ValidateAddress(address string) (va ValidateAddressResponse, err error) {
	r, err := b.client.call("", "validateaddress", []interface{}{address})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &va)
	return
}
