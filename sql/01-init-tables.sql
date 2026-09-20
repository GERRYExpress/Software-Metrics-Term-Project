DROP DATABASE IF EXISTS mymariadb;

CREATE DATABASE IF NOT EXISTS mymariadb;
USE mymariadb;

CREATE TABLE products (
	product_id INT NOT NULL UNIQUE AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(1024) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    poster_url VARCHAR(255) NOT NULL,
    PRIMARY KEY (product_id)
) ENGINE=InnoDB;

CREATE TABLE categories (
	category_id INT NOT NULL UNIQUE AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    is_hidden BOOLEAN NOT NULL DEFAULT 0,
    PRIMARY KEY (category_id)
) ENGINE=InnoDB;

CREATE TABLE accounts (
	user_id INT NOT NULL UNIQUE AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_date DATE NOT NULL,
    is_buyer BOOLEAN DEFAULT FALSE,
    total_spend DECIMAL(10,2) DEFAULT 0,
    total_game INT DEFAULT 0,
    is_seller BOOLEAN DEFAULT FALSE,
    total_sale DECIMAL(10,2) DEFAULT 0,
    PRIMARY KEY (user_id)
) ENGINE=InnoDB;

CREATE TABLE hardware_requirement (
	requirement_id INT NOT NULL UNIQUE AUTO_INCREMENT,
    operating_system VARCHAR(50) DEFAULT '-',
    processor VARCHAR(50) DEFAULT '-',
    memory VARCHAR(50) DEFAULT '-',
    graphics_card VARCHAR(50) DEFAULT '-',
    directx_version VARCHAR(50) DEFAULT '-',
    network VARCHAR(50) DEFAULT '-',
    storage VARCHAR(50) DEFAULT '-',
    soundcard VARCHAR(50) DEFAULT '-',
    additional_note VARCHAR(255) DEFAULT '-',
    type ENUM('min','recommended'),
    product_id INT NOT NULL,
    PRIMARY KEY (requirement_id),
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE images (
	image_id INT NOT NULL UNIQUE AUTO_INCREMENT,
    image_url VARCHAR(255),
    product_id INT NOT NULL,
    PRIMARY KEY (image_id),
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
) ENGINE=InnoDB;
CREATE TABLE videos (
	video_id INT NOT NULL UNIQUE AUTO_INCREMENT,
    video_url VARCHAR(255),
    product_id INT NOT NULL,
    PRIMARY KEY (video_id),
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE orders (
	order_id INT NOT NULL UNIQUE AUTO_INCREMENT,
    user_id INT NOT NULL,
    order_date DATETIME NOT NULL,
    payment_method ENUM('promptpay', 'mastercard', 'visa', 'truemoney') NOT NULL,
    final_price FLOAT NOT NULL,
    PRIMARY KEY (order_id),
    FOREIGN KEY (user_id) REFERENCES accounts(user_id)  ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE order_items (
	order_id INT NOT NULL,
    product_id INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(order_id)  ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id)  ON DELETE RESTRICT,
    PRIMARY KEY (order_id, product_id)
) ENGINE=InnoDB;

CREATE TABLE wishlists (
	user_id INT NOT NULL,
    product_id INT NOT NULL,
    wishlist_date DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES accounts(user_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE RESTRICT,
    PRIMARY KEY (user_id, product_id)
) ENGINE=InnoDB;

CREATE TABLE carts (
	user_id INT NOT NULL,
    product_id INT NOT NULL,
    cart_date DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES accounts(user_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, product_id)
) ENGINE=InnoDB;

CREATE TABLE product_lists (
	list_id INT UNIQUE NOT NULL AUTO_INCREMENT,
	product_id INT NOT NULL,
    category_id INT DEFAULT NULL,
    user_id INT NOT NULL,
    release_date DATE NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES accounts(user_id) ON DELETE RESTRICT,
    PRIMARY KEY (list_id)
)ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `sessions` (
  `session_id` varchar(128) COLLATE utf8mb4_bin NOT NULL,
  `expires` int(11) unsigned NOT NULL,
  `data` mediumtext COLLATE utf8mb4_bin,
  PRIMARY KEY (`session_id`)
) ENGINE=InnoDB
