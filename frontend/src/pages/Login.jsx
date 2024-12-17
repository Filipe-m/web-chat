import React, { useState } from "react";
import axiosClient from "../api/axiosClient";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

function Login() {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });

  const [errorMessage, setErrorMessage] = useState("");

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrorMessage("");

    try {
      const response = await axiosClient.post("/login", formData);
      const token = response.data.token;

      localStorage.setItem("token", token);

      window.location.href = "/chats";
    } catch (error) {
      if (error.response && error.response.data) {
        setErrorMessage(error.response.data.message || "Erro ao fazer login");
      } else {
        setErrorMessage("Erro ao conectar ao servidor");
      }
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-md shadow-md">
        <h1 className="text-2xl font-bold text-center text-gray-800">Login</h1>
        <form className="space-y-4" onSubmit={handleSubmit}>
          <div>
            <label
              htmlFor="email"
              className="block text-sm font-medium text-gray-700"
            >
              Email
            </label>
            <Input
              type="email"
              id="email"
              name="email"
              placeholder="Digite seu email"
              className="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              value={formData.email}
              onChange={handleChange}
            />
          </div>
          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-700"
            >
              Senha
            </label>
            <Input
              type="password"
              id="password"
              name="password"
              placeholder="Digite sua senha"
              className="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              value={formData.password}
              onChange={handleChange}
            />
          </div>
          {errorMessage && (
            <p className="text-sm text-red-500">{errorMessage}</p>
          )}
          <Button
            type="submit"
            className="w-full py-2 mt-4 text-white rounded-md "
          >
            Entrar
          </Button>
        </form>
        <p className="mt-4 text-sm text-center text-gray-600">
          Ainda não possui uma conta?{" "}
          <a
            href="/register"
            className="text-slate-900 font-bold hover:underline"
          >
            Cadastre-se
          </a>
        </p>
      </div>
    </div>
  );
}

export default Login;
