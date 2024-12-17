import ChatRepository from "./repository.js";
import ChatService from "./service.js";
import prisma from "../database/prisma.js";
import isEmpty from "../lib/isEmpty.js";
import UserRepository from "../user/repository.js";

const chatRepository = ChatRepository(prisma);
const userRepository = UserRepository(prisma);
const chatService = ChatService(chatRepository, userRepository);

// Map para gerenciar conexões WebSocket por sala
const chatRooms = new Map();

export default {
  async create(req, res) {
    const body = req.body;

    if (isEmpty(body.name)) {
      res.sendStatus(400);
      return;
    }

    const userId = req.id;

    try {
      const chat = await chatService.create(body, userId);
      res.status(201).json(chat);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  },

  async delete(req, res) {
    const chatId = parseInt(req.params["id"]);
    const userId = req.id;

    try {
      const chat = await chatService.delete(chatId, userId);

      if (chat == undefined) {
        res.sendStatus(204);
        return;
      }

      res.status(200).json(chat);
    } catch (error) {
      res.status(500).json({ error: error });
    }
  },

  async allChats(req, res) {
    try {
      const chats = await chatService.getChats();
      res.status(200).json(chats);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  },
  async chatInfo(req, res) {
    const chatId = parseInt(req.params["id"]);

    const chat = await chatService.chatInfo(chatId);
    if (chat == null) {
      res.sendStatus(404);
    }

    res.status(200).json(chat);
  },
  async messages(req, res) {
    const pages = parseInt(req.params["page"]);
    const chatId = parseInt(req.params["id"]);

    const messages = await chatService.messages(chatId, pages);
    if (messages == null) {
      res.sendStatus(404);
    }

    res.status(200).json(messages);
  },

  handleConnection(ws, chatId, userId) {
    if (!chatRooms.has(chatId)) {
      chatRooms.set(chatId, new Set());
    }

    const clients = chatRooms.get(chatId);

    if (clients.has(ws)) {
      console.log(`O cliente já está conectado à sala ${chatId}.`);
      return;
    }

    clients.add(ws);
    console.log(
      `Cliente conectado à sala ${chatId}. Total de clientes: ${clients.size}`
    );

    ws.on("message", async (m) => {
      const message = JSON.parse(m);

      await chatService.saveMessage(message.content, parseInt(chatId), userId);

      await chatService.broadcastMessage(
        chatRooms,
        chatId,
        ws,
        message.content,
        userId
      );
    });

    ws.on("close", () => {
      clients.delete(ws);
      console.log(
        `Cliente desconectado da sala ${chatId}. Restantes: ${clients.size}`
      );
      if (clients.size === 0) {
        chatRooms.delete(chatId);
      }
    });
  },
};
