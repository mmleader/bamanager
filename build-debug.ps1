$OutputEncoding = [System.Console]::OutputEncoding = [System.Text.Encoding]::UTF8
Write-Host "🚀 Starting DEBUG Build for Windows..."

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
    }
}

# Append to PATH safely
$pathsToAdd = @("$env:GOPATH\bin", "$HOME\go\bin")
foreach ($p in $pathsToAdd) {
    if ($env:PATH -notlike "*$p*") {
        $env:PATH = "$env:PATH;$p"
    }
}

# 2. Compile project with -debug flag
Write-Host "📂 Compiling application in DEBUG mode..."
wails build -debug

if ($LASTEXITCODE -ne 0) {
    Write-Error "❌ Build failed."

    exit $LASTEXITCODE
}

Write-Host ""
Write-Host "✅ DEBUG Build complete!"
Write-Host "The executable can be found in: build\bin\bsmanager.exe"
Write-Host "Run this executable from a terminal to see backend logs."
Write-Host "Right-click anywhere in the app to open Frontend DevTools."
Write-Host ""

