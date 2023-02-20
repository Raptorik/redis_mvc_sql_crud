create table product_list
(
    Id          int unsigned auto_increment
        primary key,
    Name        varchar(20)  not null,
    Price       int          not null,
    Quantity    varchar(50)  null,
    Description varchar(200) null,
    Action      varchar(50)  null
);

INSERT INTO product_list (Id, Name, Price, Quantity, Description, Action) VALUES (1, 'computer', 15, '1', 'good', 'Edit | Delete');
INSERT INTO product_list (Id, Name, Price, Quantity, Description, Action) VALUES (3, 'jewelry', 999, '1', 'excellent', 'Edit | Delete');
INSERT INTO product_list (Id, Name, Price, Quantity, Description, Action) VALUES (4, 'car', 9999, '2', 'bad', 'Edit | Delete');
