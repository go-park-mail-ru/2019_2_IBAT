DROP TABLE companies;
DROP TABLE offers;
DROP TABLE responds;
DROP TABLE resumes;
DROP TABLE favorite_vacancies;
DROP TABLE vacancies;
DROP TABLE persons;

-- DROP TABLE favorites_resumes;

CREATE TABLE persons(
    id uuid PRIMARY KEY,
    first_name VARCHAR (50) NOT NULL,
    second_name VARCHAR (70) NOT NULL,
    email VARCHAR (355) UNIQUE NOT NULL,
    password_hash VARCHAR (255) NOT NULL,
    role VARCHAR (50) NOT NULL,
    path_to_image VARCHAR (50) DEFAULT ''
);

CREATE TABLE companies(
    company_name VARCHAR (70) PRIMARY KEY,
    own_id uuid REFERENCES persons (id) ON DELETE CASCADE UNIQUE NOT NULL,
    site VARCHAR (70) NOT NULL,
    --not not null
    region  VARCHAR (70) DEFAULT '',
    phone_number VARCHAR (70) DEFAULT '',
    extra_phone_number VARCHAR (70) DEFAULT '',
    
    spheres_of_work TEXT DEFAULT '',
    -- must be done by tags
    empl_num VARCHAR(70) DEFAULT '',
    description TEXT DEFAULT ''
);

CREATE TABLE resumes(
    id uuid PRIMARY KEY,
    own_id uuid REFERENCES persons (id) ON DELETE CASCADE NOT NULL,
    email VARCHAR (355) NOT NULL,
    region  VARCHAR (70) NOT NULL,
    phone_number VARCHAR (30) NOT NULL,
    first_name VARCHAR (50) NOT NULL,
    second_name VARCHAR (70) NOT NULL,
    birth_date   DATE NOT NULL,
    sex VARCHAR (30) NOT NULL,

    type_of_employment VARCHAR (50) DEFAULT '',
	work_schedule VARCHAR (50) DEFAULT '',
    citizenship VARCHAR (70) DEFAULT '',
    -- must be done by tags
    position    VARCHAR (70) DEFAULT '',
    experience  TEXT DEFAULT '',
    education TEXT DEFAULT '',
    wage  MONEY DEFAULT '',
    about  TEXT DEFAULT ''
);

CREATE TABLE vacancies(
    id uuid PRIMARY KEY,
    own_id uuid REFERENCES persons (id) ON DELETE CASCADE NOT NULL,

    region     VARCHAR (70) DEFAULT '',
    position    VARCHAR (70),
    experience  TEXT DEFAULT '',
    wage_from  MONEY DEFAULT '',
    wage_to    MONEY DEFAULT '',
    tasks TEXT DEFAULT '',
    type_of_employment VARCHAR (70) DEFAULT '',
    work_schedule      VARCHAR (70) DEFAULT '',
    requirements TEXT DEFAULT '',
    conditions TEXT DEFAULT '',
    about  TEXT DEFAULT ''
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

INSERT INTO vacancies(id, own_id, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Mail.ru'), 'Москва', 'frontend developer',, '4 years', 125000, 130000, 'Полная занятость',
'write frontend', 'JS', 'Полный день','nice office, good team', 'the best IT company');

INSERT INTO vacancies(id, own_id, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Mail.ru'), 'Москва', 'backend developer', 'Более 6 лет', 80000, 120000, 'Полная занятость',
'write backend', 'Go', 'Полный день', 'nice office, good team', 'the best IT company');

INSERT INTO vacancies(id, own_id, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Yandex'), 'Москва', 'backend developer', 'От 3 до 6 лет', 125000, 250000, 'Полная занятость',
'write backend', 'Go', 'Удаленная работа','nice office, good team', 'top 2 IT company');


INSERT INTO vacancies(id, own_id, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Yandex'),
'Москва', 'data scientist', 'Более 6 лет', 150000, 300000, 'Полная занятость',
'write II', 'Python, math', 'Удаленная работа','nice office, good team', 'top 2 IT company');



INSERT INTO resumes(id, own_id, email,  first_name, second_name, region, phone_number, birth_date, sex,
citizenship, position, education, wage, about, type_of_employment, work_schedule, experience)
VALUES(gen_random_uuid(), (SELECT id FROM persons WHERE email = 'vladle@mail.ru'),
'vladle@mail.ru', 'Vlad', 'Lee', 'Москва', '89266211479', '1991-10-10', 'мужской', 'русское',
'фронтенд-разработчик','МГТУ', '123000', 'Хороший парень', 'Полная занятость', 'Удаленная работа',
'Более 6 лет');

INSERT INTO resumes(id, own_id, email,  first_name, second_name, region, phone_number, birth_date, sex,
citizenship, position, education, wage, about, type_of_employment, work_schedule, experience)
VALUES(gen_random_uuid(), (SELECT id FROM persons WHERE email = 'vladle@mail.ru'),
'vladle@mail.ru', 'Vlad', 'Lee', 'Москва', '89266211479', '1991-10-10', 'мужской', 'русское',
'бэкенд-разработчик', 'МГТУ', '123000', 'Хороший парень', 'Полная занятость', 'Удаленная работа',
'От 3 до 6 лет');
   
INSERT INTO resumes(id, own_id, email,  first_name, second_name, region, phone_number, birth_date, sex,
citizenship, position, education, wage, about, type_of_employment, work_schedule, experience)
VALUES(gen_random_uuid(), (SELECT id FROM persons WHERE email = 'vladle@mail.ru'),
'vasyapupkin@mail.ru', 'Вася', 'Пупкин', 'Москва', '89236211479', '1998-10-10', 'мужской', 'русское',
'дата-саентист', 'МГУ', '170000', 'Хороший парень', 'Полная занятость', 'Полный день',
'Более 6 лет');

INSERT INTO vacancies(id, own_id, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Yandex'), 
'Москва', 
'Frontend-разработчик', 
'Более 6 лет', 150000, 300000, 
'Полная занятость',
'· Программирование Web-интерфейсов с использованием указанных технологий,\n
· Участие в обсуждении требований к разрабатываемым системам', 
'· Уровень Middle;
· Отличные знания JS, HTML и CSS;
· Знание jQuery, ES, sass или less;
· Опыт использования систем сборки;
· Опыт использования любой из системы контроля версий;
· Инициативность, тяга к знаниям и ответственность;
· Умения писать качественный поддерживаемый код;
· Опыт работы в команде;', 
'Удаленная работа',
'· Работа в SCRUM команде над интересными и сложными задачами;
· Официальное оформление в полном соответствии с Трудовым Кодексом РФ;
· Удобный офис с отличным видом на Шлюзовую набережную и в шаговой доступности от м. Павелецкая;
· Гибкое начало и окончание рабочего дня;
· Заработная плата обсуждается индивидуально на собеседовании;
· Участие в конференциях и курсы повышения квалификации;
· Возможность бесплатно посещать развлекательные мероприятия;
·Собеседование проводится очно в офисе в Москве.', 
'top 2 IT company');

INSERT INTO vacancies(id, own_id, region, position, experience,
wage_from, wage_to, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Mail.ru'), 
'Москва', 
'Программист Python', 
'Более 6 лет', 100000, 200000, 
'Полная занятость',
'Разработка агрегатора для сервиса продажи билетов;
Разработка веб­-проектов компании;
Поддержка и доработка существующих проектов.', 

'Хорошее знание Python (3.х) и Django Framework;
Опыт работы с реляционными базами данных (PostgreSQL);
Опыт работы с GIT;
Опыт работы с Docker;
Бибилиотеки, которые мы используем Requests, lxml, API на DRF, различные django библиотеки и д.р.;
Понимание принципов построение API, протокола HTTP;
Умение написать правильный код, который не тормозит;
Опыт написания тестов;
Будет плюсом опыт работы с elasticsearch;
Если вы чего-то не знаете, при вашем желании вы быстро научитесь у нас.
У нас есть CI/CD, Git Flow, PR, и бекенд разработчик занимается именно написанием качественного python кода;
Мы ездим на pycon и другие конференции, оплачиваем их. 
И даже если у вас не очень много знганий, но очень много желания, 
то не стесняйтесь мы попробуем вам помочь.', 

'Удаленная работа',
'Оформление по ТК;
Заработная плата по результатам собеседования;
Удобное расположение офиса, 7 мин пешком от м. Трубная, Театральная, Чеховская.
Мы работаем в самом центре города. Офисная планировка;
8­-ми часовой рабочий день. Время начала работы выбираем сами;
Рассматриваем программиста преимущественно на удаленной основе, но можно и в офис;
Свободный стиль одежды;
Дружная команда.', 
'По успешному завершению проектов выплачиваются бонусы;
Пицца-дни. Яркие масштабные корпоративы и спортивные мероприятия / молчание, 
понимание и пледы для талантливых интровертов :).');


INSERT INTO vacancies(id, own_id, region, position, experience,
wage_from, type_of_employment, tasks, requirements, work_schedule,
conditions, about)VALUES(gen_random_uuid(),
(SELECT own_id FROM companies WHERE company_name = 'Mail.ru'), 
'Москва', 
'Back-end (PHP) разработчик Senior', 
'Более 6 лет', 200000, 
'Полная занятость',
'Проектировать и разрабатывать архитектуру продуктов компании в составе команды разработчиков;
Интеграция с внешними сервисами (API);
Анализировать и оптимизировать существующий кода.', 

'Знание принципов Solid;
Отличное знание и применение ООП
Стек: PHP 5.6, 7 (Symfony 2.8, 4.2); PostgreSQL, RabbitMQ, Doctrine, Redis, Git
Опыт работы с Git (сами работаем с GitLab)
Умение разбираться с чужим кодом и быстро включаться в процесс;
Работа исключительно в офисе', 

'Удаленная работа',
'Достойную "белую" заработную плату по результатам собеседования;
Работу в компании с вменяемым руководством и адекватным отношением к сотрудникам;
интересные задачи и очень большой проект с продвинутыми технологиями;
Посещение профильных конференций и курсов за счет компании;
Удобный график работы 5/2 время начала рабочего дня обсуждаемо',
 
'По успешному завершению проектов выплачиваются бонусы;
Пицца-дни. Яркие масштабные корпоративы и спортивные мероприятия / молчание, 
понимание и пледы для талантливых интровертов :).');
