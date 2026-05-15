# NetForge Remote Installation Script for Windows
# This script downloads the latest release binary from GitHub and installs it.

$GitHubUser = "krishna99-tech" # <--- CHANGE THIS to your GitHub username
$RepoName = "netforge"
$InstallDir = "$HOME\AppData\Local\NetForge"
$BinaryName = "netforge.exe"

Write-Host "--- NetForge Installer ---" -ForegroundColor Cyan

# 1. Create Installation Directory
if (!(Test-Path $InstallDir)) {
    Write-Host "Creating installation directory..." -ForegroundColor Gray
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# 2. Fetch Latest Release from GitHub API
try {
    Write-Host "Fetching latest release info from $GitHubUser/$RepoName..." -ForegroundColor Gray
    $ReleaseUrl = "https://api.github.com/repos/$GitHubUser/$RepoName/releases/latest"
    $Release = Invoke-RestMethod -Uri $ReleaseUrl
    
    # Look for a Windows binary in the release assets
    $Asset = $Release.assets | Where-Object { $_.name -like "*windows-amd64*" -or $_.name -eq "netforge.exe" } | Select-Object -First 1
    
    if ($null -eq $Asset) {
        throw "Could not find a valid Windows binary in the latest GitHub release."
    }

    # 3. Download the Binary
    Write-Host "Downloading $($Asset.name) ($($Release.tag_name))..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri $Asset.browser_download_url -OutFile "$InstallDir\$BinaryName"
    
    Write-Host "Downloaded successfully to $InstallDir" -ForegroundColor Green

} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Falling back to local build if source is available..." -ForegroundColor Yellow
    if (Test-Path "main.go") {
        go build -o "$InstallDir\$BinaryName" main.go
        Write-Host "Local build successful." -ForegroundColor Green
    } else {
        exit 1
    }
}

# 4. Update PATH Environment Variable
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "Updating User PATH..." -ForegroundColor Gray
    $NewPath = "$UserPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    Write-Host "PATH updated successfully!" -ForegroundColor Green
} else {
    Write-Host "NetForge is already in your PATH." -ForegroundColor Yellow
}

Write-Host "`n--- NetForge installed successfully! ---" -ForegroundColor Green
Write-Host "Please RESTART your terminal to start using 'netforge'." -ForegroundColor White
