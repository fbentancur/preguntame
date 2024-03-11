# Proyecto Preguntame

## Descripción
Proyecto sobre una aplicación imaginaria tipo ask.fm que incluye posteo de posts

## Tecnologías
- La aplicación usa SQL para manejar la base de datos en PostgreSQL
- El código está hecho en golang usando ECHO como framework http

## Endpoints
- ´POST /users/login´ 
Sirve para hacer log a la pagina y genera un token aleatorio para el usuario registrado
- ´POST /users/register´ 
Sirve para registrar al usuario añadiendo una fila a la base de datos con su usuario y contraseña
- ´GET /users/:user_id/questions´
Sirve para buscar las preguntas que hacen referencia al id de usuario correspondiente
- ´POST /users/:user_id/questions´
Sirve para hacer una pregunta a un usuario usando el id como parametro, la pregunta se almacena en la base de datos con el id de usuario como       referencia
- ´PUT /users/:user_id/questions/:question_id´
Sirve para responder una pregunta realizada el usuario, el endpoint compara que el id de usuario al que se hizo la pregunta sea coincidente con el token de logueo del usuario que responde
- ´PUT /users/:user_id/questions/:question_id/fav´
Sirve para que el usuario pueda marcar una pregunta como favorita
- ´DELETE /users/:user_id/questions/:question_id´
Sirve para hacer un hard delete a una pregunta, el endpoint compara que el id de usuario al que se hizo la pregunta sea coincidente con el token del logueo del usuario que quiere borrar la pregunta

- ´POST /users/:user_id/posts´
Sirve para crear un post y agregarlo al feed de quien lo crea, el endpoint verifica que el id del dueño del feed sea coincidente con el token de logueo del usuario que postea
- ´PATCH /users/:user_id/posts/:post_id´
Sirve para modificar un post, el endpoint verifica que el id del usuario dueño del post sea coincidente con el token de logueo del usuario que busca modificarlo
- ´DELETE /users/:user_id/posts/:post_id´
Sirve para hacer un soft delete de un post, el endpoint verifica que el id del usuario dueño del post sea coincidente con el token de logueo de usuario que busca borrarlo