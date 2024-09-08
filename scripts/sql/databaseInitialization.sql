CREATE USER 'shopping-cart-app'@'%' IDENTIFIED BY 'shoppingCartPassword!'; 

CREATE DATABASE shoppingCart;

GRANT SELECT, INSERT, UPDATE, DELETE ON shoppingCart.* TO 'shopping-cart-app'@'%';

CREATE TABLE `shoppingCart`.`cartItem` (
  -- Id could be a uuid. Being an int leaves the value exposed to potential attackers, that can know how many items are, at least,
  -- by looking at the id. A uuid makes this data more opaque. It is harder to index, though. 
  -- For simplicity, I´ve used an int. Name is a dangerous column for primary key as there could be several
  -- items with the same name. Perhaps in the future, let´s say there are 2 jackets but with different color.
  `id` int PRIMARY KEY,
  `name` varchar(50),
  `quantity` int, 
  `reservationId` varchar(50)
);