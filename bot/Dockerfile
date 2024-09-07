# Usar uma imagem base do Golang
FROM golang:1.20-alpine

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /bot

# Copiar os arquivos da API para o contêiner
COPY ./ /bot

# Instalar as dependências e compilar a aplicação
RUN go build -o bot

# Expor a porta em que a API será executada
EXPOSE 8079

# Comando para rodar a API
CMD ["./bot"]
