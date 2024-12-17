import "dotenv/config";
import express from "express";
import { WebSocketServer } from "ws";

import prisma from "./database/prisma.js";
import UserRepository from "./user/repository.js";
import UserService from "./user/service.js";
import UserHandler from "./user/handler.js";

import ChatRepository from "./chat/repository.js";
import ChatService from "./chat/service.js";
import ChatHandler from "./chat/handler.js";
import { auth, websocketAuth } from "./middleware/auth.js";
import { globalMiddlewares } from "./middleware/global.js";

const userRepository = UserRepository(prisma);
const userService = UserService(userRepository);
const userHandler = UserHandler;

const chatRepository = ChatRepository(prisma);
const chatService = ChatService(chatRepository, UserRepository);
const chatHandler = ChatHandler;

const app = express();
const PORT = process.env.PORT || 3000;
const JWT_SECRET = process.env.JWT;

if (!JWT_SECRET) {
  console.error("JWT secret não está definido no arquivo .env.");
  process.exit(1);
}

app.use(globalMiddlewares);

app.get("/", auth, (req, res) => {
  res.send({ user: req.user });
});

app.post("/login", (req, res) => userHandler.login(req, res));
app.post("/register", (req, res) => userHandler.registerUser(req, res));

app.post("/chat", auth, (req, res) => chatHandler.create(req, res));
app.delete("/chat/:id", auth, (req, res) => chatHandler.delete(req, res));
app.get("/chat", auth, (req, res) => chatHandler.allChats(req, res));
app.get("/chat/:id", auth, (req, res) => chatHandler.chatInfo(req, res));
app.get("/messages/:id/:page", auth, (req, res) => chatHandler.messages(req, res));

const server = app.listen(PORT, () => {
  console.log(`Servidor rodando na porta ${PORT}`);
});

const wss = new WebSocketServer({ server });

wss.on("connection", (ws, req) => {
  websocketAuth(ws, req, (ws, req) => {
    const url = req.url;
    const urlParts = url.split("?");

    const chatId = urlParts[0]?.split("/")[2];

    if (!chatId) {
      ws.close(1008, "ID da sala não fornecido.");
      return;
    }

    ChatHandler.handleConnection(ws, chatId, req.userId);
  });
});
