CREATE TABLE `Users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(100) NOT NULL,
  `role` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `Users` (`username`, `password`, `email`, `role`) 
VALUES ('beanAdmin', '$2a$10$xgd7lVHK77P9kZinKW52IuByK30d4FsQpouaW5EfvCGOKZfAXKcQe', 'bean@admin.com', 'admin');

SELECT * From Users;