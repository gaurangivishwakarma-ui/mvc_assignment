import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const registerPlayer = async (username, password) => {
    try {
        const response = await axios.post(`${API_URL}/register`, { username, password });
        return response.data;
    } catch (error) {
        throw error.response?.data?.error || "Registration failed";
    }
};

export const loginPlayer = async (username, password) => {
    try {
        const response = await axios.post(`${API_URL}/login`, { username, password });
        if (response.data.token) {
            localStorage.setItem('token', response.data.token);
        }
        return response.data;
    } catch (error) {
        throw error.response?.data?.error || "Login failed";
    }
};