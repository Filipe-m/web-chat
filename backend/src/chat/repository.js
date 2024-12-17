export default function ChatRepository(prisma) {
  return {
    async findById(id) {
      return await prisma.chat.findUnique({ where: { id } });
    },
    async create(data, userId) {
      return await prisma.chat.create({
        data: {
          name: data.name,
          created_by: userId,
        },
      });
    },
    async delete(chatId) {
      return await prisma.chat.delete({
        where: {
          id: chatId,
        },
      });
    },
    async getAllChats() {
      return await prisma.chat.findMany({
        include: {
          user: {
            select: {
              name: true,
            },
          },
        },
      });
    },
    async saveMessage(message, chatId, userId) {
      return await prisma.message.create({
        data: {
          content: message,
          user_id: userId,
          chat_id: chatId,
        },
      });
    },
    async chatInfo(chatId) {
      return await prisma.chat.findUnique({
        where: {
          id: chatId,
        },
      });
    },
    async messages(chatId, pages) {
      return await prisma.message.findMany({
        skip: pages * 10,
        take: 10,
        orderBy: {
          id: "desc",
        },
        include: {
          user: {
            select: {
              name: true,
            },
          },
        },
        where: {
          chat_id: chatId,
        },
      });
    },
  };
}
