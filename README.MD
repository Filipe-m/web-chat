
# Web-Chat

Web-Chat é um projeto de chat feito utilizando de Websockets para comunicação de diferentes usuário dentro de uma rede local

## Stack utilizada

**Front-end:** React, TailwindCSS, Axios

**Back-end:** Go, Fiber, JWT, Testify

**Mobile:** React Native, Expo, Axios
## Aprendizados

Ao longo do desenvolvimento, aprofundei meu conhecimento na construção de APIs eficientes com Go e Fiber, explorando autenticação JWT e práticas de testes com Testify. No front-end, refinei a utilização de React com TailwindCSS para criar interfaces responsivas e elegantes, enquanto no mobile, o uso de React Native e Expo facilitou a entrega multiplataforma. A integração entre essas camadas com Axios reforçou a importância de um fluxo de dados bem estruturado.

Quais desafios enfrentei e como os superei?

Autenticação JWT e segurança
        Desafio: Gerenciar autenticação segura e eficiente no back-end.
        Solução: Implementei middlewares para validar tokens JWT e refinei a estratégia de refresh tokens para melhorar a experiência do usuário.

Testes unitários no back-end
        Desafio: Criar testes eficientes sem depender de um banco de dados real.
        Solução: Usei a biblioteca Testify e mocks para isolar dependências e garantir testes rápidos e confiáveis.

Estilização e responsividade
        Desafio: Garantir uma interface agradável e adaptável em diferentes dispositivos.
        Solução: TailwindCSS simplificou a criação de estilos responsivos, permitindo um desenvolvimento mais ágil.

Comunicação entre o front-end e o back-end
        Desafio: Lidar com erros de requisição e tratamento de dados.
        Solução: Configurei interceptors no Axios para padronizar a manipulação de respostas e erros.

Esse projeto trouxe um grande aprendizado sobre arquitetura escalável, boas práticas e otimização de fluxo entre as tecnologias utilizadas. 🚀

## Instalação

Clone o projeto

```bash
  git clone https://github.com/Filipe-m/web-chat
```

Entre no diretório do projeto

```bash
  cd web-chat
```

Backend

```bash
  cd backend
  go mod tidy
  go run ./cmd/main.go
```

Frontend

```bash
  cd frontend
  npm install
  npm run dev
```
Mobile

```bash
  cd mobile
  npm install
  npx expo start
```
    
## Deploy

Para fazer o deploy desse projeto rode

```bash
  docker compose up
```


## Variáveis de Ambiente

Para rodar esse projeto, você vai precisar adicionar as seguintes variáveis de ambiente no seu .env

`DATABASE_URL`

`JWT`


## Rodando os testes

Para rodar os testes, rode o seguinte comando

```bash
  cd backend
  go test ./...
```

