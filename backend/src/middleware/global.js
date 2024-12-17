import cors from "cors";
import express from "express";

export const globalMiddlewares = [
  cors({
    origin: (origin, callback) => {
      const allowedOrigins = [
        "http://localhost:3000",
        "http://localhost:5173",
      ];
  
      // Verifica se a origem é local (IP na sub-rede 192.168.x.x ou 10.x.x.x)
      const isLocalNetwork = origin && (origin.startsWith("http://192.168.") || origin.startsWith("http://10."));
  
      if (!origin || allowedOrigins.includes(origin) || isLocalNetwork) {
        callback(null, true); // Permite a origem
      } else {
        console.error(`Origem não permitida pelo CORS: ${origin}`);
        callback(new Error("Not allowed by CORS")); // Bloqueia a origem
      }
    },
    credentials: true,
    methods: "GET,HEAD,PUT,PATCH,POST,DELETE",
    allowedHeaders: "Content-Type, Authorization",
  }),
  
  express.json({
    verify: (req, res, buf) => {
      if (buf.length > 0) {
        try {
          JSON.parse(buf);
        } catch {
          res.status(400).send({ message: "Bad Request" });
          throw new Error("Invalid JSON");
        }
      }
    },
  }),
];
