# Этап сборки
FROM node:20-alpine as build

WORKDIR /app

# Копируем файлы зависимостей
COPY package.json yarn.lock ./

# Устанавливаем зависимости
RUN yarn install --frozen-lockfile

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN yarn build

# Этап продакшена
FROM nginx:alpine

# Копируем собранное приложение в nginx
COPY --from=build /app/dist /usr/share/nginx/html

# Копируем конфигурацию nginx для SPA
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Открываем порт 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"] 