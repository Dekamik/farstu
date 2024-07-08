#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

apt install unclutter
git clone git@github.com:Dekamik/farstu.git

pushd ./farstu/

make
make install
cp ./example.app.toml ./app.toml

echo 'Installation complete'
echo 'Configure the app.toml and reboot'

popd
