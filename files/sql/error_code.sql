
CREATE TABLE `tberrors` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `error_code` varchar(20) NOT NULL,
  `status` int(12) NOT NULL,
  `message` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE (`error_code`),
  PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8;