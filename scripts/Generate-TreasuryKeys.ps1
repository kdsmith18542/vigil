# Generate Treasury Keys for Vigil Network
# This script generates a new key pair and outputs the necessary information
# for setting up the treasury in the mainnet parameters.

# Generate a new ECDSA key pair (secp256k1 curve)
$ecParams = [System.Security.Cryptography.ECCurve+NamedCurves]::GetByFriendlyName("secp256k1")
$ecdsa = [System.Security.Cryptography.ECDsa]::Create($ecParams)

# Get private and public key bytes
$privateKey = $ecdsa.ExportECPrivateKey()
$publicKey = $ecdsa.ExportSubjectPublicKeyInfo()

# Create a SHA256 hash of the public key
$sha256 = [System.Security.Cryptography.SHA256]::Create()
$publicKeyHash = $sha256.ComputeHash($publicKey)

# Create a RIPEMD160 hash of the SHA256 hash
$ripemd160 = [System.Security.Cryptography.RIPEMD160]::Create()
$hash160 = $ripemd160.ComputeHash($publicKeyHash)

# Create a P2SH script (OP_HASH160 <hash160> OP_EQUAL)
$script = New-Object byte[] 23
$script[0] = 0xA9       # OP_HASH160
$script[1] = 0x14       # Push 20 bytes (length of hash160)
[System.Buffer]::BlockCopy($hash160, 0, $script, 2, 20)
$script[22] = 0x87     # OP_EQUAL

# Output the results
Write-Host "Vigil Treasury Key Generation"
Write-Host "=============================="
Write-Host "Private Key (hex): $([BitConverter]::ToString($privateKey).Replace('-', ''))"
Write-Host "Public Key (hex):  $([BitConverter]::ToString($publicKey).Replace('-', ''))"
Write-Host "P2SH Script:       $([BitConverter]::ToString($script).Replace('-', ''))"

# Format for mainnetparams.go
Write-Host "`nAdd the following to mainnetparams.go:"
Write-Host "OrganizationPkScript:        []byte{" -NoNewline
for ($i = 0; $i -lt $script.Length; $i++) {
    if ($i -gt 0) { Write-Host ", " -NoNewline }
    Write-Host ("0x{0:x2}" -f $script[$i]) -NoNewline
}
Write-Host "},"
Write-Host "OrganizationPkScriptVersion: 0,"

# Save the private key to a secure file
$privateKey | Set-Content -Path "treasury_private_key.bin" -Encoding Byte -NoNewline
Write-Host "`nPrivate key has been saved to treasury_private_key.bin - KEEP THIS SECURE!"
