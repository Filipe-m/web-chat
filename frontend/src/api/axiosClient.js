import axios from "axios";

// Crie uma instância do Axios
const axiosClient = axios.create({
    baseURL: "http://192.168.0.154:3000", // Defina a base URL da sua API
    withCredentials: true, // Inclui cookies e credenciais, se necessário
    headers: {
        "Content-Type": "application/json",
    },
});

// Interceptor para adicionar o token a todas as requisições
axiosClient.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem("token");
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Interceptor para tratar respostas e erros globais
axiosClient.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        // Por exemplo, redirecionar para login se o token for inválido/expirar
        if (error.response && error.response.status === 401) {
            localStorage.removeItem("token");
            window.location.href = "/login";
        }
        return Promise.reject(error);
    }
);

export default axiosClient;
