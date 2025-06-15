const { app, BrowserWindow, ipcMain } = require('electron');
const path = require('path');
const { exec } = require('child_process');

function createWindow () {
  const win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js')
    }
  });

  win.loadFile('index.html');

  // IPC Handlers
  let dcrdProcess = null;
  let dcrwalletProcess = null;

  const dcrdPath = path.join(__dirname, '..', 'node', 'vgld.exe'); // Adjust path as needed
  const dcrwalletPath = path.join(__dirname, '..', 'wallet', 'vigilwallet.exe'); // Adjust path as needed

  ipcMain.on('start-node', (event) => {
    if (dcrdProcess && !dcrdProcess.killed) {
      event.sender.send('node-status', 'Node Status: Already Running');
      return;
    }
    event.sender.send('node-status', 'Node Status: Starting...');
    dcrdProcess = exec(dcrdPath, { cwd: path.join(__dirname, '..', 'node') }, (error, stdout, stderr) => {
      if (error) {
        console.error(`exec error: ${error}`);
        event.sender.send('node-status', `Node Status: Error - ${error.message}`);
        return;
      }
      console.log(`stdout: ${stdout}`);
      console.error(`stderr: ${stderr}`);
    });

    dcrdProcess.stdout.on('data', (data) => {
      console.log(`dcrd stdout: ${data}`);
      if (data.includes('RPC server listening on')) {
        event.sender.send('node-status', 'Node Status: Running');
      }
    });

    dcrdProcess.stderr.on('data', (data) => {
      console.error(`dcrd stderr: ${data}`);
    });

    dcrdProcess.on('close', (code) => {
      console.log(`dcrd process exited with code ${code}`);
      event.sender.send('node-status', `Node Status: Exited with code ${code}`);
    });

    dcrdProcess.on('error', (err) => {
      console.error(`Failed to start dcrd process: ${err.message}`);
      event.sender.send('node-status', `Node Status: Error - ${err.message}`);
    });
  });

  ipcMain.on('stop-node', (event) => {
    if (dcrdProcess) {
      event.sender.send('node-status', 'Node Status: Stopping...');
      dcrdProcess.kill();
      dcrdProcess = null;
    } else {
      event.sender.send('node-status', 'Node Status: Not Running');
    }
  });

  ipcMain.on('open-wallet', (event) => {
    if (dcrwalletProcess && !dcrwalletProcess.killed) {
      event.sender.send('wallet-status', 'Wallet Status: Already Open');
      return;
    }
    event.sender.send('wallet-status', 'Wallet Status: Opening...');
    dcrwalletProcess = exec(dcrwalletPath, { cwd: path.join(__dirname, '..', 'wallet') }, (error, stdout, stderr) => {
      if (error) {
        console.error(`exec error: ${error}`);
        event.sender.send('wallet-status', `Wallet Status: Error - ${error.message}`);
        return;
      }
      console.log(`stdout: ${stdout}`);
      console.error(`stderr: ${stderr}`);
    });

    dcrwalletProcess.stdout.on('data', (data) => {
      console.log(`dcrwallet stdout: ${data}`);
      if (data.includes('Wallet is unlocked')) {
        event.sender.send('wallet-status', 'Wallet Status: Opened');
      }
    });

    dcrwalletProcess.stderr.on('data', (data) => {
      console.error(`dcrwallet stderr: ${data}`);
    });

    dcrwalletProcess.on('close', (code) => {
      console.log(`dcrwallet process exited with code ${code}`);
      event.sender.send('wallet-status', `Wallet Status: Exited with code ${code}`);
    });

    dcrwalletProcess.on('error', (err) => {
      console.error(`Failed to start dcrwallet process: ${err.message}`);
      event.sender.send('wallet-status', `Wallet Status: Error - ${err.message}`);
    });
  });

  ipcMain.on('install-update', (event) => {
    event.sender.send('install-status', 'Installation Status: Not yet implemented. Please update manually.');
  });
}

app.whenReady().then(() => {
  createWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});