package compat

import (
	"github.com/kdsmith18542/vigil/blockchain/standalone/v2"
	"github.com/kdsmith18542/vigil/hdkeychain/v3"
	"github.com/kdsmith18542/vigil/txscript/v4/stdaddr"
	"github.com/kdsmith18542/vigil/wire"
)

func HD2Address(k *hdkeychain.ExtendedKey, params stdaddr.AddressParams) (*stdaddr.AddressPubKeyHashEcdsaSecp256k1V0, error) {
	pk := k.SerializedPubKey()
	hash := stdaddr.Hash160(pk)
	return stdaddr.NewAddressPubKeyHashEcdsaSecp256k1V0(hash, params)
}

// IsEitherCoinBaseTx verifies if a transaction is either a coinbase prior to
// the treasury agenda activation or a coinbse after treasury agenda
// activation.
func IsEitherCoinBaseTx(tx *wire.MsgTx) bool {
	if standalone.IsCoinBaseTx(tx, false) {
		return true
	}
	if standalone.IsCoinBaseTx(tx, true) {
		return true
	}
	return false
}
