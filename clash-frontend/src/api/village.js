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
        throw error.response?.data?.error || "Failed to train troops";
    }
};