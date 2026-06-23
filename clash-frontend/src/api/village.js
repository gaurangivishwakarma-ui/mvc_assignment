import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const getVillageProfile = async () => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.get(`${API_URL}/player/dashboard`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        throw error.response?.data?.error || "Failed to load village data";
    }
};

export const getShopCatalog = async () => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.get(`${API_URL}/shop/catalog`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        throw error.response?.data?.error || "Failed to load shop catalog";
    }
};

export const getArmyCatalog = async () => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.get(`${API_URL}/army/catalog`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        throw error.response?.data?.error || "Failed to load army catalog";
    }
};

export const getArmyStatus = async () => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.get(`${API_URL}/army/status`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        throw error.response?.data?.error || "Failed to load army status";
    }
};

export const trainTroops = async (data) => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.post(`${API_URL}/army/train`, data, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to train troops";
    }
};

export const purchaseBuilding = async (data) => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.post(`${API_URL}/village/buildings`, data, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to purchase building";
    }
};

export const upgradeVillage = async () => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.post(`${API_URL}/village/upgrade`, {}, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to upgrade village";
    }
};

export const getVillageUpgradeCost = async () => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.get(`${API_URL}/village/upgrade/cost`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to fetch upgrade cost";
    }
};

export const upgradeBuilding = async (data) => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.post(`${API_URL}/village/buildings/upgrade`, data, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to upgrade building";
    }
};

export const getBuildingUpgradeCost = async (placementId) => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.get(`${API_URL}/village/buildings/upgrade/cost?placement_id=${placementId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to fetch building upgrade cost";
    }
};

export const completeBuildingUpgrade = async (data) => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.post(`${API_URL}/village/buildings/complete`, data, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to complete building upgrade";
    }
};

export const moveBuilding = async (data) => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.post(`${API_URL}/village/move-building`, data, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to move building";
    }
};

export const getMatch = async () => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.get(`${API_URL}/battle/match`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to find match";
    }
};

export const attackOpponent = async (data) => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.post(`${API_URL}/battle/attack`, data, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Attack failed";
    }
};

export const collectResources = async () => {
    const token = localStorage.getItem('token');

    if (!token) {
        throw new Error("No commander token found. Please log in again.");
    }

    try {
        const response = await axios.post(`${API_URL}/village/collect`, {}, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error) {
        const backendMessage = typeof error.response?.data === 'string' 
            ? error.response.data.trim() 
            : error.response?.data?.error;
        throw backendMessage || "Failed to collect resources";
    }
};