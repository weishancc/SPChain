# -*- coding: utf-8 -*-
"""
Created on Mon May 10 08:48:30 2021

@author: Koma
"""
import os 
import json
import requests
from Upload_Artworks import invoke_log_CC

# -------------------------------------------------------------------
# Get owner and creator of the artwork through calling artwork_CC API
# -------------------------------------------------------------------
def get_roles(api, tokenID):
    par = {"tokenID": tokenID}
    
    try:
        r = requests.get(api, json=par)
        r.raise_for_status()
        res = json.loads(r.text)  #{'response': '{"Creator": ...}'}    
        res = json.loads(res['response'])  #{'Creator': 'Koma', 'Owner': ...}
        
        return res
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        print(e)
        return 0
        
# -------------------------------------------------------------------
# Get owner and creator of the artwork through calling artwork_CC API
# -------------------------------------------------------------------
def change_balance(api, newCollector, collector, creator, price, r_y):
    par = {"newCollector": newCollector, "collector": collector, "creator": creator, "price": price, "r_y": r_y}
    
    try:
        r = requests.post(api, json=par)
        r.raise_for_status()
        print(r.text)
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        print(e)    

# ----------------------
# Invoke artowrk_CC API
# ----------------------
def invoke_artwork_CC(api, tokenID, newCollector, multiHash):
    par = {"tokenID": tokenID, "newCollector": newCollector, "multihash": multiHash}
    try:
        r = requests.post(api, json=par)
        r.raise_for_status()
        print(r.text)
        return 1
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        print(e)
        return 0

        
if __name__== "__main__":  
    # (1) Invoke artwork_CC
    # ----------------------
    tokenID = input('Input tokenID: ')
    api = "http://140.123.105.138:8080/spchain/readArtwork/"
    res = get_roles(api, tokenID)
    
    
    # (2) Invoke wallet_CC
    # ---------------------
    creator, owner = res['Creator'], res['Owner']
    r_y = str(input('Input royalty rate: '))
    newCollector = input('Input new collector: ')
    price = str(input('Input price: '))  # We simplify query the price from orbitdb here 
    
    api = "http://140.123.105.138:8080/spchain/transferBalance/"
    change_balance(api, newCollector, owner, creator, price, r_y)
    
    
    # (3) Update orbitdb information
    # ----------------------------------------------------------------
    # -> node setup.js transfer -t $tokenID -n $newCollector -o $owner
    
    
    # (4) Invoke artwork_CC
    # ---------------------
    api = "http://140.123.105.138:8080/spchain/transferArtwork/"
    multiHash = input('multi-hash: ')
    res = invoke_artwork_CC(api, tokenID, newCollector, multiHash)


    # (5) Invoke log_CC
    # -----------------
    api = "http://140.123.105.138:8080/spchain/addLog/"
    path = os.path.join(os.getcwd(), newCollector)
    operation = "Transfer the artwork"
    invoke_log_CC(api, path, res, operation)
    
    
    
    
    
    