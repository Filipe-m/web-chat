import jwt from "jsonwebtoken";

const JWT_SECRET = process.env.JWT;

if (!JWT_SECRET) {
  console.error("JWT secret não está definido no arquivo .env.");
  process.exit(1);
}

export function auth(req, res, next) {
  const token = req.headers["authorization"]?.split(" ")[1];
  if (!token) {
    return res.status(401).send({ message: "Token não fornecido." });
  }

  jwt.verify(token, JWT_SECRET, (err, decoded) => {
    if (err) {
      return res.status(401).send({ message: "Token inválido." });
    }
    req.id = decoded.id;
    next();
  });
}

export function websocketAuth(ws, req, next) {
  const urlParams = new URLSearchParams(req.url.split('?')[1]);
  const token = urlParams.get('token');

  if (!token) {
    ws.close(1008, "Token não fornecido.");
    return;
  }

  jwt.verify(token, JWT_SECRET, (err, decoded) => {
    if (err) {
      return;
    }

    req.userId = decoded.id;
    next(ws, req);
  });
}
