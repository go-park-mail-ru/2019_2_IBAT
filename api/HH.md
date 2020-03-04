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
  - PUT - изменение текущих настроек кандидата

- "/seeker"
  - DELETE - удаление пользователя

- "/resume"
  - POST - создание резюме

- "/resume/<id>"
  - PUT - изменение резюме
  - DELETE - удаление резюмеc

- "/favorite_vacancies"
  - GET 

- "/favorite_vacancy/{id}"
  - POST - создание избранной вакансии
  - DELETE - удаление избранной вакансии (не сделано на бэке)

## EMPLOYER METHODS

- "/employer"
  - POST - регистрация компании 

- "/employer"
  - PUT - изменение текущих настроек компании

- "/employer"
  - DELETE - удаление компании

- "/vacancy"
  - POST - создание вакансии компании

- "/vacancy/<id>"
  - PUT - изменение вакансии компании
  - DELETE - удаление вакансии компании


## COMMON METHODS

- "/profile/"
  - GET - получение профиля пользователя

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


### "/seeker" GET PUT
PUT - изменение текущих настроек кандидата
  

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
        "region": "Moscow",
        "number": "12345678910",
        "birth_date": "1994-21-08",
        "sex": "male",
        "citizenship": "Russia",
        "experience": "7 years",
        "position": "programmer",
        "wage": "100500",
        "education": "MSU",
        "about": "Hello employer",
        "spheres": [
            {"first": "Бухгалтерия, управленческий учет, финансы предприятия", "second": "бухгалтер"},
            {"first": "Бухгалтерия, управленческий учет, финансы предприятия", "second": "основные средства"}
	    ]

    }


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK
        {
            "id": "id" (string)
        }

    Отрицательный

        HTTP/1.1 400

### "/resume/<id>" PUT DELETE

PUT - изменение резюме
Структура JSON тела запроса
    
    {
        "first_name": "Vova",
        "second_name": "Zyablikov",
        "region": "Moscow",
        "number": "12345678910",
        "birth_date": "1994-21-08",
        "sex": "male",
        "citizenship": "Russia",
        "experience": "7 years",
        "position": "middle",
        "wage": "100500",
        "education": "MSU",
        "about": "Hello employer",
        "spheres": [
            {"first": "Бухгалтерия, управленческий учет, финансы предприятия", "second": "бухгалтер"},
            {"first": "Бухгалтерия, управленческий учет, финансы предприятия", "second": "основные средства"}
	    ]
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

### "/respond" POST

POST - создание отклика

Структура JSON тела запроса

       {
            "vacancy_id": "322",
            "resume_id": "244",
        }

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 400


### "/respond/<id>" (не сделано на бэке)

DELETE - удаление вакансии компании  

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 400


### "/favorite_vacancies"

GET 

Ответ на запрос:

    Положительный

        {
            {   
                "id": uuid
                "company_name": "BMSsTU",
                "experience": "3 years and more",
                "profession": "baker",
                "position":  "mid",
                "tasks": "writing test",
                "requirements": "should be able writing good tests",
                "wage_from": "100500",
                "conditions": "nice team",
                "about": "you will work in the best company"
            }
            
            {
                "id": uuid
                "company_name": "BMSsTU",
                "experience": "3 years and more",
                "position":  "mid",
                "tasks": "writing test",
                "requirements": "should be able writing good tests",
                "wage_from": "100500",
                "conditions": "nice team",
                "about": "you will work in the best company"
            }
        }

    Отрицательный

        HTTP/1.1 400


###  "/favorite_vacancy/{id}"
POST - создание избранной вакансии
    
Структура JSON тела запроса


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

    Отрицательный

        HTTP/1.1 400

        
DELETE - удаление избранной вакансии (не сделано на бэке)


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
            "region": "Moscow",
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
            "region": "Moscow",
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
            "region": "Moscow",
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
        "position":  "mid",
        "tasks": "writing test",
        "requirements": "should be able writing good tests",
        "wage": "100500",
        "conditions": "nice team",
        "about": "you will work in the best company"
        "spheres": [
            {"first": "Бухгалтерия, управленческий учет, финансы предприятия", "second": "бухгалтер"},
            {"first": "Бухгалтерия, управленческий учет, финансы предприятия", "second": "основные средства"}
	    ]
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
        "spheres": [
            {"first": "Бухгалтерия, управленческий учет, финансы предприятия", "second": "бухгалтер"},
            {"first": "Бухгалтерия, управленческий учет, финансы предприятия", "second": "основные средства"}
	    ]
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

### "/profile"
GET - профиль юзера
 
    Ответ на запрос:  

        Положительный (Вариант 1)

            HTTP/1.1 200 OK

            {
                "company_name": "BMSsTU",
                "site": "bmstu.ru",
                "first_name": "Tolya",
                "second_name": "Alex",
                "email": "bmsstu@mail.com",
                "number": "12345678911",
                "extra_number": "12345678910",
                "region": "Moscow",
                "empl_num": 1830,
                "vacancies": {"id1", "id2"}  //array of strings
            }
    
        Положительный (Вариант 2)

            HTTP/1.1 200 OK

            {
                "email": "somes@mail.com",
                "first_name": "Grisha",
                "second_name": "Zyablikov",
                "resumes": {"id1", "id2"} //array of ids
            }

### "/employer/<id>"
GET - персональная страница компании


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            "id": uuid
            "company_name": "BMSsTU",
            "site": "bmstu.ru",
            "first_name": "Tolya",
            "second_name": "Alex",
            "email": "bmsstu@mail.com",
            "number": "12345678911",
            "extra_number": "12345678910",
            "region": "Moscow",
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
            "id": uuid
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
            "id": uuid
            "first_name": "Vova",
            "second_name": "Zyablikov",
            "region": "Moscow",
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

Параметры

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            {
                "id": uuid
                "company_name": "BMSsTU",
                "site": "bmstu.ru",
                "first_name": "Tolya",
                "second_name": "Alex",
                "email": "bmsstu@mail.com",
                "number": "12345678911",
                "extra_number": "12345678910",
                "region": "Moscow",
                "empl_num": 1830,
                "vacancies": {"id1", "id2"},  //array of strings
            }
            
            {
                "id": uuid
                "company_name": "BMSsTU",
                "site": "bmstu.ru",
                "first_name": "Tolya",
                "second_name": "Alex",
                "email": "bmsstu@mail.com",
                "number": "12345678911",
                "extra_number": "12345678910",
                "region": "Moscow",
                "empl_num": 1830,
                "vacancies": {"id1", "id2"}  //array of strings
            }
        }

    Отрицательный

        HTTP/1.1 400


### "/resumes"
GET - запрос списка резюме(по фильтрам)


Возможные параметры

"region"
"wage_from"
"wage_to"
"experience"
"type_of_employment"
"work_schedule"


Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            {
                "id": uuid
                "first_name": "Vova",
                "second_name": "Zyablikov",
                "region": "Moscow",
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

            {
                "id": uuid
                "first_name": "Vova",
                "second_name": "Zyablikov",
                "region": "Moscow",
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

Возможные параметры

"region"
"wage"
"experience"
"type_of_employment"
"work_schedule"

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            {   
                "id": uuid
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
            
            {
                "id": uuid
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


	router.HandleFunc("/responds", h.GetResponds).Methods(http.MethodGet, http.MethodOptions)

### "/responds"
GET - запрос списка вакансий(по фильтрам)


Параметры

    vacancy_id 
        Возвращаются responds для конкретной вакансии
        Запрос доступен для только для владельца вакансии

    resume_id 
        Возвращаются responds для конкретного резюме
        Запрос доступен для только для владельца резюме
    
    При отсутствии параметров вернутся все responds на 
    вакансии или резюме пользователя

    Возможны три состояния respond
        1. awaits
        2. accepted
        3. rejected

Ответ на запрос:

    Положительный

        HTTP/1.1 200 OK

        {
            {
                "status": awaits
                "resume_id": uuid1
                "vacancy_id": uuid2
            }
             {
                "status": accepted
                "resume_id": uuid3
                "vacancy_id": uuid4
            }
        }

    Отрицательный
    
        При наличии двух аргументов resume_id и vacancy_id 

            HTTP/1.1 400