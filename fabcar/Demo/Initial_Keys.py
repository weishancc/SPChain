# -*- coding: utf-8 -*-
"""
Created on Thu May  6 12:18:04 2021

@author: Koma
"""

import os
from Crypto.PublicKey import RSA
from Crypto import Random
from Crypto.Cipher import PKCS1_OAEP

class initialKeys():
    def __init__(self, name, path):
        self.name = name
        self.path = path
        
    # ----------------------------
    # Generate keys for end-users
    # ----------------------------
    def generate_keys(self):
        modulus_length = 1024
        random_generator = Random.new().read
        key = RSA.generate(modulus_length, random_generator)
        publicKey = key.publickey()
    
        return key, publicKey
    
    # ---------------------------------------------------------------
    # Export keys to the directory, where flag is used to distinguish 
    # different roles such as DS, DC and DP
    # ---------------------------------------------------------------
    def export_keys(self, key, publicKey, role):         
        if role == "DP":
            self.name = "DP"
            self.path = os.path.join(os.getcwd(), self.name)
        if role == "DC":
             self.name = "DC"
             self.path = os.path.join(os.getcwd(), self.name)
             
        if not os.path.exists(self.name):
            os.makedirs(self.name)      
        savePath = os.path.join(os.getcwd(), self.name)
            
        with open(savePath + '\\pk_' + role +'.pem', 'wb') as f:
            f.write(publicKey.exportKey(format='PEM'))
            
        with open(savePath + '\\sk_' + role +'.pem', 'wb') as f:
            f.write(key.exportKey(format='PEM'))
    
    # ------------------------------------------------
    # Use encryptor to encrypt the address of orbitdb
    # ------------------------------------------------
    def encrypt_address(self, publicKey, address):
        encryptor = PKCS1_OAEP.new(publicKey)
        
        with open(self.path + '\\enhash', 'wb') as f:
            enhash = encryptor.encrypt(str.encode(address))
            f.write(enhash)
        return enhash