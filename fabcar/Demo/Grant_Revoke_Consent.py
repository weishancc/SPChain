# -*- coding: utf-8 -*-
"""
Created on Tue May 11 08:45:50 2021

@author: Koma
"""

import requests
import os
from Initial_Phase import invoke_3A_CC

# ------------------
# Request the consent
# ------------------
def change_consent(api_3A, api_log, path, policy, res = 1, operation = "requestConsent"):
    # (1) invoke 3A_CC to update policy and then return sk_enc
    # --------------------------------------------------------
    TP_path = os.path.join(os.getcwd(), 'TP')
    with open(TP_path + '\\pk_TP.pem', 'rb') as f:   
        pk_TP = f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
        
    policy = policy.replace('pk_TP', pk_TP)
    ret = invoke_3A_CC(api_3A, path, policy)
    if ret == None: res = 0
    
    
    # (2) invoke log_CC   
    # ------------------
    
    # pk_DC
    DC_path = os.path.join(os.getcwd(), 'DC')
    with open(DC_path + '\\pk_DC.pem', 'rb') as f:   
        pk_DC = f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')   
 
    # sk_data
    with open(path + '\\sk_data.pem', 'rb') as f:
        sk_data =  f.read().decode('utf-8').replace('-----BEGIN RSA PRIVATE KEY-----\n', '').replace('\n-----END RSA PRIVATE KEY-----', '')
        pk_DS = sk_data
        
    # pk_DS (pk_enc)
    with open(path + '\\pk_enc.pem', 'rb') as f:
        pk_DS =  f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
        
    # Now we put it all together and start to request API
    par = {"pk_DS": pk_DS, "pk_DC": pk_DC, "pk_DP": pk_TP, "sk_data": sk_data, "status": str(res), "operation": operation}
    try:
        r = requests.post(api_log, json=par)
        r.raise_for_status()
        print(r.text)
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        raise SystemExit(e)
    
    # In the end, return sk_enc
    with open(path + '\\sk_enc.pem', 'rb') as f:
        sk_enc =  f.read().decode('utf-8').replace('-----BEGIN PUBLIC KEY-----\n', '').replace('\n-----END PUBLIC KEY-----', '')
    return sk_enc
            

if __name__== "__main__":  
    # Request Consent from DP, we temporarily don't write permissoin 
    # toward orbitdb here since orbitdb is still being completed, please refer to: 
    # https://github.com/orbitdb/orbit-db/blob/master/GUIDE.md#access-control
    name = input('Name of DS: ')
    api_3A = "http://140.123.105.138:8080/spchain/grantConsent/"
    api_log = "http://140.123.105.138:8080/spchain/addLog/"
    path = os.path.join(os.getcwd(), name)
    policy = "{\"R\":\"+pk_TP\", \"U\":\"+pk_TP\"}"  # example policy, can include "C", "R", "U", "D"
    sk_enc = change_consent(api_3A, api_log, path, policy)
    
    # Revoke Consent from DP, we don't remove write permission due to
    # same reason of orbitdb mentioned above
    policy = "{\"R\":\"-pk_TP\", \"U\":\"-pk_TP\"}"
    _ = change_consent(api_3A, api_log, path, policy)

    