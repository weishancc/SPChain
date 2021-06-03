# -*- coding: utf-8 -*-
"""
Created on Mon May 17 17:33:24 2021

@author: Koma
"""
from collections import defaultdict
from matplotlib import pyplot as plt

def helper(filename, filterFunction, metricPosition):
    # Read experiment figure from 'result.txt'
    data = []
    
    with open(filename, 'r') as f:
        data = f.readlines()
    data = [x.strip('\n') for x in data]
    data = [x.split('\t') for x in data]
            
    # Delete unnecessary funcions and metrics (7=tps)
    data_dic = {k[0]: float(k[metricPosition]) for k in data}
    for k in filterFunction:
        del data_dic[k]
    
    return data_dic

# Concat differnet txDuration data
def concat_result(filterFunc, metricPosition):
    data_dic_10 = helper('txDuration_10.txt', filterFunc, metricPosition)
    data_dic_20 = helper('txDuration_20.txt', filterFunc, metricPosition)
    data_dic_30 = helper('txDuration_30.txt', filterFunc, metricPosition)
    data_dic_40 = helper('txDuration_40.txt', filterFunc, metricPosition)
    data_dic_50 = helper('txDuration_50.txt', filterFunc, metricPosition)   
    data_dic_60 = helper('txDuration_60.txt', filterFunc, metricPosition)   
    data_dic_70 = helper('txDuration_70.txt', filterFunc, metricPosition)   

    data_dic = defaultdict(list)
    for d in (data_dic_10, data_dic_20, data_dic_30, data_dic_40, data_dic_50, data_dic_60, data_dic_70): # you can list as many input dicts as you want here
        for key, value in d.items():
            data_dic[key].append(value)
    
    return data_dic

# Plot bar chart for different functions performance between differnet txDuration
def bar_plot(ax, data, colors=None, total_width=0.8, single_width=1):

    # Check if colors where provided, otherwhise use the default color cycle
    if colors is None:
        colors = plt.rcParams['axes.prop_cycle'].by_key()['color']
    
    n_bars = len(data)
    bar_width = total_width / n_bars
    
    bars = []

    for i, (name, values) in enumerate(data.items()):
    
        # The offset in x direction of that bar
        x_offset = (i - n_bars / 2) * bar_width + bar_width / 2

        # Draw a bar for every value of that type
        for x, y in enumerate(values):
            bar = ax.bar(x + x_offset, y, width=bar_width * single_width, color=colors[i % len(colors)])
        bars.append(bar[0])
        
    ax.legend(bars, data.keys(), loc=(1.04, 0.25))
    ax.grid(which='major', color='gray', linestyle='--', linewidth=.5)

# Plot bar chart for different functions latency between differnet txDuration
def plot_latency(ax, data):
    colors = plt.rcParams['axes.prop_cycle'].by_key()['color']
    markers=[">", "<", "o", "^", "^", "8", "D"]
    
    for i, v in enumerate(data):
        ax.plot(range(len(data[v])), data[v], colors[i % len(colors)], linestyle='dashed', 
                marker=markers[i % len(markers)], label=str(v), linewidth=2)
        
    ax.legend(data.keys(), loc=(1.04, 0.25))  
    ax.grid(which='major', color='gray', linestyle='--', linewidth=.5)
