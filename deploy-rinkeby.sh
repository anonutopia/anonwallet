#!/bin/bash

mkdir walletnew

cp anonwallet walletnew/

cp -R templates walletnew/

cp -R public walletnew/

cp -R abi walletnew/

tar -czf walletnew.tar.gz walletnew

scp walletnew.tar.gz kriptokuna@164.132.25.95:

rm -rf walletnew*

ssh -t kriptokuna@164.132.25.95 'tar -xzf walletnew.tar.gz && rm -rf walletnew.tar.gz && cp wallet-rinkeby/config.json walletnew/ && rm -rf wallet-rinkeby && mv walletnew wallet-rinkeby && sudo supervisorctl reload'
