# Usar uma imagem base do Golang
FROM golang:1.18-alpine

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Baixar as dependências (caso existam)
RUN go mod download

# Copiar os arquivos da API para o contêiner
COPY ./server /app

# Instalar as dependências e compilar a aplicação
RUN go build -o server

# Expor a porta em que a API será executada
EXPOSE 8080

# Comando para rodar a API
CMD ["./server"]

# Adicionar a etapa de verificação dos arquivos copiados e do build
RUN ls -al /app
RUN echo "Verificando a construção do binário"
RUN file /app/server
