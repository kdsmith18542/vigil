const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electronAPI', {
  startNode: () => ipcRenderer.send('start-node'),
  stopNode: () => ipcRenderer.send('stop-node'),
  openWallet: () => ipcRenderer.send('open-wallet'),
  installUpdate: () => ipcRenderer.send('install-update'),
  onNodeStatus: (callback) => ipcRenderer.on('node-status', (event, status) => callback(status)),
  onWalletStatus: (callback) => ipcRenderer.on('wallet-status', (event, status) => callback(status)),
  onInstallStatus: (callback) => ipcRenderer.on('install-status', (event, status) => callback(status))
});