# Étape 1 : Utiliser une image Golang de base
FROM golang:1.21

# Étape 2 : Créer un répertoire de travail
WORKDIR /app

# Étape 3 : Copier les fichiers Go dans le container
COPY . .

# Étape 4 : Télécharger les dépendances
RUN go mod tidy

# Étape 5 : Compiler ton programme Go
RUN go build -o main .

# Étape 6 : Exposer le port de ton serveur
EXPOSE 8080

# Étape 7 : Lancer l'application
CMD ["./main"]
