import isEmpty from "../lib/isEmpty.js";

export default function ChatService(chatRepository, userRepository) {
  return {
    async create(chat, userId) {
      return await chatRepository.create(chat, userId);
    },
    async delete(chatId, userId) {
      const chat = await chatRepository.findById(chatId);

      if (chat == null || chat.created_by != userId) {
        return undefined;
      }

      return await chatRepository.delete(chatId);
    },
    async getChats(chatId) {
      return await chatRepository.getAllChats(chatId);
    },
    async broadcastMessage(chatRooms, chatId, ws, message, userId) {
      const user = await userRepository.findById(userId);

      const payload = {
        content: message,
        user: {
          name: user.name,
        },
        user_id: user.id,
        created_at: new Date().toISOString(),
        chat_id: chatId,
      };

      const clients = chatRooms.get(chatId);
      if (clients) {
        clients.forEach((client) => {
          if (client !== ws && client.readyState === ws.OPEN) {
            client.send(JSON.stringify(payload));
          }
        });
      }
    },
    async saveMessage(message, chatId, userId) {
      if (isEmpty(message) || isEmpty(chatId) || isEmpty(userId)) {
        return undefined;
      }
      return await chatRepository.saveMessage(message, chatId, userId);
    },
    async chatInfo(chatId) {
      return await chatRepository.chatInfo(chatId);
    },
    async messages(chatId, pages) {
      return await chatRepository.messages(chatId, pages);
    },
  };
}
