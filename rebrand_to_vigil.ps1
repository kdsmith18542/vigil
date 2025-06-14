# This script automates the rebranding of 'Vigil' to 'Vigil' within the specified directory.
# This script automates the rebranding of "Vigil" to "Vigil" in Go source files.
# It replaces text occurrences, updates import paths, and modifies copyright notices.

param(
    [string]$TargetPath = ""
)

function Replace-Content {
    param(
        [string]$FilePath,
        [string]$OldString,
        [string]$NewString
    )
    (Get-Content $FilePath) | ForEach-Object { $_ -replace $OldString, $NewString } | Set-Content $FilePath
}

if (-not $TargetPath) {
    Write-Host "Usage: .\rebrand_to_vigil.ps1 -TargetPath <path_to_directory>"
    exit 1
}

Write-Host "Starting rebranding process in: $TargetPath"

# Define common replacements
$replacements = @(
    @{Pattern = "(?i)// Copyright \(c\) \d{4}(-\d{4})? The Vigil developers"; Replacement = "// Copyright (c) 2024 The Vigil developers"},
    @{Pattern = "vigil.network/vgl/cspp/v2"; Replacement = "vigil.network/vgl/cspp/v2"},
    @{Pattern = "vigil.network/vgl/wallet/v5/rpc/walletrpc"; Replacement = "vigil.network/vgl/wallet/v5/rpc/walletrpc"},
    @{Pattern = "vigil.network/vgl/wallet/v\d"; Replacement = "vigil.network/vgl/wallet"},
    @{Pattern = "vigil.network/vgl/wallet"; Replacement = "vigil.network/vgl/wallet"},
    @{Pattern = "github.com/vigilnetwork/vgl"; Replacement = "github.com/vigilnetwork/vgl"},
    @{Pattern = "github.com/vigilnetwork/vgl/slog"; Replacement = "github.com/vigilnetwork/vgl/slog"},
    @{Pattern = "github.com/vigilnetwork/vgl/vspd"; Replacement = "github.com/vigilnetwork/vgl/vspd"},
    @{Pattern = "Vigil signed message"; Replacement = "Vigil signed message"},
    @{Pattern = "vigil coin amounts"; Replacement = "vigil coin amounts"},
    @{Pattern = "data.vigil.network"; Replacement = "data.vigil.network"},
    @{Pattern = "vigil.network"; Replacement = "vigil.network"},
    @{Pattern = "Vigil"; Replacement = "Vigil"},
    @{Pattern = "Vigil"; Replacement = "vigil"},
    @{Pattern = "VGLP0005"; Replacement = "VGLP0005"},
    @{Pattern = "vglwallet"; Replacement = "vglwallet"},
    @{Pattern = "Vigiliton"; Replacement = "Vigilant"},
    @{Pattern = "vglctl"; Replacement = "vglctl"},
    @{Pattern = "vgld"; Replacement = "vgld"},
    @{Pattern = "run_vgl_tests.sh"; Replacement = "run_vgl_tests.sh"},
    @{Pattern = "Vigil-release"; Replacement = "vigil-release"},
    @{Pattern = "Vigil developers"; Replacement = "Vigil developers"},
    @{Pattern = "vigil.network"; Replacement = "vigil.network"},
    @{Pattern = "Vigil"; Replacement = "vigil"},
    @{Pattern = "Vigil"; Replacement = "Vigil"},
    @{Pattern = "vglwallet"; Replacement = "vglwallet"},
    @{Pattern = "vglctl"; Replacement = "vglctl"},
    @{Pattern = "vgld"; Replacement = "vgld"},
    @{Pattern = "VGLP"; Replacement = "VGLP"},
    @{Pattern = "Vigil Full Node"; Replacement = "Vigil Full Node"},
    @{Pattern = "jcv@vigil.network"; Replacement = "jcv@vigil.network"},
    @{Pattern = "cjepson@vigil.network"; Replacement = "cjepson@vigil.network"},
    @{Pattern = "ay-p@vigil.network"; Replacement = "ay-p@vigil.network"},
    @{Pattern = "A full-node Vigil implementation written in Go"; Replacement = "A full-node Vigil implementation written in Go"},
    @{Pattern = "VGL"; Replacement = "VGL"},
    @{Pattern = "Vigiliteia"; Replacement = "Vigiliteia"},
    @{Pattern = "vgldata"; Replacement = "vgldata"},
    @{Pattern = "VGLpool"; Replacement = "vglpool"},
    @{Pattern = "Vigil/vgld"; Replacement = "vigilnetwork/vgl"},
    @{Pattern = "Vigil/vglctl"; Replacement = "vigilnetwork/vglctl"},
    @{Pattern = "Vigil"; Replacement = "vigil"},
    @{Pattern = "Vigil-node"; Replacement = "vigil-node"},
    @{Pattern = "Vigil genesis block"; Replacement = "Vigil genesis block"},
    @{Pattern = "Faster than Vigil"; Replacement = "Faster than Vigil"},
    @{Pattern = "make Vigil one of the most profitable KawPoW coins"; Replacement = "make Vigil one of the most profitable KawPoW coins"},
    @{Pattern = "ensuring Vigil's long-term survival"; Replacement = "ensuring Vigil's long-term survival"},
    @{Pattern = "Vigil Wallet and Vigil Explorer"; Replacement = "Vigil Wallet and Vigil Explorer"},
    @{Pattern = "Vigil Launcher"; Replacement = "Vigil Launcher"},
    @{Pattern = "official Vigil Explorer"; Replacement = "official Vigil Explorer"},
    @{Pattern = "official vigil.network website"; Replacement = "official vigil.network website"},
    @{Pattern = "What is Vigil?"; Replacement = "What is Vigil?"},
    @{Pattern = "Vigil Testnet wallet/node"; Replacement = "Vigil Testnet wallet/node"},
    @{Pattern = "audit of the Vigil codebase"; Replacement = "audit of the Vigil codebase"},
    @{Pattern = "official launch versions of the Vigil software"; Replacement = "official launch versions of the Vigil software"},
    @{Pattern = "mining pool for Vigil"; Replacement = "mining pool for Vigil"},
    @{Pattern = "Fork of Vigil's VGLpool"; Replacement = "Fork of Vigil's vglpool"},
    @{Pattern = "native compatibility with Vigil's hybrid consensus"; Replacement = "native compatibility with Vigil's hybrid consensus"},
    @{Pattern = "expand the Vigil ecosystem"; Replacement = "expand the Vigil ecosystem"},
    @{Pattern = "ExecStart=/home/vgld/Vigil/vgld"; Replacement = "ExecStart=/home/vgld/vigil/vgld"},
    @{Pattern = "/opt/Vigil/bin/vgld"; Replacement = "/opt/vigil/bin/vgld"},
    @{Pattern = "/var/vgld"; Replacement = "/var/vgld"},
    @{Pattern = "John C. Vernaleo <jcv@vigil.network>"; Replacement = "John C. Vernaleo <jcv@vigil.network>"},
    @{Pattern = "cjepson <cjepson@vigil.network>"; Replacement = "cjepson <cjepson@vigil.network>"},
    @{Pattern = "Alex Yocom-Piatt <ay-p@vigil.network>"; Replacement = "Alex Yocom-Piatt <ay-p@vigil.network>"},
    @{Pattern = "Description=Vigil Full Node"; Replacement = "Description=Vigil Full Node"},
    @{Pattern = "ENV USER=Vigil"; Replacement = "ENV USER=vigil"},
    @{Pattern = "WORKDIR /go/src/github.com/vigilnetwork/vgl"; Replacement = "WORKDIR /go/src/github.com/vigilnetwork/vgl"},
    @{Pattern = "RUN git clone --branch ${vgld_BUILD_TAG} -c advice.detachedHead=false https://github.com/vigilnetwork/vgl . && \\"; Replacement = "RUN git clone --branch ${VGL_BUILD_TAG} -c advice.detachedHead=false https://github.com/vigilnetwork/vgl . && \\"},
    @{Pattern = "WORKDIR /go/src/github.com/Vigil/vglctl"; Replacement = "WORKDIR /go/src/github.com/vigilnetwork/vglctl"},
    @{Pattern = "RUN git clone --branch ${vglctl_BUILD_TAG} -c advice.detachedHead=false https://github.com/Vigil/vglctl . && \\"; Replacement = "RUN git clone --branch ${VGLCTL_BUILD_TAG} -c advice.detachedHead=false https://github.com/vigilnetwork/vglctl . && \\"},
    @{Pattern = "WORKDIR /go/src/github.com/vigilnetwork/vgl/contrib/docker/entrypoint"; Replacement = "WORKDIR /go/src/github.com/vigilnetwork/vgl/contrib/docker/entrypoint"},
    @{Pattern = "ENV Vigil_DATA=/home/Vigil"; Replacement = "ENV VIGIL_DATA=/home/vigil"},
    @{Pattern = "COPY --from=builder --chown=Vigil /emptydatadir /tmp"; Replacement = "COPY --from=builder --chown=vigil /emptydatadir /tmp"},
    @{Pattern = "USER Vigil"; Replacement = "USER vigil"},
    @{Pattern = "#VOLUME [ `"/home/Vigil`" ]"; Replacement = "#VOLUME [ `"/home/vigil`" ]"},
    @{Pattern = "A full-node Vigil implementation written in Go"; Replacement = "A full-node Vigil implementation written in Go"},
    @{Pattern = "Full KawPoW Integration (Vigil-node)"; Replacement = "Full KawPoW Integration (vigil-node)"},
    @{Pattern = "Finalized and coded economic parameters for the Vigil genesis block."; Replacement = "Finalized and coded economic parameters for the Vigil genesis block."},
    @{Pattern = "Block Time: ~2.5 minutes (Faster than Vigil for a better user experience)."; Replacement = "Block Time: ~2.5 minutes (Faster than Vigil for a better user experience)."},
    @{Pattern = "50% (10 VGL) to PoW Miners (KawPoW): (Viral Launch) An aggressive share to make Vigil one of the most profitable KawPoW coins at launch, bootstrapping network hashrate and attracting a large initial community."; Replacement = "50% (10 VGL) to PoW Miners (KawPoW): (Viral Launch) An aggressive share to make Vigil one of the most profitable KawPoW coins at launch, bootstrapping network hashrate and attracting a large initial community."},
    @{Pattern = "10% (2 VGL) to Vigil Treasury (Vigiliteia): (Sustainability) A self-funding mechanism, controlled by stakeholders, to pay for all future development, marketing, and ecosystem growth, ensuring Vigil's long-term survival."; Replacement = "10% (2 VGL) to Vigil Treasury (Vigiliteia): (Sustainability) A self-funding mechanism, controlled by stakeholders, to pay for all future development, marketing, and ecosystem growth, ensuring Vigil's long-term survival."},
    @{Pattern = "Fully rebrand the vglwallet and vgldata forks to Vigil Wallet and Vigil Explorer."; Replacement = "Fully rebrand the vglwallet and vgldata forks to Vigil Wallet and Vigil Explorer."},
    @{Pattern = "Develop a one-click ""Vigil Launcher"" application that bundles the node and wallet for less technical users."; Replacement = "Develop a one-click ""Vigil Launcher"" application that bundles the node and wallet for less technical users."},
    @{Pattern = "Task 2.4: [VIRAL FEATURE] The ""Vigil Staking"" Dashboard"; Replacement = "Task 2.4: [VIRAL FEATURE] The ""Vigil Staking"" Dashboard"},
    @{Pattern = "Concept: A simple, web-based dashboard integrated into the official Vigil Explorer."; Replacement = "Concept: A simple, web-based dashboard integrated into the official Vigil Explorer."},
    @{Pattern = "Launch the official vigil.network website and open the Discord to the public."; Replacement = "Launch the official vigil.network website and open the Discord to the public."},
    @{Pattern = "Publish a series of educational blog posts and Twitter threads: `"What is Vigil?`", `"Why KawPoW?`", `"How PoS Secures Your Coins,`" `"Understanding Vigiliteia.`" "; Replacement = "Simplified Replacement String"},
    @{Pattern = "Easy-to-use installers for the Vigil Testnet wallet/node."; Replacement = "Easy-to-use installers for the Vigil Testnet wallet/node."},
    @{Pattern = "Security Audit: Commission and publish a full third-party audit of the Vigil codebase."; Replacement = "Security Audit: Commission and publish a full third-party audit of the Vigil codebase."},
    @{Pattern = "Launch Binaries: Compile, sign, and prepare the official launch versions of the Vigil software."; Replacement = "Launch Binaries: Compile, sign, and prepare the official launch versions of the Vigil software."},
    @{Pattern = "Objective: To provide a stable, trusted, and high-performance official mining pool for Vigil at launch."; Replacement = "Objective: To provide a stable, trusted, and high-performance official mining pool for Vigil at launch."},
    @{Pattern = "Pool Engine: Fork of Vigil's VGLpool (Go). This is critical for native compatibility with Vigil's hybrid consensus and treasury rules. The core development task is to modify it to support the KawPoW algorithm."; Replacement = "Pool Engine: Fork of Vigil's vglpool (Go). This is critical for native compatibility with Vigil's hybrid consensus and treasury rules. The core development task is to modify it to support the KawPoW algorithm."},
    @{Pattern = "Launch Time: Publish the genesis block hash to all official channels. Release the final software on GitHub. Launch the official Vigil mining pool website."; Replacement = "Launch Time: Publish the genesis block hash to all official channels. Release the final software on GitHub. Launch the official Vigil mining pool website."},
    @{Pattern = "Concept: A formal process where the Treasury actively encourages and funds community-led projects that expand the Vigil ecosystem."; Replacement = "Concept: A formal process where the Treasury actively encourages and funds community-led projects that expand the Vigil ecosystem."}
)


# Function to check if a file contains any of the old patterns
function Test-FileForOldPatterns($FilePath) {
    foreach ($r in $replacements) {
        $pattern = $r.Pattern
        if (Select-String -Path $FilePath -Pattern $pattern -Quiet -ErrorAction SilentlyContinue) {
            return $true
        }
    }
    return $false
}

$filesToProcess = Get-ChildItem -Path $TargetPath -Recurse -Include @('*.go', '*.md', '*.txt', '*.ps1', '*.sum', '*.conf', '*.yml', '*.toml', '*.json', '*.sh', '*.bat', '*.html', '*.css', '*.js', '*.xml', '*.proto', '*.s', '*.c', '*.h', '*.cpp', '*.hpp', '*.rs', '*.java', '*.kt', '*.py', '*.rb', '*.php', '*.pl', '*.pm', '*.sql', '*.vue', '*.ts', '*.tsx', '*.jsx')

# Process all files
foreach ($file in $filesToProcess) {
    $filePath = $file.FullName
    if (Test-FileForOldPatterns $filePath) {
        Write-Host "Processing file: $filePath"
        foreach ($r in $replacements) {
            Replace-Content -FilePath $filePath -OldString $r.Pattern -NewString $r.Replacement
        }
    } else {
        Write-Host "Skipping file (already rebranded): $filePath"
    }
}

Write-Host "Rebranding process completed successfully!"
