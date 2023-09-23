CREATE DATABASE IF NOT EXISTS productdb;

USE productdb;

CREATE TABLE products (
    id INT(6) AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50),
    description VARCHAR(100),
    categories VARCHAR(100),
    size VARCHAR(50),
    price FLOAT,
    tax FLOAT,
    available BOOLEAN,
    store_id INT
)