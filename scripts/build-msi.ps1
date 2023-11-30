./scripts/build-go.ps1 -tags="d3d,tray" -ldflags -H=windowsgui
go-msi make --msi openstadia.msi --version 0.2.0