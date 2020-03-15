# Argument 0 is the destination path. If this path does not exist the server will be deployed to Documents/sona-server
$Dest = $args[0]

if (-not ($Dest))
{
    $Dest = Join-Path ([Environment]::GetFolderPath('MyDocuments')) '/sona-server'
}

if (-not Test-Path -Path $Dest)
{
    New-Item $Dest -ItemType 'directory'
}

# Ensure Chocolaty is installed
$ChocoInstalled = choco -v

if (-not ($ChocoInstalled))
{
    Write-Output "Choco is not installed, Installing now..."
    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
}
else
{
    Write-Output "Choco is Installed"
}

# Ensure git is installed.
$GitInstalled = git version

if (-not $GitInstalled)
{
    Write-Output "Git is Installed"
}
else
{
    Write-Output "Git is not Installed, installing now..."
    choco install git
}

# Ensure go is installed.
$GoInstalled = go version

if (-not ($GoInstalled))
{
    Write-Output "Go is not installed, Installing now..."
    choco install golang
}
else
{
    Write-Output "Go is Installed"
}

# Clone repo
New-Item -Path $env:Temp -Name 'sonabuild' -ItemType 'directory'
$BuildArea = Join-Path $env:Temp 'sonabuild' 

Set-Location -Path $BuildArea
git clone https://github.com/JeffreyRiggle/sona-server

# Build application
$SRCDIR = Join-Path $BuildArea '/sona-server/src'
Set-Location -Path $SRCDIR

go get -v -d -t ./...
go build -v .

# Copy output
Move-Item -Path (Join-Path $SRCDIR 'src.exe') -Destination (Join-Path $Dest 'src.exe')
Set-Location $Dest

# Cleanup
Remove-Item $BuildArea -Force

Write-Output "Sona Server can be found at $Dest"