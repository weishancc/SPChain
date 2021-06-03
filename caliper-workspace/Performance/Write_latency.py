# -*- coding: utf-8 -*-
"""
Created on Mon May 17 15:27:10 2021

@author: Koma
"""

# Plotted write functions include: 
#    - initialConsent
#    - grantRevokeConsent
#    - uploadArtwork
#    - transferArtwork
#    - addLog
#    - addModel
#    - addWallet
#    - transferBlalance


import matplotlib.pyplot as plt
import util

if __name__ == "__main__":
    filterFunc = ['readConsent', 'readArtwork', 'getHistoryForArtwork', 'readLog', 'readModel', 'readBalance', 'invokeModel']
    data = util.concat_result(filterFunc, 6)  # 6 for latency
 
    fig, ax = plt.subplots(figsize=(10, 6))
    util.plot_latency(ax, data)
    
    plt.xticks(range(7), ["10", "20", "30", "40", "50", "60", "70"])
    plt.xlabel('txDuration (sec)')
    plt.ylabel('Latency (sec)')
    plt.title('Write latency different functions under different transaction duration')
    plt.show()