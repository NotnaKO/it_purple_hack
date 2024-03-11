## Создание бд
Запустить run.sh в текущей директории там создается БД:\
\
Название БД: postgres\
User: postgres\
Password: 1234\
Схема: matrix


## Таблицы:
1) baseline_matrix_1.sql - цена вида (x, y, y) - для тестов
2) baseline_matrix_2.sql - цены в локациях и категоряих рандомные
3) baseline_matrix_3.sql - храним цены откаты(хот тейбл мб будет)
4) discount_matrix_1.sql - используем для скидок основная по локациям
4) discount_matrix_2.sql - храним скидку для региона
4) discount_matrix_3.sql - хранится скидка для рута(например на всю страну)