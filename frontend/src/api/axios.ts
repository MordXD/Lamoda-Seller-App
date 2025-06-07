import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'http://localhost:3001/api', // URL бэкенда
});

// Интерцептор, который будет добавлять токен ко всем запросам
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('active_token'); // Активный токен
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default apiClient; 