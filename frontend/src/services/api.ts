import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

export const fetchContainerData = async () => {
    try {
        const response = await axios.get(`${API_URL}/pings`);
        return response.data;
    } catch (error) {
        console.error('Ошибка при получении данных:', error);
        return [];
    }
};
