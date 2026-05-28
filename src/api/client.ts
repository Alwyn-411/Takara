import axios from 'axios';
import { useUserStore } from '../store/User';

export const axiosInstance = axios.create({
    baseURL: 'http://localhost:8080',
    headers: {
        'Content-Type': 'application/json',
    },
});

axiosInstance.interceptors.request.use((config) => {
    const token = useUserStore.getState().token;
    if (token) config.headers.Authorization = `Bearer ${token}`;
    return config;
});

export interface PingResponse {
    message: string;
}

export const ping = async () => {
    const response = await axiosInstance.get<PingResponse>('/ping');
    return response.data;
};
