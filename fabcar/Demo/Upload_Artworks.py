# -*- coding: utf-8 -*-
"""
Created on Thu May  6 13:18:47 2021

@author: Koma
"""
import os 
import requests

# ----------------------
# Invoke artowrk_CC API
# ----------------------
def invoke_artwork_CC(api, tokenID, multiHash, owner, creator):
    par = {"tokenID": tokenID, "multihash": multiHash, "owner": owner, "creator": creator}
    try:
        r = requests.post(api, json=par)
        r.raise_for_status()
        print(r.text)
        return 1
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        print(e)
        return 0
        
# ----------------------
# Invoke log_CC API
# ----------------------
def invoke_log_CC(api, path, res, operation):
    # pk_DC (for now is also pk_DP)
    DS_path = os.path.join(os.getcwd(), 'DC')
    with open(DS_path + '\\pk_DC.pem', 'rb') as f:   
        pk_DC = f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
            
    # sk_data
    with open(path + '\\sk_data.pem', 'rb') as f:
        sk_data =  f.read().decode('utf-8').replace('-----BEGIN RSA PRIVATE KEY-----\n', '').replace('\n-----END RSA PRIVATE KEY-----', '')
        pk_DS = sk_data
        
    # pk_DS (pk_enc)
    with open(path + '\\pk_enc.pem', 'rb') as f:
        pk_DS =  f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
        
    # Now we put it all together and start to request API
    par = {"pk_DS": pk_DS, "pk_DC": pk_DC, "pk_DP": pk_DC, "sk_data": sk_data, "status": str(res), "operation": operation}
    try:
        r = requests.post(api, json=par)
        r.raise_for_status()
        print(r.text)
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        raise SystemExit(e)
        

if __name__== "__main__":   
    hash = input('Input created orbitdb address: ')
    name = hash.split('/')[-1]
    path = os.path.join(os.getcwd(), name)

    # Get result from "setup.js upload" and then invoke artwork_CC
    tokenID = input('tokenID: ')
    multiHash = input('multi-hash: ')
    api = "http://140.123.105.138:8080/spchain/uploadArtwork/"
    res = invoke_artwork_CC(api, tokenID, multiHash, name, name)
    
    # Invoke log_CC
    api = "http://140.123.105.138:8080/spchain/addLog/"
    path = os.path.join(os.getcwd(), name)
    operation = "Upload the artwork"
    invoke_log_CC(api, path, res, operation)
    

#     p = subprocess.Popen(['/home/blockchain/.nvm/versions/node/v12.14.1/bin/node', '/home/blockchain/SPChain/orbit/setup.js', 'upload -a \"{\"imagePath\":'+ imagePath
#                           + ',\"name\":' + name
#                           + ',\"desc\":' + desc
#                           + ',\"price\":' + price
# 						  + '}\" ' 
#                           + '-d ' + 'testdb'], stdout=subprocess.PIPE)
    

  