#! /bin/bash
# Ensure git is installed
gitInstalled=`which git`

if [ -z "$gitInstalled" ]
then
    echo "Git not found installing now..."
    apt update
    apt upgrade -y
    apt install git-all -y
else
    echo "Git Installed"
fi

# Ensure golang is installed
goInstalled=`which go`
if [ -z "$goInstalled" ]
then
    echo "Go not found installing now..."
    apt update
    apt upgrade -y
    apt install wget
    wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz

    tar -xvf go1.13.3.linux-amd64.tar.gz
    mv go /usr/local
    
    export GOROOT=/usr/local/go
    export PATH=$PATH:$GOROOT/bin
else
    echo "Go Installed"
fi

# Clone repo
mkdir tempBuild
cd tempBuild
git clone https://github.com/JeffreyRiggle/sona-server

# Build src
cd sona-server/src
go get -v -d -t ./...
go build -v .

# Copy output
mv src ../../../

# Cleanup
rm -r tempBuild