export default function UserRepository(prisma) {
  return {
    async findById(id) {
      return await prisma.user.findUnique({ where: { id } });
    },
    async create(data) {
      return await prisma.user.create({ data });
    },
    async findByEmail(email) {
      return await prisma.user.findFirst({
        where: {
          email: email,
        },
      });
    },
  };
}
