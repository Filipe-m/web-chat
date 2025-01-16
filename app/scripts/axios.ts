import axios from "axios";
import * as SecureStore from "expo-secure-store";

const api = axios.create({
  baseURL: process.env.EXPO_PUBLIC_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(
  async (config) => {
    const accessToken = await SecureStore.getItemAsync("session");
    if (accessToken) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }

    // console.log("Request:", {
    //   url: config.url,
    //   method: config.method,
    //   headers: config.headers,
    //   data: config.data,
    // });
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// api.interceptors.response.use(
//   (response) => response,
//   async (error) => {
//     const originalRequest = error.config;
//     if (error.response?.status === 401 && !originalRequest._retry) {
//       originalRequest._retry = true;

//       const refreshToken = await SecureStore.getItemAsync("refreshToken");
//       if (refreshToken) {
//         try {
//           const { data } = await axios.post(
//             `${process.env.EXPO_PUBLIC_API_URL}/refresh`,
//             {
//               refreshToken,
//             }
//           );

//           const { accessToken: newAccessToken } = data;

//           await SecureStore.setItemAsync("accessToken", newAccessToken);
//           originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
//           return api(originalRequest);
//         } catch (refreshError) {
//           console.error("Erro ao renovar token:", refreshError);
//         }
//       }
//     }

//     return Promise.reject(error);
//   }
// );

export default api;
