# Build script for NetForge CLI

$ProjectName = "netforge"
$OutputDirectory = "bin"

if (!(Test-Path $OutputDirectory)) {
    New-Item -ItemType Directory -Path $OutputDirectory
}

# Local Build (for your current PC)
Write-Host "Building for Local Machine..." -ForegroundColor Green
$env:GOOS = ""
$env:GOARCH = ""
go build -o "$OutputDirectory/$ProjectName.exe" main.go

# Windows amd64 (64-bit Intel/AMD)
Write-Host "Building for Windows (amd64)..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o "$OutputDirectory/$ProjectName-windows-amd64.exe" main.go

# Windows 386 (32-bit)
Write-Host "Building for Windows (386)..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:GOARCH = "386"
go build -o "$OutputDirectory/$ProjectName-windows-386.exe" main.go

# Windows arm64 (Windows on ARM)
Write-Host "Building for Windows (arm64)..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:GOARCH = "arm64"
go build -o "$OutputDirectory/$ProjectName-windows-arm64.exe" main.go

# Linux
Write-Host "Building for Linux (amd64)..." -ForegroundColor Cyan
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o "$OutputDirectory/$ProjectName-linux-amd64" main.go

# macOS (Intel)
Write-Host "Building for macOS (amd64)..." -ForegroundColor Cyan
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o "$OutputDirectory/$ProjectName-darwin-amd64" main.go

# macOS (Apple Silicon)
Write-Host "Building for macOS (arm64)..." -ForegroundColor Cyan
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o "$OutputDirectory/$ProjectName-darwin-arm64" main.go

Write-Host "Builds completed! Check the '$OutputDirectory' folder." -ForegroundColor Green
