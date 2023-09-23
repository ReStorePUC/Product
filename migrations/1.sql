USE productdb;

CREATE TABLE images (
    id INT(6) AUTO_INCREMENT PRIMARY KEY,
    image_path varchar(100),
    product_id INT(6),
    FOREIGN KEY (product_id) REFERENCES products(id)
);
