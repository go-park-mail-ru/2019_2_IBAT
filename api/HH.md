# Документация к HH


<!-- ## API взаимодействия клиента и сервера HH IBAT -->


# SHORT SPECIFICATION

## AUTH

- "/auth"
  - POST - email и пароль для авторизации пользователя

- "/auth"
  - DELETE - выход пользователя из сессии   

## SEEKER METHODS

- "/seeker"
  - POST - регистрация пользователя

- "/seeker"
  - GET - получение текущих настроек кандидата
  - PUT - изменение текущих настроек кандидата

- "/seeker"
  - DELETE - удаление пользователя

- "/resume/<id>"
  - POST - создание резюме
  - PUT - изменение резюме
  - DELETE - удаление резюме

## EMPLOYER METHODS

- "/employer"
  - POST - регистрация компании 

- "/employer"
  - GET - получение текущих настроек компании
  - PUT - изменение текущих настроек компании

- "/employer"
  - DELETE - удаление компании

- "/vacancy/<id>"
  - POST - создание вакансии компании
  - PUT - изменение вакансии компании
  - DELETE - удаление вакансии компании
  
## COMMON METHODS

- "/employer/<id>"
  - GET - персональная страница компании

- "/vacancy/<id>"
  - GET - персональная страница вакансии

- "/resume/<id>"
  - GET - персональная страница соискателя

- "/employers"
  - GET - запрос списка работодателей(по фильтрам)

- "/resumes"
  - GET - запрос списка сооискателей(по фильтрам)

- "/vacancies"
  - GET - запрос списка вакансий(по фильтрам)


# METHODS DESCRIPTION

## AUTH

###  "/auth" POST
POST - email и пароль для авторизации пользователя

Структура JSON тела запроса

        {
            "email": __message__(string)
            "password": __password__(string)
        }

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK
        Cookie: name=value


    Отрицательный

        HTTP/1.1 400 Bad request

###  "/auth" DELETE
DELETE - разрыв сессии

Структура JSON тела запроса

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 204 

## SEEKER METHODS

### "/seeker" POST
POST - регистрация пользователя

Структура JSON тела запроса

       {
            "email": __message__(string)
            "first_name": __first_name__(string)
            "second_name": __second_name__(string)
            "password": __password__(string)
        }

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 400 


 


<!-- ## "/auth/"

    TODO -->


### "/seeker" GET PUT
GET - получение текущих настроек кандидата
PUT - изменение текущих настроек кандидата

GET
Структура JSON тела запроса

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            TODO implement settings data according to DB structure
        }

    Отрицательный

        HTTP/1.1 400   

PUT
Структура JSON тела запроса
    
    {
        TODO implement settings data according to DB structure
    }


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK


    Отрицательный

        HTTP/1.1 400 

### "/seeker" DELETE
DELETE - удаление пользователя

Структура JSON тела запроса

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 204

### "/resume/<id>" POST PUT DELETE
POST - создание резюме

Структура JSON тела запроса
    
    {
        TODO resume fields
    }


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 400


PUT - изменение резюме
Структура JSON тела запроса
    
    {
        TODO implement resume data according to DB structure
    }

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK


    Отрицательный

        HTTP/1.1 400 


DELETE - удаление резюме

Структура JSON тела запроса

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 204 

## EMPLOYER METHODS

### "/employer" POST
POST - регистрация компании

Структура JSON тела запроса

       {
            "email": __message__(string)
            "first_name": __first_name__(string)
            "second_name": __second_name__(string)
            "org_name":  __org_name__(string)
            "org_format":  __org_format__(string)
            "password": __password__(string)
        }

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 400

### "/employer" GET PUT
GET - получение текущих настроек работодателя
PUT - изменение текущих настроек работодателя

GET
Структура JSON тела запроса

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            TODO implement settings data according to DB structure
        }

    Отрицательный

        HTTP/1.1 400   

PUT
Структура JSON тела запроса
    
    {
        TODO implement settings data according to DB structure
    }


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK


    Отрицательный

        HTTP/1.1 400 

### "/employer" DELETE
DELETE - удаление компании

Структура JSON тела запроса

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 204



### "/vacancy/<id>" POST PUT DELETE

POST - создание вакансии
Структура JSON тела запроса
    
    {
        TODO vacancy fields
    }


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 400  


PUT - изменение вакансии
Структура JSON тела запроса
    
    {
        TODO implement vacancy data according to DB structure
    }


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK


    Отрицательный

        HTTP/1.1 400 



DELETE - удаление вакансии
Структура JSON тела запроса

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 204 

## COMMON METHODS

### "/employer/<id>"
GET - персональная страница компании


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            TODO implement employer info according to DB structure
        }

    Отрицательный

        HTTP/1.1 400 


### "/vacancy/<id>"
GET - персональная страница вакансии

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            TODO implement vacancy info according to DB structure
        }

    Отрицательный

        HTTP/1.1 400

### "/resume/<id>"
GET - резюме соискателя 

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            TODO implement seeker info according to DB structure
        }

    Отрицательный

        HTTP/1.1 400

### "/employers"
GET - запрос списка работодателей(по фильтрам)

TODO describe possible GET request flags

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            TODO implement employer info according to DB structure
        }

    Отрицательный

        HTTP/1.1 400


### "/resumes"
GET - запрос списка резюме(по фильтрам)

TODO describe possible GET request flags

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            TODO implement seeker info according to DB structure
        }

    Отрицательный

        HTTP/1.1 400



### "/vacancies"
GET - запрос списка вакансий(по фильтрам)

TODO describe possible GET request flags

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            TODO implement seeker info according to DB structure
        }

    Отрицательный

        HTTP/1.1 400
