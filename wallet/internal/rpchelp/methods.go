// Copyright (c) 2015 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//go:build !generate

package rpchelp

import (
	"github.com/kdsmith18542/vigil/wallet/rpc/jsonrpc/types"
	vgldtypes "github.com/kdsmith18542/vigil/rpc/jsonrpc/types/v4"
)

// Common return types.
var (
	returnsBool        = []any{(*bool)(nil)}
	returnsNumber      = []any{(*float64)(nil)}
	returnsString      = []any{(*string)(nil)}
	returnsStringArray = []any{(*[]string)(nil)}
	returnsLTRArray    = []any{(*[]types.ListTransactionsResult)(nil)}
)

// Methods contains all methods and result types that help is generated for,
// for every locale.
var Methods = []struct {
	Method      string
	ResultTypes []any
}{
	{"abandontransaction", nil},
	{"accountaddressindex", []any{(*int)(nil)}},
	{"accountsyncaddressindex", nil},
	{"accountunlocked", []any{(*types.AccountUnlockedResult)(nil)}},
	{"addmultisigaddress", returnsString},
	{"addtransaction", nil},
	{"auditreuse", []any{(*map[string][]string)(nil)}},
	{"consolidate", returnsString},
	{"createmultisig", []any{(*types.CreateMultiSigResult)(nil)}},
	{"createnewaccount", nil},
	{"createrawtransaction", returnsString},
	{"createsignature", []any{(*types.CreateSignatureResult)(nil)}},
	{"debuglevel", returnsString},
	{"disapprovepercent", []any{(*uint32)(nil)}},
	{"discoverusage", nil},
	{"dumpprivkey", returnsString},
	{"fundrawtransaction", []any{(*types.FundRawTransactionResult)(nil)}},
	{"getaccount", returnsString},
	{"getaccountaddress", returnsString},
	{"getaddressesbyaccount", returnsStringArray},
	{"getbalance", []any{(*types.GetBalanceResult)(nil)}},
	{"getbestblock", []any{(*vgldtypes.GetBestBlockResult)(nil)}},
	{"getbestblockhash", returnsString},
	{"getblockcount", returnsNumber},
	{"getblockhash", returnsString},
	{"getblockheader", []any{(*vgldtypes.GetBlockHeaderVerboseResult)(nil)}},
	{"getblock", []any{(*vgldtypes.GetBlockVerboseResult)(nil)}},
	{"getcoinjoinsbyacct", []any{(*map[string]uint32)(nil)}},
	{"getcurrentnet", []any{(*uint32)(nil)}},
	{"getinfo", []any{(*types.InfoWalletResult)(nil)}},
	{"getmasterpubkey", []any{(*string)(nil)}},
	{"getmultisigoutinfo", []any{(*types.GetMultisigOutInfoResult)(nil)}},
	{"getnewaddress", returnsString},
	{"getpeerinfo", []any{(*types.GetPeerInfoResult)(nil)}},
	{"getrawchangeaddress", returnsString},
	{"getreceivedbyaccount", returnsNumber},
	{"getreceivedbyaddress", returnsNumber},
	{"getstakeinfo", []any{(*types.GetStakeInfoResult)(nil)}},
	{"gettickets", []any{(*types.GetTicketsResult)(nil)}},
	{"gettransaction", []any{(*types.GetTransactionResult)(nil)}},
	{"gettxout", []any{(*vgldtypes.GetTxOutResult)(nil)}},
	{"getunconfirmedbalance", returnsNumber},
	{"getvotechoices", []any{(*types.GetVoteChoicesResult)(nil)}},
	{"getwalletfee", returnsNumber},
	{"getcfilterv2", []any{(*types.GetCFilterV2Result)(nil)}},
	{"help", append(returnsString, returnsString[0])},
	{"importcfiltersv2", nil},
	{"importprivkey", nil},
	{"importpubkey", nil},
	{"importscript", nil},
	{"importxpub", nil},
	{"listaccounts", []any{(*map[string]float64)(nil)}},
	{"listaddresstransactions", returnsLTRArray},
	{"listalltransactions", returnsLTRArray},
	{"listlockunspent", []any{(*[]vgldtypes.TransactionInput)(nil)}},
	{"listreceivedbyaccount", []any{(*[]types.ListReceivedByAccountResult)(nil)}},
	{"listreceivedbyaddress", []any{(*[]types.ListReceivedByAddressResult)(nil)}},
	{"listsinceblock", []any{(*types.ListSinceBlockResult)(nil)}},
	{"listtransactions", returnsLTRArray},
	{"listunspent", []any{(*types.ListUnspentResult)(nil)}},
	{"lockaccount", nil},
	{"lockunspent", returnsBool},
	{"mixaccount", nil},
	{"mixoutput", nil},
	{"processunmanagedticket", nil},
	{"purchaseticket", returnsString},
	{"redeemmultisigout", []any{(*types.RedeemMultiSigOutResult)(nil)}},
	{"redeemmultisigouts", []any{(*types.RedeemMultiSigOutResult)(nil)}},
	{"renameaccount", nil},
	{"rescanwallet", nil},
	{"sendfrom", returnsString},
	{"sendfromtreasury", returnsString},
	{"sendmany", returnsString},
	{"sendrawtransaction", returnsString},
	{"sendtoaddress", returnsString},
	{"sendtomultisig", returnsString},
	{"sendtotreasury", returnsString},
	{"setaccountpassphrase", nil},
	{"setdisapprovepercent", nil},
	{"settreasurypolicy", nil},
	{"settspendpolicy", nil},
	{"settxfee", returnsBool},
	{"setvotechoice", nil},
	{"signmessage", returnsString},
	{"signrawtransaction", []any{(*types.SignRawTransactionResult)(nil)}},
	{"signrawtransactions", []any{(*types.SignRawTransactionsResult)(nil)}},
	{"spendoutputs", returnsString},
	{"sweepaccount", []any{(*types.SweepAccountResult)(nil)}},
	{"syncstatus", []any{(*types.SyncStatusResult)(nil)}},
	{"ticketinfo", []any{(*[]types.TicketInfoResult)(nil)}},
	{"treasurypolicy", []any{(*[]types.TreasuryPolicyResult)(nil), (*types.TreasuryPolicyResult)(nil)}},
	{"tspendpolicy", []any{(*[]types.TSpendPolicyResult)(nil), (*types.TSpendPolicyResult)(nil)}},
	{"unlockaccount", nil},
	{"validateaddress", []any{(*types.ValidateAddressWalletResult)(nil)}},
	{"validatepreVGLP0005cf", returnsBool},
	{"verifymessage", returnsBool},
	{"version", []any{(*map[string]vgldtypes.VersionResult)(nil)}},
	{"walletinfo", []any{(*types.WalletInfoResult)(nil)}},
	{"walletislocked", returnsBool},
	{"walletlock", nil},
	{"walletpassphrase", nil},
	{"walletpassphrasechange", nil},
	{"walletpubpassphrasechange", nil},
}

// HelpDescs contains the locale-specific help strings along with the locale.
var HelpDescs = []struct {
	Locale   string // Actual locale, e.g. en_US
	GoLocale string // Locale used in Go names, e.g. EnUS
	Descs    map[string]string
}{
	{"en_US", "EnUS", helpDescsEnUS}, // helpdescs_en_US.go
}
