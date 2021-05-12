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
	ARTWORK_CC_SRC_PATH=github.com/chaincode/artworks
	WALLET_CC_SRC_PATH=github.com/chaincode/wallets
	MODEL_CC_SRC_PATH=github.com/chaincode/models
	THREE_CC_SRC_PATH=github.com/chaincode/3A
	LOG_CC_SRC_PATH=github.com/chaincode/logs
elif [ "$CC_SRC_LANGUAGE" = "java" ]; then
	CC_RUNTIME_LANGUAGE=java
	CC_SRC_PATH=/opt/gopath/src/github.com/chaincode/logs/java
elif [ "$CC_SRC_LANGUAGE" = "javascript" ]; then
	CC_RUNTIME_LANGUAGE=node # chaincode runtime language is node.js
	CC_SRC_PATH=/opt/gopath/src/github.com/chaincode/logs/javascript
elif [ "$CC_SRC_LANGUAGE" = "typescript" ]; then
	CC_RUNTIME_LANGUAGE=node # chaincode runtime language is node.js
	CC_SRC_PATH=/opt/gopath/src/github.com/chaincode/logs/typescript
	echo Compiling TypeScript code into JavaScript ...
	pushd ../chaincode/logs/typescript
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
echo y | ./byfn.sh up -a -n -s couchdb -o etcdraft -i 1.4.4

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

#############################################
# Begin to install and instantiate chaincodes
##############################################

#--------------------------------------------------------------#
echo "Begin to installing artwork_CC"
set -x
echo "\nInstalling artwork_CC on peer0.collector.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CollectorMSP \
  -e CORE_PEER_ADDRESS=peer0.collector.spchain.com:7051 \
  -e CORE_PEER_MSPCONFIGPATH=${COLLECTOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${COLLECTOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artworks \
    -v 1.0 \
    -p "$ARTWORK_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo "\nInstalling artwork_CC on peer1.collector.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CollectorMSP \
  -e CORE_PEER_ADDRESS=peer1.collector.spchain.com:8051 \
  -e CORE_PEER_MSPCONFIGPATH=${COLLECTOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${COLLECTOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artworks \
    -v 1.0 \
    -p "$ARTWORK_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo "\nInstalling artwork_CC on peer0.creator.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CreatorMSP \
  -e CORE_PEER_ADDRESS=peer0.creator.spchain.com:9051 \
  -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${CREATOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artworks \
    -v 1.0 \
    -p "$ARTWORK_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo "\nInstalling artwork_CC on peer1.creator.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CreatorMSP \
  -e CORE_PEER_ADDRESS=peer1.creator.spchain.com:10051 \
  -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${CREATOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artworks \
    -v 1.0 \
    -p "$ARTWORK_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo -e "\nInstalling artwork_CC on peer0.gallery.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_ADDRESS=peer0.gallery.spchain.com:13051 \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artworks \
    -v 1.0 \
    -p "$ARTWORK_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling artwork_CC on peer1.gallery.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_ADDRESS=peer1.gallery.spchain.com:14051 \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n artworks \
    -v 1.0 \
    -p "$ARTWORK_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstantiating artwork_CC on mychannel"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  cli \
  peer chaincode instantiate \
    -o orderer.spchain.com:7050 \
    -C mychannel \
    -n artworks \
    -l "$CC_RUNTIME_LANGUAGE" \
    -v 1.0 \
    -c '{"Args":[]}' \
    -P "OR ('CollectorMSP.member','CreatorMSP.member','MdMSP.member','GalleryMSP.member')" \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
    --peerAddresses peer0.gallery.spchain.com:13051 \
    --tlsRootCertFiles ${GALLERY_TLS_ROOTCERT_FILE}

echo -e "\nWaiting for instantiation request to be committed ..."
sleep 10
echo ""

echo -e "\nSubmitting initLedger transaction to smart contract on mychannel"
echo "The transaction is sent to the peers with the chaincode installed so that chaincode is built before receiving the following requests"
echo "artwork_CC instantiated !!"

#--------------------------------------------------------------#
echo "Begin to installing wallet_CC"
set -x
echo "\nInstalling wallet_CC on peer0.collector.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=CollectorMSP \
 -e CORE_PEER_ADDRESS=peer0.collector.spchain.com:7051 \
 -e CORE_PEER_MSPCONFIGPATH=${COLLECTOR_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${COLLECTOR_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n wallets \
   -v 1.0 \
   -p "$WALLET_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"

echo "\nInstalling wallet_CC on peer1.collector.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=CollectorMSP \
 -e CORE_PEER_ADDRESS=peer1.collector.spchain.com:8051 \
 -e CORE_PEER_MSPCONFIGPATH=${COLLECTOR_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${COLLECTOR_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n wallets \
   -v 1.0 \
   -p "$WALLET_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"

echo "\nInstalling wallet_CC on peer0.creator.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=CreatorMSP \
 -e CORE_PEER_ADDRESS=peer0.creator.spchain.com:9051 \
 -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${CREATOR_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n wallets \
   -v 1.0 \
   -p "$WALLET_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo "\nInstalling wallet_CC on peer1.creator.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=CreatorMSP \
 -e CORE_PEER_ADDRESS=peer1.creator.spchain.com:10051 \
 -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${CREATOR_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n wallets \
   -v 1.0 \
   -p "$WALLET_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"

echo -e "\nInstalling wallet_CC on peer0.md.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=MdMSP \
 -e CORE_PEER_ADDRESS=peer0.md.spchain.com:11051 \
 -e CORE_PEER_MSPCONFIGPATH=${MD_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${MD_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n wallets \
   -v 1.0 \
   -p "$WALLET_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling wallet_CC on peer1.md.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=MdMSP \
 -e CORE_PEER_ADDRESS=peer1.md.spchain.com:12051 \
 -e CORE_PEER_MSPCONFIGPATH=${MD_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${MD_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n wallets \
   -v 1.0 \
   -p "$WALLET_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling wallet_CC on peer0.gallery.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=GalleryMSP \
 -e CORE_PEER_ADDRESS=peer0.gallery.spchain.com:13051 \
 -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n wallets \
   -v 1.0 \
   -p "$WALLET_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling wallet_CC on peer1.gallery.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=GalleryMSP \
 -e CORE_PEER_ADDRESS=peer1.gallery.spchain.com:14051 \
 -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n wallets \
   -v 1.0 \
   -p "$WALLET_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstantiating wallet_CC on mychannel"
docker exec \
 -e CORE_PEER_LOCALMSPID=GalleryMSP \
 -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
 cli \
 peer chaincode instantiate \
   -o orderer.spchain.com:7050 \
   -C mychannel \
   -n wallets \
   -l "$CC_RUNTIME_LANGUAGE" \
   -v 1.0 \
   -c '{"Args":[]}' \
   -P "OR ('CollectorMSP.member','CreatorMSP.member','MdMSP.member','GalleryMSP.member')" \
   --tls \
   --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
   --peerAddresses peer0.gallery.spchain.com:13051 \
   --tlsRootCertFiles ${GALLERY_TLS_ROOTCERT_FILE}

echo -e "\nWaiting for instantiation request to be committed ..."
sleep 10
echo ""

echo -e "\nSubmitting initLedger transaction to smart contract on mychannel"
echo "The transaction is sent to the peers with the chaincode installed so that chaincode is built before receiving the following requests"
echo "wallet_CC instantiated !!"

#--------------------------------------------------------------#
echo "Begin to installing model_CC"
set -x
echo -e "\nInstalling model_CC on peer0.gallery.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_ADDRESS=peer0.gallery.spchain.com:13051 \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n models \
    -v 1.0 \
    -p "$MODEL_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling model_CC on peer1.gallery.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_ADDRESS=peer1.gallery.spchain.com:14051 \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n models \
    -v 1.0 \
    -p "$MODEL_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling model_CC on peer0.md.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=MdMSP \
 -e CORE_PEER_ADDRESS=peer0.md.spchain.com:11051 \
 -e CORE_PEER_MSPCONFIGPATH=${MD_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${MD_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n models \
   -v 1.0 \
   -p "$MODEL_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling model_CC on peer1.md.spchain.com"
docker exec \
 -e CORE_PEER_LOCALMSPID=MdMSP \
 -e CORE_PEER_ADDRESS=peer1.md.spchain.com:12051 \
 -e CORE_PEER_MSPCONFIGPATH=${MD_MSPCONFIGPATH} \
 -e CORE_PEER_TLS_ROOTCERT_FILE=${MD_TLS_ROOTCERT_FILE} \
 cli \
 peer chaincode install \
   -n models \
   -v 1.0 \
   -p "$MODEL_CC_SRC_PATH" \
   -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "Instantiating model_CC on mychannel"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  cli \
  peer chaincode instantiate \
    -o orderer.spchain.com:7050 \
    -C mychannel \
    -n models \
    -l "$CC_RUNTIME_LANGUAGE" \
    -v 1.0 \
    -c '{"Args":[]}' \
    -P "OR ('CollectorMSP.member','CreatorMSP.member','MdMSP.member','GalleryMSP.member')" \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
    --peerAddresses peer0.gallery.spchain.com:13051 \
    --tlsRootCertFiles ${GALLERY_TLS_ROOTCERT_FILE}

echo -e "\nWaiting for instantiation request to be committed ..."
sleep 10
echo ""

echo -e "\nSubmitting initLedger transaction to smart contract on mychannel"
echo "The transaction is sent to the peers with the chaincode installed so that chaincode is built before receiving the following requests"
echo "model_CC instantiated !!"

#--------------------------------------------------------------#
echo "Begin to installing 3A_CC"
set -x
echo "\nInstalling 3A_CC on peer0.collector.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CollectorMSP \
  -e CORE_PEER_ADDRESS=peer0.collector.spchain.com:7051 \
  -e CORE_PEER_MSPCONFIGPATH=${COLLECTOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${COLLECTOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n 3A \
    -v 1.0 \
    -p "$THREE_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo "\nInstalling 3A_CC on peer1.collector.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CollectorMSP \
  -e CORE_PEER_ADDRESS=peer1.collector.spchain.com:8051 \
  -e CORE_PEER_MSPCONFIGPATH=${COLLECTOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${COLLECTOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n 3A \
    -v 1.0 \
    -p "$THREE_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo "\nInstalling 3A_CC on peer0.creator.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CreatorMSP \
  -e CORE_PEER_ADDRESS=peer0.creator.spchain.com:9051 \
  -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${CREATOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n 3A \
    -v 1.0 \
    -p "$THREE_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo "\nInstalling 3A_CC on peer1.creator.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CreatorMSP \
  -e CORE_PEER_ADDRESS=peer1.creator.spchain.com:10051 \
  -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${CREATOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n 3A \
    -v 1.0 \
    -p "$THREE_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo -e "\nInstalling 3A_CC on peer0.gallery.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_ADDRESS=peer0.gallery.spchain.com:13051 \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n 3A \
    -v 1.0 \
    -p "$THREE_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling 3A_CC on peer1.gallery.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_ADDRESS=peer1.gallery.spchain.com:14051 \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n 3A \
    -v 1.0 \
    -p "$THREE_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "Instantiating 3A_CC on mychannel"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  cli \
  peer chaincode instantiate \
    -o orderer.spchain.com:7050 \
    -C mychannel \
    -n 3A \
    -l "$CC_RUNTIME_LANGUAGE" \
    -v 1.0 \
    -c '{"Args":[]}' \
    -P "OR ('CollectorMSP.member','CreatorMSP.member','MdMSP.member','GalleryMSP.member')" \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
    --peerAddresses peer0.gallery.spchain.com:13051 \
    --tlsRootCertFiles ${GALLERY_TLS_ROOTCERT_FILE}

echo -e "\nWaiting for instantiation request to be committed ..."
sleep 10
echo ""

echo -e "\nSubmitting initLedger transaction to smart contract on mychannel"
echo "The transaction is sent to the peers with the chaincode installed so that chaincode is built before receiving the following requests"
echo "3A_CC instantiated !!"

#--------------------------------------------------------------#
echo "Begin to installing log_CC"
set -x
echo "\nInstalling log_CC on peer0.collector.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CollectorMSP \
  -e CORE_PEER_ADDRESS=peer0.collector.spchain.com:7051 \
  -e CORE_PEER_MSPCONFIGPATH=${COLLECTOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${COLLECTOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n logs \
    -v 1.0 \
    -p "$LOG_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo "\nInstalling log_CC on peer1.collector.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CollectorMSP \
  -e CORE_PEER_ADDRESS=peer1.collector.spchain.com:8051 \
  -e CORE_PEER_MSPCONFIGPATH=${COLLECTOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${COLLECTOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n logs \
    -v 1.0 \
    -p "$LOG_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo "\nInstalling log_CC on peer0.creator.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CreatorMSP \
  -e CORE_PEER_ADDRESS=peer0.creator.spchain.com:9051 \
  -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${CREATOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n logs \
    -v 1.0 \
    -p "$LOG_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo "\nInstalling log_CC on peer1.creator.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=CreatorMSP \
  -e CORE_PEER_ADDRESS=peer1.creator.spchain.com:10051 \
  -e CORE_PEER_MSPCONFIGPATH=${CREATOR_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${CREATOR_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n logs \
    -v 1.0 \
    -p "$LOG_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"

echo -e "\nInstalling log_CC on peer0.gallery.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_ADDRESS=peer0.gallery.spchain.com:13051 \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n logs \
    -v 1.0 \
    -p "$LOG_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstalling log_CC on peer1.gallery.spchain.com"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_ADDRESS=peer1.gallery.spchain.com:14051 \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${GALLERY_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n logs \
    -v 1.0 \
    -p "$LOG_CC_SRC_PATH" \
    -l "$CC_RUNTIME_LANGUAGE"
echo ""

echo -e "\nInstantiating log_CC on mychannel"
docker exec \
  -e CORE_PEER_LOCALMSPID=GalleryMSP \
  -e CORE_PEER_MSPCONFIGPATH=${GALLERY_MSPCONFIGPATH} \
  cli \
  peer chaincode instantiate \
    -o orderer.spchain.com:7050 \
    -C mychannel \
    -n logs \
    -l "$CC_RUNTIME_LANGUAGE" \
    -v 1.0 \
    -c '{"Args":[]}' \
    -P "OR ('CollectorMSP.member','CreatorMSP.member','MdMSP.member','GalleryMSP.member')" \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
    --peerAddresses peer0.gallery.spchain.com:13051 \
    --tlsRootCertFiles ${GALLERY_TLS_ROOTCERT_FILE}

echo -e "\nWaiting for instantiation request to be committed ..."
sleep 10
echo ""

echo -e "\nSubmitting initLedger transaction to smart contract on mychannel"
echo "The transaction is sent to the peers with the chaincode installed so that chaincode is built before receiving the following requests"
echo "log_CC instantiated !!"

set +x

cat <<EOF

===============================================================
Good! All Chaincodes Installed  and Instantiated Successfully!!
===============================================================

EOF
