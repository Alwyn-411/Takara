import axios from "axios";

export const axiosInstance = axios.create({
  baseURL: "http://localhost:8080",
  headers: {
    "Content-Type": "application/json",
  },
});

export interface PingResponse {
  message: string;
}

export const ping = async () => {
  const response = await axiosInstance.get<PingResponse>("/ping");
  return response.data;
};
