import React, { useEffect, useState, useRef, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { useParams, useNavigate } from "react-router-dom";
import axiosClient from "../api/axiosClient";
import { Input } from "@/components/ui/input";

export default function Chat() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [messages, setMessages] = useState([]);
  const [message, setMessage] = useState("");
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef(null);
  const chatContainerRef = useRef(null);
  const [chatInfo, setChatInfo] = useState({});
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [isLoading, setIsLoading] = useState(false);

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

  const addUniqueMessages = (newMessages) => {
    setMessages((prevMessages) => {
      const existingIds = new Set(prevMessages.map((msg) => msg.id));
      const filteredMessages = newMessages.filter(
        (msg) => !existingIds.has(msg.id)
      );
      return [...prevMessages, ...filteredMessages];
    });
  };

  const connectWebSocket = () => {
    const socketUrl = `ws://192.168.0.154:3000/chat/${id}?token=${token}`;
    socketRef.current = new WebSocket(socketUrl);

    socketRef.current.onopen = () => {
      setIsConnected(true);
    };

    socketRef.current.onmessage = (event) => {
      const newMessage = JSON.parse(event.data);
      addUniqueMessages([newMessage]);
    };

    socketRef.current.onerror = (error) => {
      console.error("Erro na conexão WebSocket:", error);
    };

    socketRef.current.onclose = () => {
      setIsConnected(false);
    };
  };

  const fetchChat = async () => {
    try {
      const response = await axiosClient.get(`/chat/${id}`);
      setChatInfo(response.data);
    } catch (err) {
      console.error(err);
    }
  };

  const fetchMessages = async (currentPage) => {
    if (isLoading || !hasMore) return;
    setIsLoading(true);

    try {
      const response = await axiosClient.get(`/messages/${id}/${currentPage}`);
      const newMessages = response.data;

      if (newMessages.length === 0) {
        setHasMore(false);
      } else {
        addUniqueMessages(newMessages);
      }
    } catch (err) {
      console.error("Erro ao buscar mensagens:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const scrollToBottom = () => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop =
        chatContainerRef.current.scrollHeight;
    }
  };

  useEffect(() => {
    connectWebSocket();
    fetchChat();
    fetchMessages(page);

    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, [id]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const sendMessage = () => {
    if (socketRef.current && message.trim()) {
      const newMessage = {
        content: message,
        author: "Eu",
        user_id: userId,
        created_at: new Date().toISOString(),
        id: Date.now(),
      };

      socketRef.current.send(JSON.stringify(newMessage));
      setMessages((prevMessages) => [...prevMessages, newMessage]);
      setMessage("");
    }
  };

  const handleScroll = () => {
    if (!chatContainerRef.current) return;

    const { scrollTop } = chatContainerRef.current;
    if (scrollTop === 0 && hasMore && !isLoading) {
      const nextPage = page + 1;
      setPage(nextPage);
      fetchMessages(nextPage);
    }
  };

  return (
    <div className="h-screen w-screen grid place-items-center">
      <Button
        onClick={() => navigate("/chats")}
        className="absolute left-0 m-4 top-0"
      >
        Voltar
      </Button>
      <div className="bg-zinc-400 p-4 rounded-lg w-full max-w-2xl">
        <h1 className=" grid place-items-center font-bold text-2xl">
          {chatInfo.name}
        </h1>
        <div
          ref={chatContainerRef}
          onScroll={handleScroll}
          className="messages-container py-4 px-6 flex flex-col overflow-y-scroll max-h-96"
        >
          {messages
            .slice()
            .sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
            .map((msg) => (
              <div key={msg.id} className="message">
                <strong>
                  {msg.user_id === userId
                    ? "Eu"
                    : msg.user.name ?? "Desconhecido"}
                  :
                </strong>
                {msg.content ?? "Mensagem sem conteúdo"}{" "}
                <div className="text-sm text-gray-500">
                  {new Date(msg.created_at).toLocaleString()}{" "}
                </div>
              </div>
            ))}
        </div>
        <div className="message-input flex gap-3 mt-4">
          <Input
            className="rounded-lg pl-2 flex-grow"
            type="text"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Digite uma mensagem..."
            onKeyDown={(e) => {
              if (e.key === "Enter") sendMessage();
            }}
          />
          <Button onClick={sendMessage} disabled={!isConnected}>
            Enviar
          </Button>
        </div>
        {!isConnected && <p>Conectando ao chat...</p>}
      </div>
    </div>
  );
}
