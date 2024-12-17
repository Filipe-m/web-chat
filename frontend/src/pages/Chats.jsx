import React, { useEffect, useState, useCallback } from "react";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import axiosClient from "../api/axiosClient";
import { useNavigate } from "react-router-dom";
import { CreateServer } from "../components/createServer";

export default function Chats() {
  const [chats, setChats] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [reload, setReload] = useState(false);

  const navigate = useNavigate();

  const token = localStorage.getItem("token");

  const parseJwt = useCallback((token) => {
    try {
      return JSON.parse(atob(token.split(".")[1]));
    } catch (e) {
      return null;
    }
  }, []);

  const decoded = parseJwt(token);
  const userId = decoded?.id;

  const fetchChats = useCallback(async () => {
    try {
      setLoading(true);
      const response = await axiosClient.get("/chat");
      setChats(response.data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchChats();
  }, [fetchChats, reload]);

  const handleClick = (id) => {
    navigate(`/chat/${id}`);
  };

  const deleteChat = useCallback(async (id) => {
    try {
      await axiosClient.delete(`/chat/${id}`);
      setReload((prev) => !prev); // Atualiza a lista após exclusão
    } catch (err) {
      console.error("Erro ao deletar chat:", err);
    }
  }, []);

  if (loading) return <p className="p-16">Carregando chats...</p>;
  if (error) return <p className="p-16 text-red-500">Erro: {error}</p>;

  return (
    <div className="p-2 py-4 lg:p-16">
      <CreateServer setReload={setReload} />
      <Table className="my-4">
        <TableCaption>Lista de chats disponíveis</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead className="hidden lg:table-cell md:table-cell">
              ID
            </TableHead>
            <TableHead>Nome</TableHead>
            <TableHead>Criado por</TableHead>
            <TableHead className="hidden lg:table-cell md:table-cell">
              Criado em
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {chats.map((chat) => (
            <TableRow key={chat.id}>
              <TableCell className=" hidden lg:table-cell md:table-cell font-medium">
                {chat.id}
              </TableCell>
              <TableCell>{chat.name}</TableCell>
              <TableCell>{chat.user.name}</TableCell>
              <TableCell className="hidden lg:table-cell md:table-cell">
                {new Date(chat.created_at).toLocaleString()}
              </TableCell>
              <TableCell>
                <Button
                  className="hover:ring-4 hover:ring-slate-300"
                  onClick={() => handleClick(chat.id)}
                >
                  Conversa
                </Button>
              </TableCell>
              <TableCell>
                {chat.created_by === userId && (
                  <Button
                    className="bg-red-600 hover:bg-red-600 hover:ring-4 hover:ring-red-200"
                    onClick={() => deleteChat(chat.id)}
                  >
                    Deletar
                  </Button>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
