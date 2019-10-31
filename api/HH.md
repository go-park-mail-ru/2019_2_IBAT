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

- "/resume"
  - POST - создание резюме

- "/resume/<id>"
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

- "/vacancy"
  - POST - создание вакансии компании

- "/vacancy/<id>"
  - PUT - изменение вакансии компании
  - DELETE - удаление вакансии компании
  
## COMMON METHODS

- "/employer/<id>"
  - GET - персональная страница компании

- "/vacancy/<id>"
  - GET - персональная страница вакансии

- "/resume/<id>"
  - GET - персональная страница соискателя

- "/seeker/<id>"
  - GET - персональная страница соискателя

- "/employers"
  - GET - запрос списка работодателей(по фильтрам)

- "/resumes"
  - GET - запрос списка сооискателей(по фильтрам)

- "/vacancies"
  - GET - запрос списка вакансий(по фильтрам)

  - "/responds"
  - GET - запрос списка вакансий(по фильтрам)


# METHODS DESCRIPTION

Примечание:
JSON структура на каждый отрицательный респонс
{
    "error": "message"
} 
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
        Set-Cookie: name=value
        {
            "class": "seeker/employer"
        }


    Отрицательный

        HTTP/1.1 400 Bad request

###  "/auth" DELETE
DELETE - разрыв сессии

Структура JSON тела запроса

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK
        Set-Cookie: name=value expires=new_date

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
        {
            "class": "seeker"
        }

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
                "email": "somes@mail.com",
                "first_name": "Grisha",
                "second_name": "Zyablikov",
                "password": "111111",
                "resumes": {id, id} //array of ids
        }

    Отрицательный

        HTTP/1.1 400   

PUT
Структура JSON тела запроса
    
    {
        "email": "somes@mail.com",
        "first_name": "Grisha",
        "second_name": "Zyablikov",
        "password": "111111",
    }

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK


    Отрицательный

        HTTP/1.1 400
        {
            "error": "message"
        }

### "/seeker" DELETE
DELETE - удаление пользователя

Структура JSON тела запроса

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 204

### "/resume" POST
POST - создание резюме

Структура JSON тела запроса   

    {
        "first_name": "Vova",
        "second_name": "Zyablikov",
        "city": "Moscow",
        "number": "12345678910",
        "birth_date": "1994-21-08",
        "sex": "male",
        "citizenship": "Russia",
        "experience": "7 years",
        "profession": "programmer",
        "position": "middle",
        "wage": "100500",
        "education": "MSU",
        "about": "Hello employer"
    }


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK
        {
            "id": "id" (string)
        }

    Отрицательный

        HTTP/1.1 400

### "/resume/<id>"PUT DELETE

PUT - изменение резюме
Структура JSON тела запроса
    
    {
        "first_name": "Vova",
        "second_name": "Zyablikov",
        "city": "Moscow",
        "number": "12345678910",
        "birth_date": "1994-21-08",
        "sex": "male",
        "citizenship": "Russia",
        "experience": "7 years",
        "profession": "programmer",
        "position": "middle",
        "wage": "100500",
        "education": "MSU",
        "about": "Hello employer"
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
            "company_name": "BMSsTU",
            "site":"bmstu.ru",
            "email":"bmsstu@mail.com",
            "first_name":"Tolya",
            "second_name": "Alex",
            "password":"1234",
            "number": "12345678911",
            "extra_number": "12345678910",
            "city": "Moscow",
            "empl_num": 1830
        }

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK
        {
            "class": "employer"
        }

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
            "company_name": "BMSsTU",
            "site": "bmstu.ru",
            "first_name": "Tolya",
            "second_name": "Alex",
            "email": "bmsstu@mail.com",
            "number": "12345678911",
            "extra_number": "12345678910",
            "password": "1234",
            "city": "Moscow",
            "empl_num": 1830
        }

    Отрицательный

        HTTP/1.1 400   

PUT
Структура JSON тела запроса
    
        {
            "company_name": "BMSsTU",
            "site": "bmstu.ru",
            "first_name": "Tolya",
            "second_name": "Alex",
            "email": "bmsstu@mail.com",
            "number": "12345678911",
            "extra_number": "12345678910",
            "password": "1234",
            "city": "Moscow",
            "empl_num": 1830
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



### "/vacancy" POST 

POST - создание вакансии
Структура JSON тела запроса
    
    {
        "company_name": "BMSsTU",
        "experience": "3 years and more",
        "profession": "baker",
        "position":  "mid",
        "tasks": "writing test",
        "requirements": "should be able writing good tests",
        "wage": "100500",
        "conditions": "nice team",
        "about": "you will work in the best company"
    }

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK
        {
            "id": "id" (string)
        }

    Отрицательный

        HTTP/1.1 400  

### "/vacancy/<id>" PUT DELETE

PUT - изменение вакансии
Структура JSON тела запроса
    
   {
        "company_name": "BMSsTU",
        "experience": "3 years and more",
        "profession": "baker",
        "position":  "mid",
        "tasks": "writing test",
        "requirements": "should be able writing good tests",
        "wage": "100500",
        "conditions": "nice team",
        "about": "you will work in the best company"
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
            "company_name": "BMSsTU",
            "site": "bmstu.ru",
            "first_name": "Tolya",
            "second_name": "Alex",
            "email": "bmsstu@mail.com",
            "number": "12345678911",
            "extra_number": "12345678910",
            "password": "",                     //field should be eliminated
            "city": "Moscow",
            "empl_num": 1830,
            "vacancies": {"id1", "id2"}  //array of strings
        }

    Отрицательный

        HTTP/1.1 400 


### "/vacancy/<id>"
GET - персональная страница вакансии

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            "company_name": "BMSsTU",
            "experience": "3 years and more",
            "profession": "baker",
            "position":  "mid",
            "tasks": "writing test",
            "requirements": "should be able writing good tests",
            "wage": "100500",
            "conditions": "nice team",
            "about": "you will work in the best company"
        }

    Отрицательный

        HTTP/1.1 400

### "/resume/<id>"
GET - резюме соискателя 

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            "first_name": "Vova",
            "second_name": "Zyablikov",
            "city": "Moscow",
            "number": "12345678910",
            "birth_date": "1994-21-08",
            "sex": "male",
            "citizenship": "Russia",
            "experience": "7 years",
            "profession": "programmer",
            "position": "middle",
            "wage": "100500",
            "education": "MSU",
            "about": "Hello employer"
        }

    Отрицательный

        HTTP/1.1 400

###  "/seeker/<id>"
GET - персональная страница соискателя

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            "email": "somes@mail.com",
            "first_name": "Grisha",
            "second_name": "Zyablikov",
            "password": "",
            "resumes": {"id1", "id2"} //array of ids
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
            id1: {
                "company_name": "BMSsTU",
                "site": "bmstu.ru",
                "first_name": "Tolya",
                "second_name": "Alex",
                "email": "bmsstu@mail.com",
                "number": "12345678911",
                "extra_number": "12345678910",
                "password": "",                     //field should be eliminated
                "city": "Moscow",
                "empl_num": 1830,
                "vacancies": {"id1", "id2"},  //array of strings
            }
            
            id2: {
                "company_name": "BMSsTU",
                "site": "bmstu.ru",
                "first_name": "Tolya",
                "second_name": "Alex",
                "email": "bmsstu@mail.com",
                "number": "12345678911",
                "extra_number": "12345678910",
                "password": "",                     //field should be eliminated
                "city": "Moscow",
                "empl_num": 1830,
                "vacancies": {"id1", "id2"}  //array of strings
            }
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
            id: {
                "first_name": "Vova",
                "second_name": "Zyablikov",
                "city": "Moscow",
                "number": "12345678910",
                "birth_date": "1994-21-08",
                "sex": "male",
                "citizenship": "Russia",
                "experience": "7 years",
                "profession": "programmer",
                "position": "middle",
                "wage": "100500",
                "education": "MSU",
                "about": "Hello employer"
            }

            id: {
                "first_name": "Vova",
                "second_name": "Zyablikov",
                "city": "Moscow",
                "number": "12345678910",
                "birth_date": "1994-21-08",
                "sex": "male",
                "citizenship": "Russia",
                "experience": "7 years",
                "profession": "programmer",
                "position": "middle",
                "wage": "100500",
                "education": "MSU",
                "about": "Hello employer"
            }   
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
            id: {
                "company_name": "BMSsTU",
                "experience": "3 years and more",
                "profession": "baker",
                "position":  "mid",
                "tasks": "writing test",
                "requirements": "should be able writing good tests",
                "wage": "100500",
                "conditions": "nice team",
                "about": "you will work in the best company"
            }
            
            id: {
                "company_name": "BMSsTU",
                "experience": "3 years and more",
                "profession": "baker",
                "position":  "mid",
                "tasks": "writing test",
                "requirements": "should be able writing good tests",
                "wage": "100500",
                "conditions": "nice team",
                "about": "you will work in the best company"
            }
        }

    Отрицательный

        HTTP/1.1 400
