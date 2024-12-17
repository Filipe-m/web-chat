import UserRepository from "./repository.js";
import UserService from "./service.js";
import prisma from "../database/prisma.js";
import isEmpty from "../lib/isEmpty.js";

const userRepository = UserRepository(prisma);
const userService = UserService(userRepository);

export default {
  async registerUser(req, res) {
    const body = req.body;

    if (isEmpty(body.name) || isEmpty(body.email) || isEmpty(body.password)) {
      res.sendStatus(400);
    }

    try {
      const user = await userService.register(body);

      if (user === undefined) {
        res.status(400).json({ error: "Esse email já está cadastrado" });
      }

      res.status(201).json(user);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  },

  async login(req, res) {
    const body = req.body;

    if (isEmpty(body.email) || isEmpty(body.password)) {
      res.sendStatus(400);
      return;
    }

    try {
      const token = await userService.login(body);

      if(token === ""){
        res.sendStatus(401);
        return;
      }

      res.status(200).json({
        token: token,
      });
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  },

  async getUser(req, res) {
    try {
      const user = await userService.getUserById(req.params.id);
      if (!user) {
        return res.status(404).json({ message: "User not found" });
      }
      res.json(user);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  },
};
