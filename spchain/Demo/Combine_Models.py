# -*- coding: utf-8 -*-
"""
Created on Mon May 10 12:45:03 2021

@author: Koma
"""

import requests
import json

# -------------------------------
# Model developers upload models
# -------------------------------
def upload_models(api):
    # Inpupt information
    name = input('Name of the model: ')
    url = input('Address of model in Docker Registry: ')
    md = input('Model developer: ')
    desc = input('Description of the model: ')
    
    # Invoke model_CC API
    par = {"name": name, "address": url, "creator": md, "desc": desc}
    
    try:
        r = requests.post(api, json=par)
        r.raise_for_status()
        print(r.text)
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        print(e)

# ---------------
# Read the model
# ---------------
def reading_models(api, name):
    par = {"name": name}
    
    try:
        r = requests.get(api, json=par)
        r.raise_for_status()
        print(r.text)
        res = json.loads(r.text)
        res = json.loads(res['response'])

        return res
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        print(e)
        
# -------------------------
# Creators invoke the model
# -------------------------
def invoking_models(api, artwork):
    # Query the address of the selected model first
    par = {"api": api, "artwork": artwork}
    
    try:
        r = requests.post(api, json=par)
        r.raise_for_status()
        print(r.text)
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        print(e)
     
# ----------------------
# Invoke wallet_CC API
# ----------------------
def invoke_wallet_CC(api, model_developer, premium):
    par = {"creator": model_developer, "premium": premium}
    try:
        r = requests.post(api, json=par)
        r.raise_for_status()
        print(r.text)
    except requests.exceptions.RequestException as e:  
        # Whoops 200 isn't return
        print(e)


# -----------------------------------------------------
# With regard to combingin models there are two parts: 
# (1) Upload the model (2) Invoking the modle
# -----------------------------------------------------
if __name__== "__main__":  
    # (1) Upload the model
    api = "http://140.123.105.138:8080/spchain/addModel/"
    upload_models(api)
    
    # (2-1) Read address of selected model first
    api = "http://140.123.105.138:8080/spchain/readModel/"
    name = input('Name of the model: ')
    res = reading_models(api, name)
    url = res['Address']
    
    # (2-2) Running application container, and this step provides "service api"
    
    # (2-3) Invoking the model
    service_api = "http://140.123.105.138:5000/mosaic"
    artwork = input('Path of the artowrk: ')
    artifact = invoking_models(service_api, artwork)
    
    # (2-4) Invoke wallet_CC
    api = "http://140.123.105.138:8080/spchain/addBalance/"
    md = res['Creator']
    premium = str(input('Premium: '))
    invoke_wallet_CC(api, md, premium)
    
    
    