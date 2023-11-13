$env:MSYS2_PATH_TYPE = "inherit"
$env:MSYSTEM = "UCRT64"
$env:CHERE_INVOKING = 1

# Build the Go project for Windows
Invoke-Expression "C:\msys64\usr\bin\bash.exe -l -c 'GOOS=windows GOARCH=amd64 go build $( $args )'"

# Check for errors
if ($LASTEXITCODE -ne 0)
{
    Write-Host "Error building Go project" -ForegroundColor Red
}