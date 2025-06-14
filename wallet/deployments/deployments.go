// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package deployments

import (
	"context"

	"github.com/vigilnetwork/vgl/wallet/errors"
	"github.com/vigilnetwork/vgl/chaincfg/v3"
	vgldtypes "github.com/vigilnetwork/vgl/rpc/jsonrpc/types/v4"
	"github.com/vigilnetwork/vgl/wire"
)

// HardcodedDeployment specifies hardcoded block heights that a deployment
// activates at.  If the value is negative, the deployment is either inactive or
// can't be determined due to the uniqueness properties of the network.
//
// Since these are hardcoded deployments, and cannot support every possible
// network, conditional logic should only be applied when a deployment is
// active, not when it is inactive.
type HardcodedDeployment struct {
	MainNetActivationHeight  int32
	TestNet2ActivationHeight int32
	TestNet3ActivationHeight int32
	SimNetActivationHeight   int32
}

// VGLP0001 specifies hard forking changes to the stake difficulty algorithm as
// defined by https://github.com/Vigil/VGLPs/blob/master/VGLP-0001/VGLP-0001.mediawiki.
var VGLP0001 = HardcodedDeployment{
	MainNetActivationHeight:  149248,
	TestNet2ActivationHeight: 46128,
	TestNet3ActivationHeight: 0,
	SimNetActivationHeight:   0,
}

// VGLP0002 specifies the activation of the OP_SHA256 hard fork as defined by
// https://github.com/Vigil/VGLPs/blob/master/VGLP-0002/VGLP-0002.mediawiki.
var VGLP0002 = HardcodedDeployment{
	MainNetActivationHeight:  189568,
	TestNet2ActivationHeight: 151968,
	TestNet3ActivationHeight: 0,
	SimNetActivationHeight:   0,
}

// VGLP0003 specifies the activation of a CSV soft fork as defined by
// https://github.com/Vigil/VGLPs/blob/master/VGLP-0003/VGLP-0003.mediawiki.
var VGLP0003 = HardcodedDeployment{
	MainNetActivationHeight:  189568,
	TestNet2ActivationHeight: 151968,
	TestNet3ActivationHeight: 0,
	SimNetActivationHeight:   0,
}

// Active returns whether the hardcoded deployment is active at height on the
// network specified by params.  Active always returns false for unrecognized
// networks.
func (d *HardcodedDeployment) Active(height int32, net wire.CurrencyNet) bool {
	var activationHeight int32 = -1
	switch net {
	case wire.MainNet:
		activationHeight = d.MainNetActivationHeight
	case 0x48e7a065: // testnet2
		activationHeight = d.TestNet2ActivationHeight
	case wire.TestNet3:
		activationHeight = d.TestNet3ActivationHeight
	case wire.SimNet:
		activationHeight = d.SimNetActivationHeight
	}
	return activationHeight >= 0 && height >= activationHeight
}

const (
	lockedinStatus = vgldtypes.AgendaInfoStatusLockedIn
	activeStatus   = vgldtypes.AgendaInfoStatusActive
)

// Querier defines the interface for a chain backend that can (trustfully)
// query for agenda deployment information.
type Querier interface {
	// Deployments should return information about existing agendas,
	// including their deployment status.
	Deployments(context.Context) (map[string]vgldtypes.AgendaInfo, error)
}

// VGLP0010Active returns whether the consensus rules for the next block with the
// current chain tip height requires the subsidy split as specified in VGLP0010.
// VGLP0010 is always active on simnet, and requires the RPC syncer to detect
// activation on mainnet and testnet3.
func VGLP0010Active(ctx context.Context, height int32, params *chaincfg.Params,
	querier Querier) (bool, error) {

	net := params.Net
	rcai := int32(params.RuleChangeActivationInterval)

	if net == wire.SimNet {
		return true, nil
	}
	if net != wire.MainNet && net != wire.TestNet3 {
		return false, nil
	}
	if querier == nil {
		return false, errors.E(errors.Bug, "VGLP0010 activation check requires RPC syncer")
	}
	deployments, err := querier.Deployments(ctx)
	if err != nil {
		return false, err
	}
	d, ok := deployments[chaincfg.VoteIDChangeSubsidySplit]
	if !ok {
		return false, nil
	}
	switch {
	case d.Status == lockedinStatus && height == int32(d.Since)+rcai-1:
		return true, nil
	case d.Status == activeStatus:
		return true, nil
	default:
		return false, nil
	}
}

// VGLP0012Active returns whether the consensus rules for the next block with the
// current chain tip height requires the version 2 subsidy split as specified in
// VGLP0012.  VGLP0012 requires the RPC syncer to detect activation on mainnet,
// testnet3 and simnet.
func VGLP0012Active(ctx context.Context, height int32, params *chaincfg.Params,
	querier Querier) (bool, error) {

	net := params.Net
	rcai := int32(params.RuleChangeActivationInterval)

	if net != wire.MainNet && net != wire.TestNet3 && net != wire.SimNet {
		return false, nil
	}
	if querier == nil {
		return false, errors.E(errors.Bug, "VGLP0012 activation check requires RPC syncer")
	}
	deployments, err := querier.Deployments(ctx)
	if err != nil {
		return false, err
	}
	d, ok := deployments[chaincfg.VoteIDChangeSubsidySplitR2]
	if !ok {
		return false, nil
	}
	switch {
	case d.Status == lockedinStatus && height == int32(d.Since)+rcai-1:
		return true, nil
	case d.Status == activeStatus:
		return true, nil
	default:
		return false, nil
	}
}
