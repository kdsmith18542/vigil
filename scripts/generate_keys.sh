#!/bin/bash

# Generate a new private key
openssl ecparam -name secp256k1 -genkey -noout -out private.pem

# Extract the public key
openssl ec -in private.pem -pubout -out public.pem

# Get the private key in hex format
PRIVKEY_HEX=$(openssl ec -in private.pem -text -noout 2>/dev/null | grep priv -A 3 | tail -n +2 | tr -d '\n ' | sed 's/://g' | tr '[:upper:]' '[:lower:]')

# Get the public key in compressed format
PUBKEY_HEX=$(openssl ec -in private.pem -pubout -text 2>/dev/null | grep -A 5 "pub:" | tail -n +2 | tr -d '\n ' | sed 's/://g' | tr '[:upper:]' '[:lower:]')

# Create a SHA-256 hash of the public key
PUBKEY_HASH=$(echo -n $PUBKEY_HEX | xxd -r -p | openssl dgst -sha256 -binary | xxd -p -c 32)

# Create a RIPEMD160 hash of the SHA-256 hash
RIPEMD160_HASH=$(echo -n $PUBKEY_HASH | xxd -r -p | openssl dgst -rmd160 -binary | xxd -p -c 20)

# Create a P2SH script (OP_HASH160 <hash> OP_EQUAL)
P2SH_SCRIPT="a914${RIPEMD160_HASH}87"

# Output the results
echo "Vigil Treasury Key Generation"
echo "=============================="
echo "Private Key (hex): $PRIVKEY_HEX"
echo "Public Key (hex):  $PUBKEY_HEX"
echo "P2SH Script:       $P2SH_SCRIPT"
echo
echo "Add the following to mainnetparams.go:"
echo -n "OrganizationPkScript:        []byte{"
for i in $(seq 0 2 $((${#P2SH_SCRIPT}-2))); do
    if [ $i -gt 0 ]; then
        echo -n ", "
    fi
    echo -n "0x${P2SH_SCRIPT:$i:2}"
done
echo "}"
echo "OrganizationPkScriptVersion: 0,"

# Save the private key to a secure file
echo -e "\nPrivate key has been saved to private.pem - KEEP THIS SECURE!"
