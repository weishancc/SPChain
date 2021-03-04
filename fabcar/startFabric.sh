#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
CC_SRC_LANGUAGE=${1:-"go"}
CC_SRC_LANGUAGE=`echo "$CC_SRC_LANGUAGE" | tr [:upper:] [:lower:]`
if [ "$CC_SRC_LANGUAGE" = "go" -o "$CC_SRC_LANGUAGE" = "golang"  ]; then
	CC_RUNTIME_LANGUAGE=golang
	CC_SRC_PATH=github.com/chaincode/artworks
elif [ "$CC_SRC_LANGUAGE" = "java" ]; then
	CC_RUNTIME_LANGUAGE=java
	CC_SRC_PATH=/opt/gopath/src/github.com/chaincode/artwork/java
elif [ "$CC_SRC_LANGUAGE" = "javascript" ]; then
	CC_RUNTIME_LANGUAGE=node # chaincode runtime language is node.js
	CC_SRC_PATH=/opt/gopath/src/github.com/chaincode/artwork/javascript
elif [ "$CC_SRC_LANGUAGE" = "typescript" ]; then
	CC_RUNTIME_LANGUAGE=node # chaincode runtime language is node.js
	CC_SRC_PATH=/opt/gopath/src/github.com/chaincode/artwork/typescript
	echo Compiling TypeScript code into JavaScript ...
	pushd ../chaincode/artwork/typescript
	npm install
	npm run build
	popd
	echo Finished compiling TypeScript code into JavaScript
else
	echo The chaincode language ${CC_SRC_LANGUAGE} is not supported by this script
	echo Supported chaincode languages are: go, javascript, and typescript
	exit 1
fi


# clean the keystore
rm -rf ./hfc-key-store

# launch network; create channel and join peer to channel
cd ../first-network
echo y | ./byfn.sh down
echo y | ./byfn.sh up -a -n -s couchdb -o etcdraft

CONFIG_ROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer

COLLECTOR_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/collector.spchain.com/users/Admin@collector.spchain.com/msp
COLLECTOR_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/collector.spchain.com/peers/peer0.collector.spchain.com/tls/ca.crt

CREATOR_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/creator.spchain.com/users/Admin@creator.spchain.com/msp
CREATOR_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/creator.spchain.com/peers/peer0.creator.spchain.com/tls/ca.crt

MD_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/md.spchain.com/users/Admin@md.spchain.com/msp
MD_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/md.spchain.com/peers/peer0.md.spchain.com/tls/ca.crt

GALLERY_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/gallery.spchain.com/users/Admin@gallery.spchain.com/msp
GALLERY_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/gallery.spchain.com/peers/peer0.gallery.spchain.com/tls/ca.crt

ORDERER_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/ordererOrganizations/spchain.com/orderers/orderer.spchain.com/msp/tlscacerts/tlsca.spchain.com-cert.pem

set -x
echo "Installing smart contract on peer0.collector.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CollectorMSP \
  -e CORE_PEER_ADDRESS=peer0.collector.spchain.com:7051 \
  -e CORE_PEER_MSPCONFIGPATH=${COLLECTOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${COLLECTOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artwork \
    -v 1.0 \
    -p "$CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo -e "\nInstalling smart contract on peer0.creator.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CreatorMSP \
  -e CORE_PEER_ADDRESS=peer0.creator.spchain.com:9051 \
  -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${CREATOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artwork \
    -v 1.0 \
    -p "$CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling smart contract on peer0.md.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=MdMSP \
  -e CORE_PEER_ADDRESS=peer0.md.spchain.com:11051 \
  -e CORE_PEER_MSPCONFIGPATH=${MD_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${MD_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artwork \
    -v 1.0 \
    -p "$CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling smart contract on peer0.gallery.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_ADDRESS=peer0.gallery.spchain.com:13051 \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artwork \
    -v 1.0 \
    -p "$CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstantiating smart contract on mychannel"
docker exec \
  -e CORE_PEER_LOCALMSPID=CreatorMSP \
  -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
  cli \
  peer chaincode instantiate \
    -o orderer.spchain.com:7050 \
    -C mychannel \
    -n artwork \
    -l "$CC_RUNTIME_LANGUAGE" \
    -v 1.0 \
    -c '{"Args":[]}' \
    -P "OR ('CollectorMSP.member','CreatorMSP.member','MdMSP.member','GalleryMSP.member')" \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
    --peerAddresses peer0.creator.spchain.com:9051 \
    --tlsRootCertFiles ${CREATOR_TLS_ROOTCERT_FILE}

echo -e "\nWaiting for instantiation request to be committed ..."
sleep 10
echo ""

echo -e "\nSubmitting initLedger transaction to smart contract on mychannel"
echo "The transaction is sent to the peers with the chaincode installed so that chaincode is built before receiving the following requests"

#docker exec \
#  -e CORE_PEER_LOCALMSPID=CreatorMSP \
#  -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
#  cli \
#  peer chaincode invoke \
#    -o orderer.spchain.com:7050 \
#    -C mychannel \
#    -n artwork \
#    -c '{"Args":["uploadArtwork","tokenID","en-pointer","owner","creator"]}' \
#    --waitForEvent \
#    --tls \
#    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
#    --peerAddresses peer0.creator.spchain.com:9051 \
#    --tlsRootCertFiles ${CREATOR_TLS_ROOTCERT_FILE}
set +x

cat <<EOF

Total setup execution time : $(($(date +%s) - starttime)) secs ...

Next, use the FabCar applications to interact with the deployed FabCar contract.
The FabCar applications are available in multiple programming languages.
Follow the instructions for the programming language of your choice:

JavaScript:

  Start by changing into the "javascript" directory:
    cd javascript

  Next, install all required packages:
    npm install

  Then run the following applications to enroll the admin user, and register a new user
  called user1 which will be used by the other applications to interact with the deployed
  FabCar contract:
    node enrollAdmin
    node registerUser

  You can run the invoke application as follows. By default, the invoke application will
  create a new car, but you can update the application to submit other transactions:
    node invoke

  You can run the query application as follows. By default, the query application will
  return all cars, but you can update the application to evaluate other transactions:
    node query

EOF
