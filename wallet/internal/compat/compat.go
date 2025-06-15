package compat

import (
	"github.com/Vigil-Labs/vgl/blockchain/standalone"
"github.com/Vigil-Labs/vgl/hdkeychain"
"github.com/Vigil-Labs/vgl/txscript/stdaddr"
"github.com/Vigil-Labs/vgl/wire"
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




