# -*- coding: utf-8 -*-
"""
Created on Tue May 11 12:17:40 2021

@author: Koma
"""

import os 
import requests
import json
import ast
from Crypto.PublicKey import RSA
from Crypto.Cipher import PKCS1_OAEP


# ---------------------
# Read policy (Consent)
# ---------------------
def read_policy(api, path, pk_TP):
    # pk_DC   
    DS_path = os.path.join(os.getcwd(), 'DC')
    with open(DS_path + '\\pk_DC.pem', 'rb') as f:   
        pk_DC = f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
            
    # pk_DS (pk_enc)
    with open(path + '\\pk_enc.pem', 'rb') as f:
        pk_enc =  f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
        pk_DS = pk_enc
    
    # Now we put it all together and start to request API
    par = {"pk_DS": pk_DS, "pk_DC": pk_DC}
    try:
        r = requests.get(api, json=par)
        r.raise_for_status()
        res = json.loads(r.text)  #{'response': '{"Enhash": ...}'}
        res = json.loads(res['response'])  #{'Enhash': 'xxx', 'Policy': ...}
        #print('Policy: ', res['Policy'])
        
        # Check if TP has corresponding policy
        policy_result = check_policy(res['Policy'], pk_TP)
        if policy_result != None:
            return res['Enhash'], pk_DC, pk_DS
        return None, pk_DC, pk_DS
       
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        raise SystemExit(e)

# ------------------------------------
# Check if TP in the policy dictionary
# ------------------------------------
def check_policy(dic, pk_TP):
    return [v for v in dic if pk_TP in dic[v]]
    
# --------------------
# Log and get sk_data
# --------------------
def read_sk_data(api, pk_DS, pk_DC, pk_TP):
    # Now we put it all together and start to request API
    par = {"pk_DS": pk_DS, "pk_DC": pk_DC, "pk_DP": pk_TP}
    try:
        r = requests.get(api, json=par)
        r.raise_for_status()
        res = json.loads(r.text)  #{'response': '{"SKData": ...}'}
        res = json.loads(res['response'])  #{'SKData': 'xxx', 'Operation': ...}
        
        return res['SKData']
    
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        raise SystemExit(e)
    
    
if __name__== "__main__":  
    # (1) Invoke 3A_CC to check policy (consent) and get enhash
    # ---------------------------------------------------------
    name = input('Name of DS: ')
    path = os.path.join(os.getcwd(), name)
    api = "http://140.123.105.138:8080/spchain/readConsent/"   
    with open(os.path.join(os.getcwd(), 'TP') + '\\pk_TP.pem', 'rb') as f:   
        pk_TP = f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
        
    enhash, pk_DC, pk_DS = read_policy(api, path, pk_TP)
    
    
    # (2) Invoke log_CC to log and get sk_data, in the end,
    #     TP got -> sk_enc, sk_data, enhash 
    # -----------------------------------------------------
    with open(path + '\\sk_enc.pem', 'rb') as f:
        sk_enc =  RSA.importKey(f.read().decode('utf-8'))
    api = "http://140.123.105.138:8080/spchain/readLog/"
    sk_data = read_sk_data(api, pk_DS, pk_DC, pk_TP)
    
    # Use decryptor to decrypt enhash -> get original orbitdb address (i.e., decrypted)
    decryptor = PKCS1_OAEP.new(sk_enc)
    decrypted = decryptor.decrypt(ast.literal_eval(str(enhash)))
    print('Decrupted Address: ', decrypted)
    
    
    # (3) Now we have decrypted address, query this address to get data,
    #     and then decrypted data with sk_data -> get final ciphertext
    # -----------------------------------------------------------------
    # node request.js query -a /orbitdb/...
    # node request.js decrypt -s sk_data.pem -a /orbitdb/...
    

    
    