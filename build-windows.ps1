Write-Host "🚀 Starting Build for Windows..."

# 1. Setup Environment
Write-Host "🔍 Detecting Go environment..."
try {
    $env:GOPATH = go env GOPATH
} catch {
    $env:GOPATH = ""
}

if (-not $env:GOPATH) {
    if (Test-Path "$HOME\go") {
        $env:GOPATH = "$HOME\go"
    } else {
        Write-Warning "⚠️ GOPATH not found and $HOME\go does not exist."
    }
}

# Append to PATH safely
$pathsToAdd = @("$env:GOPATH\bin", "$HOME\go\bin")
foreach ($p in $pathsToAdd) {
    if ($env:PATH -notlike "*$p*") {
        $env:PATH = "$env:PATH;$p"
    }
}

# 2. Check for Wails
if (-not (Get-Command wails -ErrorAction SilentlyContinue)) {
    Write-Error "❌ Error: 'wails' command not found even after searching GOPATH\bin."
    Write-Host "Please ensure Wails is installed (go install github.com/wailsapp/wails/v2/cmd/wails@latest)"
    Read-Host "Press Enter to exit..."
    exit 1
}

# 3. Compile project
Write-Host "📂 Compiling application..."
wails build

if ($LASTEXITCODE -ne 0) {
    Write-Error "❌ Build failed."
    Read-Host "Press Enter to exit..."
    exit $LASTEXITCODE
}

Write-Host ""
Write-Host "✅ Build complete!"
Write-Host "The executable can be found in: build\bin\bsmanager.exe"
Write-Host ""
Read-Host "Press Enter to exit..."
