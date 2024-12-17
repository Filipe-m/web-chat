import React from "react";
import { Navigate } from "react-router-dom";

const ProtectedRoute = ({ children }) => {
    const token = localStorage.getItem("token");

    if (!token) {
        return <Navigate to="/login" />;
    }

    try {
        // Decodifica o token JWT
        const parseJwt = (token) => {
            try {
                return JSON.parse(atob(token.split('.')[1]));
            } catch (e) {
                console.error("Erro ao decodificar o token:", e);
                return null;
            }
        };
        
        const decoded = parseJwt(token);

        // Verifica se o token está expirado
        const isTokenExpired = decoded.exp * 1000 < Date.now();

        if (isTokenExpired) {
            console.error("Token expirado");
            return <Navigate to="/login" />;
        }
    } catch (error) {
        console.error("Token inválido:", error);
        return <Navigate to="/login" />;
    }

    return children;
};

export default ProtectedRoute;
