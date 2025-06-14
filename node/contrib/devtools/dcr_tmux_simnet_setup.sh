#!/usr/bin/env bash
#
# Copyright (c) 2020 The Vigil developers
# Use of this source code is governed by an ISC
# license that can be found in the LICENSE file.
#
# Tmux script to create 2 vgld nodes (named vgld1 and vgld2) connected in series
# along with 2 wallets (named wallet1 and wallet2) configured such that wallet1
# is connected via JSON-RPC to vgld1 and, likewise, wallet2 to vgld2.
#
# Both wallet1 and wallet2 use the same seed, however, wallet1 is configured to
# automatically buy tickets and vote, while wallet2 is only configured to vote.
#
# The primary vgld node (vgld1) is configured as the primary mining node.
#
# Network layout:
# vgld1 (p2p: localhost:19555) <-- vgld2 (p2p: localhost:19565)
#
# RPC layout:
# vgld1 (JSON-RPC: localhost:19556)
#     ^--- wallet1 (JSON-RPC: locahost:19557, gRPC: localhost:19558)
# vgld2 (JSON-RPC: localhost:19566)
#     ^--- wallet2 (JSON-RPC: locahost:19567, gRPC: None)

set -e

SESSION="vgld-simnet-nodes"
NODES_ROOT=${VGL_SIMNET_ROOT:-${HOME}/vgldsimnetnodes}
RPCUSER="USER"
RPCPASS="PASS"
WALLET_SEED="b280922d2cffda44648346412c5ec97f429938105003730414f10b01e1402eac"
WALLET_MINING_ADDR="SsXciQNTo3HuV5tX3yy4hXndRWgLMRVC7Ah" # NOTE: This must be changed if the seed is changed.
WALLET_XFER_ADDR="SsWKp7wtdTZYabYFYSc9cnxhwFEjA5g4pFc" # same as above
WALLET_CREATE_CONFIG="y
n
y
${WALLET_SEED}
"
TSPEND_PRIMARY_WIF=PsUUktzTqNKDRudiz3F4Chh5CKqqmp5W3ckRDhwECbwrSuWZ9m5fk
TSPEND_SECONDARY_WIF=PsUVZDkMHvsH8RmYtCxCWs78xsLU9qAyZyLvV9SJWAdoiJxSFhvFx

if [ -d "${NODES_ROOT}" ] ; then
  rm -R "${NODES_ROOT}"
fi

PRIMARY_vgld_NAME=vgld1
SECONDARY_vgld_NAME=vgld2
PRIMARY_WALLET_NAME=wallet1
SECONDARY_WALLET_NAME=wallet2
mkdir -p "${NODES_ROOT}/${PRIMARY_vgld_NAME}"
mkdir -p "${NODES_ROOT}/${SECONDARY_vgld_NAME}"
mkdir -p "${NODES_ROOT}/${PRIMARY_WALLET_NAME}"
mkdir -p "${NODES_ROOT}/${SECONDARY_WALLET_NAME}"

cat > "${NODES_ROOT}/vgld.conf" <<EOF
rpcuser=${RPCUSER}
rpcpass=${RPCPASS}
simnet=1
logdir=./log
datadir=./data
debuglevel=TXMP=debug,MINR=debug
txindex=1
EOF

cat > "${NODES_ROOT}/vglctl.conf" <<EOF
rpcuser=${RPCUSER}
rpcpass=${RPCPASS}
simnet=1
EOF

cat > "${NODES_ROOT}/wallet.conf" <<EOF
username=${RPCUSER}
password=${RPCPASS}
simnet=1
logdir=./log
appdata=./data
pass=123
enablevoting=1
EOF

cd ${NODES_ROOT} && tmux -2 new-session -d -s $SESSION

################################################################################
# Setup the primary vgld node
################################################################################

PRIMARY_vgld_P2P=127.0.0.1:19555
PRIMARY_vgld_RPC=127.0.0.1:19556
tmux rename-window -t $SESSION:0 "${PRIMARY_vgld_NAME}"
tmux split-window -v
tmux select-pane -t 0
tmux send-keys "cd ${NODES_ROOT}/${PRIMARY_vgld_NAME}" C-m
tmux send-keys "vgld -C ../vgld.conf --listen ${PRIMARY_vgld_P2P} --miningaddr=${WALLET_MINING_ADDR}" C-m
tmux resize-pane -D 15
tmux select-pane -t 1
tmux send-keys "cd ${NODES_ROOT}/${PRIMARY_vgld_NAME}" C-m

cat > "${NODES_ROOT}/${PRIMARY_vgld_NAME}/ctl" <<EOF
#!/usr/bin/env bash
vglctl -C ../vglctl.conf "\$@"
EOF
chmod +x "${NODES_ROOT}/${PRIMARY_vgld_NAME}/ctl"

# Script to mine a specified number of blocks with a delay in between them
# Defaults to 1 block
cat > "${NODES_ROOT}/${PRIMARY_vgld_NAME}/mine" <<EOF
#!/usr/bin/env bash
NUM=1
case \$1 in
  ''|*[!0-9]*)  ;;
  *) NUM=\$1 ;;
esac

while [ \$((NUM--)) != 0 ]; do
  ./ctl generate 1
  sleep 1
done
EOF
chmod +x "${NODES_ROOT}/${PRIMARY_vgld_NAME}/mine"
sleep 3
tmux send-keys "./ctl generate 32" C-m

################################################################################
# Setup the primary wallet
################################################################################

PRIMARY_WALLET_RPC=127.0.0.1:19557
PRIMARY_WALLET_GRPC=127.0.0.1:19558
tmux new-window -t $SESSION:1 -n "${PRIMARY_WALLET_NAME}"
tmux split-window -v
tmux select-pane -t 0
tmux resize-pane -D 15
tmux send-keys "cd ${NODES_ROOT}/${PRIMARY_WALLET_NAME}" C-m
tmux send-keys "echo \"${WALLET_CREATE_CONFIG}\" | vglwallet -C ../wallet.conf --create; tmux wait-for -S ${PRIMARY_WALLET_NAME}" C-m
tmux wait-for ${PRIMARY_WALLET_NAME}
tmux send-keys "vglwallet -C ../wallet.conf --enableticketbuyer --ticketbuyer.limit=10" C-m
tmux select-pane -t 1
tmux send-keys "cd ${NODES_ROOT}/${PRIMARY_WALLET_NAME}" C-m

cat > "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/ctl" <<EOF
#!/usr/bin/env bash
vglctl -C ../vglctl.conf --wallet -c ./data/rpc.cert "\$@"
EOF
chmod +x "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/ctl"

# Script to manually purchase tickets
cat > "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/tickets" <<EOF
#!/usr/bin/env bash
NUM=1
case \$1 in
  ''|*[!0-9]*) ;;
  *) NUM=\$1 ;;
esac

./ctl purchaseticket default 999999 1 \`./ctl getnewaddress\` \$NUM
EOF
chmod +x "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/tickets"

# Script to transfer funds with a specified fee rate
# Defaults to a fee rate of 0.0001
cat > "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/xfer" <<EOF
#!/usr/bin/env bash
DEFAULTFEE=0.0001
FEE=\$DEFAULTFEE
case \$1 in
  ''|*[!0-9\.]*) ;;
  *) FEE=\$1 ;;
esac
if [ "\$FEE" != "\$DEFAULTFEE" ]; then
	./ctl settxfee \$FEE
fi
./ctl sendtoaddress ${WALLET_XFER_ADDR} 0.1
if [ "\$FEE" != "\$DEFAULTFEE" ]; then
	./ctl settxfee \$DEFAULTFEE
fi
EOF
chmod +x "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/xfer"

sleep 1
tmux send-keys "./ctl importprivkey ${TSPEND_PRIMARY_WIF} imported false; ./ctl importprivkey ${TSPEND_SECONDARY_WIF} imported false" C-m

################################################################################
# Setup the serially connected secondary vgld node
################################################################################

SECONDARY_vgld_P2P=127.0.0.1:19565
SECONDARY_vgld_RPC=127.0.0.1:19566
cat > "${NODES_ROOT}/${SECONDARY_vgld_NAME}/ctl" <<EOF
#!/usr/bin/env bash
vglctl -C ../vglctl.conf -s ${SECONDARY_vgld_RPC} "\$@"
EOF
chmod +x "${NODES_ROOT}/${SECONDARY_vgld_NAME}/ctl"

cp "${NODES_ROOT}/${PRIMARY_vgld_NAME}/mine" "${NODES_ROOT}/${SECONDARY_vgld_NAME}/"
chmod +x "${NODES_ROOT}/${SECONDARY_vgld_NAME}/mine"

# Script to force reorg
cat > "${NODES_ROOT}/${SECONDARY_vgld_NAME}/reorg" <<EOF
#!/usr/bin/env bash
./ctl node remove ${PRIMARY_vgld_P2P}
./mine
cd "${NODES_ROOT}/${PRIMARY_vgld_NAME}"
./mine 2
cd "${NODES_ROOT}/${SECONDARY_vgld_NAME}"
./ctl node connect ${PRIMARY_vgld_P2P} perm
EOF
chmod +x "${NODES_ROOT}/${SECONDARY_vgld_NAME}/reorg"

tmux new-window -t $SESSION:2 -n "${SECONDARY_vgld_NAME}"
tmux split-window -v
tmux select-pane -t 0
tmux resize-pane -D 15
tmux send-keys "cd ${NODES_ROOT}/${SECONDARY_vgld_NAME}" C-m
tmux send-keys "vgld -C ../vgld.conf --listen ${SECONDARY_vgld_P2P} --rpclisten ${SECONDARY_vgld_RPC} --connect ${PRIMARY_vgld_P2P}  --miningaddr=${WALLET_MINING_ADDR}" C-m
tmux select-pane -t 1
tmux send-keys "cd ${NODES_ROOT}/${SECONDARY_vgld_NAME}" C-m

################################################################################
# Setup the secondary wallet
################################################################################

SECONDARY_WALLET_RPC=127.0.0.1:19567
tmux new-window -t $SESSION:3 -n "${SECONDARY_WALLET_NAME}"
tmux split-window -v
tmux select-pane -t 0
tmux resize-pane -D 15
tmux send-keys "cd ${NODES_ROOT}/${SECONDARY_WALLET_NAME}" C-m
tmux send-keys "echo \"${WALLET_CREATE_CONFIG}\" | vglwallet -C ../wallet.conf --create; tmux wait-for -S ${SECONDARY_WALLET_NAME}" C-m
tmux wait-for ${SECONDARY_WALLET_NAME}
tmux send-keys "vglwallet -C ../wallet.conf --rpcconnect=${SECONDARY_vgld_RPC} --rpclisten=${SECONDARY_WALLET_RPC} --nogrpc" C-m
tmux select-pane -t 1
tmux send-keys "cd ${NODES_ROOT}/${SECONDARY_WALLET_NAME}" C-m

cat > "${NODES_ROOT}/${SECONDARY_WALLET_NAME}/ctl" <<EOF
#!/usr/bin/env bash
vglctl -C ../vglctl.conf -c ./data/rpc.cert -s ${SECONDARY_WALLET_RPC} "\$@"
EOF
chmod +x "${NODES_ROOT}/${SECONDARY_WALLET_NAME}/ctl"

cp "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/tickets" "${NODES_ROOT}/${SECONDARY_WALLET_NAME}/"
chmod +x "${NODES_ROOT}/${SECONDARY_WALLET_NAME}/tickets"

cp "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/xfer" "${NODES_ROOT}/${SECONDARY_WALLET_NAME}/"
chmod +x "${NODES_ROOT}/${SECONDARY_WALLET_NAME}/xfer"

sleep 1
tmux send-keys "./ctl importprivkey ${TSPEND_PRIMARY_WIF} imported false; ./ctl importprivkey ${TSPEND_SECONDARY_WIF} imported false" C-m

################################################################################
# Setup helper script to stop everything
################################################################################

cat > "${NODES_ROOT}/${PRIMARY_vgld_NAME}/stopall" <<EOF
#!/usr/bin/env bash
function countdown {
  secs=\$1
  msg=\$2
  while [ \$secs -gt 0  ]; do
    echo -ne "Seconds \$msg: \$secs\033[0K\r"
    sleep 1
    : \$((secs--))
  done
}

cd "${NODES_ROOT}/${PRIMARY_vgld_NAME}"
./ctl stop 2>/dev/null
cd "${NODES_ROOT}/${PRIMARY_WALLET_NAME}"
./ctl stop 2>/dev/null
cd "${NODES_ROOT}/${SECONDARY_vgld_NAME}"
./ctl stop 2>/dev/null
cd "${NODES_ROOT}/${SECONDARY_WALLET_NAME}"
./ctl stop 2>/dev/null

DELAY=3
countdown \$DELAY "until shutdown"
tmux kill-session -t $SESSION
EOF
cp "${NODES_ROOT}/${PRIMARY_vgld_NAME}/stopall" "${NODES_ROOT}/${SECONDARY_vgld_NAME}/"
cp "${NODES_ROOT}/${PRIMARY_vgld_NAME}/stopall" "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/"
cp "${NODES_ROOT}/${PRIMARY_vgld_NAME}/stopall" "${NODES_ROOT}/${SECONDARY_WALLET_NAME}/"
chmod +x "${NODES_ROOT}/${PRIMARY_vgld_NAME}/stopall"
chmod +x "${NODES_ROOT}/${SECONDARY_vgld_NAME}/stopall"
chmod +x "${NODES_ROOT}/${PRIMARY_WALLET_NAME}/stopall"
chmod +x "${NODES_ROOT}/${SECONDARY_WALLET_NAME}/stopall"

################################################################################
# Attach
################################################################################

tmux select-window -t $SESSION:0
tmux attach-session -t $SESSION
