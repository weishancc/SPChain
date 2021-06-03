# -*- coding: utf-8 -*-
"""
Created on Mon May 17 11:38:10 2021

@author: Koma
"""

# Plotted read functions include: 
#    - readConsent
#    - readArtwork
#    - getHistoryForArtwork
#    - readLog
#    - readModel
#    - readBalance
#    - invokeModel

import util
from matplotlib import pyplot as plt
 
if __name__ == "__main__":
    filterFunc = ['initialConsent', 'grantRevokeConsent', 'uploadArtwork', 'transferArtwork', 'addLog', 'addModel', 'addWallet', 'transferBlalance']
    data = util.concat_result(filterFunc, 7)  # 7 for tps
    
    fig, ax = plt.subplots(figsize=(10, 6))
    util.bar_plot(ax, data, total_width=.75, single_width=.9)
    
    plt.xticks(range(7), ["10", "20", "30", "40", "50", "60", "70"])
    plt.xlabel('txDuration (sec)')
    plt.ylabel('Throughtput (tps)')
    plt.title('Read performance of different functions under different transaction duration')
    plt.show()
    

    
    
    
    
    