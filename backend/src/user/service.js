import { encryptPassword, verifyPassword } from "../lib/bcrypt.js";
import jwt from "jsonwebtoken";
import "dotenv/config";

const JWT_SECRET = process.env.JWT;

export default function UserService(userRepository) {
  return {
    async register(data) {
      const findUser = await userRepository.findByEmail(data.email);

      if (findUser !== null) {
        return undefined;
      }

      data.password = await encryptPassword(data.password);

      return await userRepository.create(data);
    },

    async login(data) {
      const user = await userRepository.findByEmail(data.email);

      const samePassword = await verifyPassword(data.password, user.password);

      if (!samePassword || user === null) {
        return "";
      }

      const payload = {
        id: user.id
      };

      const token = jwt.sign(payload, JWT_SECRET, { expiresIn: "1h" });

      return token;
    },
  };
}
