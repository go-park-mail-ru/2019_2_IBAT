DROP TABLE companies;
DROP TABLE offers;
DROP TABLE responds;
DROP TABLE resumes;
DROP TABLE vacancies;
DROP TABLE persons;
DROP TABLE favorites_resumes;
DROP TABLE favorites_vacancies;

CREATE TABLE persons(
    id uuid PRIMARY KEY,
    first_name VARCHAR (50) NOT NULL,
    second_name VARCHAR (70) NOT NULL,
    email VARCHAR (355) UNIQUE NOT NULL,
    password_hash VARCHAR (70) NOT NULL,
    role VARCHAR (50) NOT NULL,
    path_to_image VARCHAR (50)
);

CREATE TABLE companies(
    company_name VARCHAR (70) PRIMARY KEY,
    own_id uuid REFERENCES persons (id) ON DELETE CASCADE UNIQUE NOT NULL,
    site VARCHAR (70) NOT NULL,
    --not not null
    region  VARCHAR (70),
    phone_number VARCHAR (70),
    extra_phone_number VARCHAR (70),
    
    spheres_of_work TEXT,
    -- must be done by tags
    empl_num VARCHAR(70),
    description TEXT
);

CREATE TABLE resumes(
    id uuid PRIMARY KEY,
    own_id uuid REFERENCES persons (id) ON DELETE CASCADE NOT NULL,
    email VARCHAR (355) NOT NULL,
    region  VARCHAR (70) NOT NULL,
    phone_number VARCHAR (30) NOT NULL,
    first_name VARCHAR (50) NOT NULL,
    second_name VARCHAR (70) NOT NULL,
    birth_date   DATE,
    sex VARCHAR (30) NOT NULL,

    type_of_employment VARCHAR (50),
	work_schedule VARCHAR (50),
    citizenship VARCHAR (70),
    profession VARCHAR (70),
    -- must be done by tags
    position    VARCHAR (70),
    experience  TEXT,
    education TEXT,
    wage  MONEY,
    about  TEXT
);

CREATE TABLE vacancies(
    id uuid PRIMARY KEY,
    own_id uuid REFERENCES persons (id) ON DELETE CASCADE NOT NULL,

    region     VARCHAR (70),
    profession VARCHAR (70) NOT NULL,
    position    VARCHAR (70),
    experience  TEXT,
    wage_from  MONEY,
    wage_to    MONEY,
    tasks TEXT,
    type_of_employment VARCHAR (70),
    work_schedule      VARCHAR (70),
    requirements TEXT,
    conditions TEXT,
    about  TEXT
);


CREATE TABLE offers(
    status VARCHAR (70) NOT NULL,
    resume_id uuid REFERENCES resumes (id) ON DELETE CASCADE NOT NULL ,
    vacancy_id uuid REFERENCES vacancies (id) ON DELETE CASCADE NOT NULL ,
    PRIMARY KEY(resume_id, vacancy_id)
);

CREATE TABLE responds(
    status VARCHAR (70) NOT NULL,
    resume_id uuid REFERENCES resumes (id) ON DELETE CASCADE NOT NULL,
    vacancy_id uuid REFERENCES vacancies (id) ON DELETE CASCADE NOT NULL ,
    PRIMARY KEY(resume_id, vacancy_id)
);

-- CREATE TABLE favorites_resumes(
--     status VARCHAR (70) NOT NULL,
--     employer_id uuid REFERENCES persons (id) ON DELETE CASCADE NOT NULL,
--     resume_id uuid REFERENCES vacancies (id) ON DELETE CASCADE NOT NULL ,
--     PRIMARY KEY(resume_id, vacancy_id)
-- );

CREATE TABLE favorite_vacancies(
    person_id uuid REFERENCES persons (id) ON DELETE CASCADE NOT NULL,
    vacancy_id uuid REFERENCES vacancies (id) ON DELETE CASCADE NOT NULL ,
    PRIMARY KEY(person_id, vacancy_id)
);


INSERT INTO persons(id, email, first_name, second_name, password_hash, role)
VALUES(gen_random_uuid(), 'vladle@mail.ru', 'Vlad', 'Lee', 'iearoiqwejfka', 'seeker');

INSERT INTO persons(id, email, first_name, second_name, password_hash, role)
VALUES(gen_random_uuid(), 'vasyapupkin@mail.ru', 'Vasya', 'Pupkin', 'iearoiqdsfwejfka', 'seeker');


INSERT INTO persons(id, email, first_name, second_name, password_hash, role)
VALUES(gen_random_uuid(), 'main@mail.ru', 'Vladimir', 'Lenin', 'iearoiqdsfwejfka', 'employer');

INSERT INTO companies(own_id, company_name, site, region, spheres_of_work, empl_num,
description, phone_number, extra_phone_number)
VALUES((SELECT id FROM persons WHERE email = 'main@mail.ru'), 'Mail.ru', 'Mail.ru',
'Moscow', 'IT, business', 'more than 1000', 'best company that ever existed', '89266239478', '8926639479');

INSERT INTO persons(id, email, first_name, second_name, password_hash, role)
VALUES(gen_random_uuid(), 'yandex@mail.ru', 'Sasha', 'Koen', 'iearoiqdsfwejfka', 'employer');

INSERT INTO companies(own_id, company_name, site, spheres_of_work, empl_num,
description, phone_number, extra_phone_number)
VALUES((SELECT id FROM persons WHERE email = 'yandex@mail.ru'), 'Yandex', 'Yandex.ru',
'IT, business', 'more than 1000', 'not so bad', '89266239479', '8926639479');

INSERT INTO vacancies(id, own_id, profession, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Mail.ru'), 'frontend developer', 'middle', 'Moscow', '4 years', 125000, 130000, 'Полная занятость',
'write frontend', 'JS', 'Полный день','nice office, good team', 'the best IT company');

INSERT INTO vacancies(id, own_id, profession, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Mail.ru'),
'backend developer', 'middle', 'Moscow', 'Более 6 лет', 80000, 120000, 'Полная занятость',
'write backend', 'Go', 'Полный день', 'nice office, good team', 'the best IT company');

INSERT INTO vacancies(id, own_id, profession, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Yandex'), 'backend developer', 'middle', 'Moscow', 'От 3 до 6 лет', 125000, 250000, 'Полная занятость',
'write backend', 'Go', 'Удаленная работа','nice office, good team', 'top 2 IT company');


INSERT INTO vacancies(id, own_id, profession, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Yandex'),
'data scientist', 'middle', 'Moscow', 'Более 6 лет', 150000, 300000, 'Полная занятость',
'write II', 'Python, math', 'Удаленная работа','nice office, good team', 'top 2 IT company');



INSERT INTO resumes(id, own_id, email,  first_name, second_name, region, phone_number, birth_date, sex,
citizenship, profession, position, education, wage, about, type_of_employment, work_schedule, experience)
VALUES(gen_random_uuid(), (SELECT id FROM persons WHERE email = 'vladle@mail.ru'),
'vladle@mail.ru', 'Vlad', 'Lee', 'Москва', '89266211479', '1991-10-10', 'мужской', 'русское',
'фронтенд-разработчик', 'middle', 'МГТУ', '123000', 'Хороший парень', 'Полная занятость', 'Удаленная работа',
'Более 6 лет');

INSERT INTO resumes(id, own_id, email,  first_name, second_name, region, phone_number, birth_date, sex,
citizenship, profession, position, education, wage, about, type_of_employment, work_schedule, experience)
VALUES(gen_random_uuid(), (SELECT id FROM persons WHERE email = 'vladle@mail.ru'),
'vladle@mail.ru', 'Vlad', 'Lee', 'Москва', '89266211479', '1991-10-10', 'мужской', 'русское',
'бэкенд-разработчик', 'middle', 'МГТУ', '123000', 'Хороший парень', 'Полная занятость', 'Удаленная работа',
'От 3 до 6 лет');
   
INSERT INTO resumes(id, own_id, email,  first_name, second_name, region, phone_number, birth_date, sex,
citizenship, profession, position, education, wage, about, type_of_employment, work_schedule, experience)
VALUES(gen_random_uuid(), (SELECT id FROM persons WHERE email = 'vladle@mail.ru'),
'vasyapupkin@mail.ru', 'Вася', 'Пупкин', 'Москва', '89236211479', '1998-10-10', 'мужской', 'русское',
'дата-саентист', 'middle', 'МГУ', '170000', 'Хороший парень', 'Полная занятость', 'Полный день',
'Более 6 лет');
