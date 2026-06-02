# install.ps1
$ErrorActionPreference = 'Stop'

Write-Host "Fetching latest Brisk release..." -ForegroundColor Cyan

$repo = "matyik/brisk"
$releaseUrl = "https://api.github.com/repos/$repo/releases/latest"
$release = Invoke-RestMethod -Uri $releaseUrl

$asset = $release.assets | Where-Object { $_.name -match 'Windows' -and $_.name -match 'x86_64\.zip$' }

if (-not $asset) {
    Write-Error "Could not find a Windows x86_64 release asset."
    exit 1
}

$downloadUrl = $asset.browser_download_url
$installDir = "$HOME\.brisk\bin"
$tempZip = "$env:TEMP\brisk.zip"

if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

Write-Host "Downloading $($asset.name)..."
Invoke-WebRequest -Uri $downloadUrl -OutFile $tempZip
Write-Host "Extracting..."
Expand-Archive -Path $tempZip -DestinationPath $installDir -Force
Remove-Item $tempZip

$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$installDir*") {
    Write-Host "Adding Brisk to User PATH..."
    $newPath = "$installDir;$userPath"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    Write-Host "Please restart your terminal to use the 'brisk' command." -ForegroundColor Yellow
}

Write-Host "Brisk installed successfully to $installDir!" -ForegroundColor Green