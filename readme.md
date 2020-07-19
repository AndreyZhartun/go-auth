# Сервис аутентификации на Go с Access и Refresh токенами
## Четыре REST маршрута:
1. Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
2. Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов
3. Третий маршрут удаляет конкретный Refresh токен из базы
4. Четвертый маршрут удаляет все Refresh токены из базы для конкретного пользователя
Access токен типа JWT, алгоритм SHA512.
Refresh токен типа JWT, хранится в базе исключительно в виде bcrypt хеша, защищен от изменения на стороне клиента.  
Access, Refresh токены обоюдно связаны случайно генерируемой строкой (см. *receive.go*), Refresh операцию для Access токена
можно выполнить только тем Refresh токеном который был выдан вместе с ним.
## Технологии: 
- Golang  
- MongoDB Atlas (Replica Set, транзакции)
- Heroku
## Результат:
Работающее приложение на [Heroku](https://sheltered-reef-38969.herokuapp.com).
HTTP маршруты могут быть проверены запросами к:
1. /receive?guid=xxxxxxxx-...-xxxx
2. /refresh
3. /remove
4. /removeall
# Authentication with Access&Refresh tokens in Golang
## Four REST routes:
1. Returns a pair of Access&Refresh tokens for an user if the user's GUID is supplied in GET params
2. Refreshes the Access token if a valid Access&Refresh token pair is supplied in HTTP-only cookie
3. Removes the Refresh token from the database if a valid Access&Refresh token pair is supplied in HTTP-only cookie
4. Removes all Refresh tokens from the database if a valid Access&Refresh token pair is supplied in HTTP-only cookie
Access token is JWT with SHA512 alg.
Refresh token is JWT, hashed with bcrypt and stored in the database.
Access, Refresh are linked with a secure random string (-> *receive.go*), preventing a possibility of refreshing Access token with a Refresh token that wasn't issued along with it.
## Used: 
- Golang  
- MongoDB Atlas (Replica Set, transactions)
- Heroku
## Results:
The app on [Heroku](https://sheltered-reef-38969.herokuapp.com).
HTTP routes are:
1. /receive?guid=xxxxxxxx-...-xxxx
2. /refresh
3. /remove
4. /removeall