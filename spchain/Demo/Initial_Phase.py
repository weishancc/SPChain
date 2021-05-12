# -*- coding: utf-8 -*-
"""
Created on Wed May  5 13:43:40 2021

@author: Koma
"""

import requests
import os
from Initial_Keys import initialKeys

# -----------------
# Invoke 3A_CC API
# -----------------
def invoke_3A_CC(api, path, policy):
    # pk_DC   
    DS_path = os.path.join(os.getcwd(), 'DC')
    with open(DS_path + '\\pk_DC.pem', 'rb') as f:   
        pk_DC = f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
        
    # ehash
    with open(path + '\\enhash', 'rb') as f:
        enhash = str(f.read())
    
    # pk_DS (pk_enc)
    with open(path + '\\pk_enc.pem', 'rb') as f:
        pk_enc =  f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
        pk_DS = pk_enc
    
    # Now we put it all together and start to request API
    par = {"pk_DS": pk_DS, "pk_DC": pk_DC, "policy": policy, "enhash": enhash, "pk_enc": pk_enc}
    try:
        r = requests.post(api, json=par)
        r.raise_for_status()
        return r.text
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        raise SystemExit(e)
    
        
if __name__== "__main__":   
    hash = input('Input created orbitdb address: ')
    name = hash.split('/')[-1]
    path = os.path.join(os.getcwd(), name)
    
    #Generate key-pairs of data and enc
    initial = initialKeys(name, path)
    sk_enc, pk_enc = initial.generate_keys()
    sk_data, pk_data = initial.generate_keys()
    sk_DC, pk_DC = initial.generate_keys()
    sk_TP, pk_TP = initial.generate_keys()
    
    initial.export_keys(sk_enc, pk_enc, 'enc') 
    initial.export_keys(sk_data, pk_data, 'data' ) 
    enhash = initial.encrypt_address(pk_enc ,hash)
    initial.export_keys(sk_DC, pk_DC, 'DC' ) 
    initial.export_keys(sk_TP, pk_TP, 'TP') 
    
    # Next we inoke 3A_CC, for "DC" we don't need 
    # to set policy because in 3A_CC set all by default
    api = "http://140.123.105.138:8080/spchain/grantConsent/"
    policy = ""
    ret = invoke_3A_CC(api, path, policy)   