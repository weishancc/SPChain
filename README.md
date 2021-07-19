# SPChain Implementation
<img src="https://github.com/weishancc/SPChain/blob/master/High-level%20System%20Architecture.png" width="657" length="427"> 

* [Requirements](#requirements)
* [Start up blockchain network and install chaincodes](#start-up-blockchain-network-and-install-chaincodes)
  * [Local Network](#local-network)
  * [API Server](#api-server)
* [Initial Phase](#initial-phase)
* [Uploading Artworks](#uploading-artworks)
* [Transferring Artworks](#transferring-artworks)
* [Combining Models](#combining-models)
* [Grant and Revoke Consent](#grant-and-revoke-consent)
* [Request Data](#request-data)


## [Requirements](#requirements)
* [Hyperledger Fabric Images](https://hub.docker.com/search?q=hyperledger%2Ffabric&type=image) (amd-v1.4.4)
* [Go](https://golang.org/) (v1.11.0 or later)
* [Docker](https://docs.docker.com/get-docker/) (v18.03 or later)
* [Docker Compose](https://docs.docker.com/compose/) (v1.14.0 or later)
* [orbit-db](https://github.com/orbitdb/orbit-db) (Put in folder "orbit/")


## [Start up blockchain network and install chaincodes](#start-up-blockchain-network-and-install-chaincodes)

### [Local Network](#local-network)
Start up blockchain network locally:
```
cd spchain
./startFabric.sh
```
The script will start up the network, install, and instantiate all chaincodes, please wait for a moment until you see the successful message &#61;&#61;&#61;&#61; __Good! All Chaincodes Installed  and Instantiated Successfully!! &#61;&#61;&#61;&#61;__

### [API Server](#api-server)
Now, let's wrap the blockchain network by nodejs, in that way we can interact with the network via RESTful API:
```
cd javascript
npm install
```

Enroll Admin and Regist User:
```
node enrollGalleryAdmin.js
node node registerGalleryUser.js
```

Start up API server:
```
node apiserver.js
```

---
*### Running following instructions under path [SPChain/spchain/Demo/](https://github.com/weishancc/SPChain/tree/master/spchain/Demo) ###*  
*### Remember to change "api" in all python files to your own api (change ip address) ###*
## [Initial Phase](#initial-phase)
This step generates the keys of Service Provider, Third-Party, and Data Subject, then invokes 3A_CC  
```
python Initial_Phase.py  
```
After that, you will see folders containing keys under path "SPChain/spchain/Demo/" and see the payload in the blockchain server

## [Uploading Artworks](#uploading-artworks)
This step demonstrates flow of uploading artowrks
```
python Upload_Artworks.py
```
then follow the instructions described in Upload_Artworks.py:  
(1) node setup.js upload -a..  
(2) Get result from "setup.js upload" and then invoke artwork_CC  
(3) Invoke log_CC  

After that, artwork_CC and log_CC will be invoked and showd the payload in the blockchain server

## [Transferring Artworks](#transferring-artworks)
This step demonstrates flow of transferring artowrks
```
python Transfer_Artworks.py
```
then follow the instructions described in Transfer_Artworks.py:  
(1) Invoke artwork_CC  
(2) Invoke wallet_CC  
(3) node setup.js transfer..  
(4) Invoke artwork_CC  
(5) Invoke log_CC  

## [Combining Models](#combining-models)
This step demonstrates flow of uploading models and invoking models
```
python Combine_Models.py
```
then follow the instructions described in Combine_Models.py:  
(1) Upload the model  
(2) Read address of selected model first  
(3) Running application container, and this step provides "service api"  
(4) Invoking the model  
(5) Invoke wallet_CC  

## [Grant and Revoke Consent](#grant-and-revoke-consent)
This step demonstrates flow of granting consents and revoking consents from (by) Data Subject
```
python Grant_Revoke_Consent.py
```
then follow the instructions described in Grant_Revoke_Consent.py:  
(1) Request Consent from DP  
(2) Revoke Consent from DP  

## [Request Data](#request-data)
This step demonstrates flow of requesting data from Third-Party
```
python Request_Data.py
```
then follow the instructions described in Request_Data.py:  
(1) Invoke 3A_CC to check policy (consent) and get enhash  
(2) Invoke log_CC to log and get sk_data  
(3) Use decryptor to decrypt enhash -> get original orbitdb address (i.e., decrypted)  
(4) Now we have decrypted address, query this address to get data, and then decrypted data with sk_data -> get final ciphertext  
