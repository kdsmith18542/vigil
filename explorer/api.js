// api.js

// Function to fetch data from a given API endpoint
async function fetchDataFromAPI(endpoint) {
    try {
        const response = await fetch(endpoint);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error(`Error fetching data from ${endpoint}:`, error);
        return null;
    }
}

// Example: Function to fetch mining info
// async function fetchMiningInfo() {
//     return fetchDataFromAPI('/api/getmininginfo');
// }

// Example: Function to fetch ticket pool value
// async function fetchTicketPoolValue() {
//     return fetchDataFromAPI('/api/getticketpoolvalue');
// }

// Function to fetch staking info
async function getStakingInfo() {
    const stakingInfo = await fetchDataFromAPI('/api/getstakinginfo');
    if (stakingInfo) {
        // Assuming the backend returns data like:
        // { TotalStaked: ..., SecurityScore: ..., ProjectedROI: ... }
        return {
            totalVGLStaked: stakingInfo.TotalStaked,
            securityScore: stakingInfo.SecurityScore,
            projectedROI: stakingInfo.ProjectedROI
        };
    }
    return null;
}

// Function to request VGL from the faucet
async function requestFaucetVGL(address) {
    try {
        const response = await fetch('/api/requestfaucetvgl', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ address: address })
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error(`Error requesting VGL from faucet:`, error);
        return { success: false, message: error.message };
    }
}