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