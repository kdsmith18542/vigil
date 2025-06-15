$files = Get-ChildItem -Path 'C:\Users\Keith\vgl' -Recurse -Include *.go,*.md

foreach ($file in $files) {
    try {
        $content = Get-Content $file.FullName -Raw
        $newContent = $content -replace 'github.com/kdsmith18542/vigil', 'github.com/Vigil-Labs/vgl'
        $newContent = $newContent -replace 'github.com/Vigil-Labs/vgl/([a-zA-Z0-9_\-/]+?)/v\d+', 'github.com/Vigil-Labs/vgl/$1'
        $newContent = $newContent -replace 'github.com/kdsmith18542/vigil-Labs/vgl/node/chaincfg/chainhash', 'github.com/kdsmith18542/vigil/chaincfg/chainhash'
        $newContent = $newContent -replace 'github.com/kdsmith18542/vigil-Labs/vigil/chaincfg/chainhash', 'github.com/kdsmith18542/vigil/chaincfg/chainhash'
        $newContent = $newContent -replace 'github.com/kdsmith18542/vigil-Labs/vgl/node/kawpow', 'github.com/kdsmith18542/vigil/kawpow'
        $newContent = $newContent -replace 'github.com/Vigil-Labs/vgl/node/wire', 'github.com/kdsmith18542/vigil/wire'
        Set-Content $file.FullName -Value $newContent -Force
        Write-Host "Processed $($file.FullName)"
    } catch {
        Write-Error "Error processing $($file.FullName): $($_.Exception.Message)"
    }
}