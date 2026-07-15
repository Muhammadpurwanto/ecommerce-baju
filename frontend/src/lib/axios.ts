// File: frontend/src/lib/axios.ts

import axios, { InternalAxiosRequestConfig, AxiosError } from "axios";
import Cookies from "js-cookie";

const baseURL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

const api = axios.create({
  baseURL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Flag untuk mencegah loop saat banyak request gagal secara bersamaan
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (reason?: any) => void;
}> = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => { 
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

// 1. Request Interceptor: Sisipkan Access Token
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    let token;
    if (typeof window !== "undefined") {  // cek apakah browser atau server yang request (SSR)
      token = Cookies.get("access_token");
    }

    if (token && config.headers) {  // jika ada token, tambahkan ke header
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 2. Response Interceptor: Tangkap 401 dan jalankan Refresh Token
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    // originalRequest = Menyimpan request asli yang gagal
    /* {
      url: "/profile",
      method: "GET",
      headers: {...}
    } */
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    // Jika error 429 (Rate Limiter)
    if (error.response?.status === 429) {
      console.error("Terlalu banyak request. Harap tunggu beberapa saat.");
      if (typeof window !== "undefined") {
        alert("Anda mengirim terlalu banyak request. Harap tunggu sebentar.");
      }
      return Promise.reject(error);
    }

    // Jika error 401 dan belum pernah di-retry
    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) { //cek apakah error 401 dan belum pernah di-retry
      if (isRefreshing) { //cek apakah sedang refresh
        // Jika sedang refresh, request yang gagal ini ikut antre menunggu hasil
        return new Promise(function (resolve, reject) { //buat promise untuk menampung request yang gagal
          failedQueue.push({ resolve, reject });
        })
          .then((token) => { //setelah selesai refresh
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`;
            }
            return api(originalRequest); // ulangi request dengan token baru
          })
          .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refreshToken = Cookies.get("refresh_token");

      if (!refreshToken) {
        // Tidak ada refresh token, log out otomatis
        Cookies.remove("access_token");
        isRefreshing = false;
        if (typeof window !== "undefined") window.location.href = "/login";
        return Promise.reject(error);
      }

      try {
        // Tembak endpoint refresh token bawaan backend kita
        const { data } = await axios.post(`${baseURL}/auth/refresh`, {
          refresh_token: refreshToken,
        });

        const newAccessToken = data.data.access_token;
        const newRefreshToken = data.data.refresh_token;

        // Simpan token baru ke cookies
        Cookies.set("access_token", newAccessToken, { expires: 1 }); // 1 hari
        Cookies.set("refresh_token", newRefreshToken, { expires: 7 }); // 7 hari

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        }

        processQueue(null, newAccessToken);
        
        // Ulangi request aslinya dengan token baru
        return api(originalRequest);
        
      } catch (refreshError) {
        // Refresh token kedaluwarsa atau tidak valid
        processQueue(refreshError, null);
        Cookies.remove("access_token");
        Cookies.remove("refresh_token");
        if (typeof window !== "undefined") window.location.href = "/login";
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

export default api;
