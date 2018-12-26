
CREATE TABLE `driver_tracking` (
  `id` bigint(12) unsigned NOT NULL AUTO_INCREMENT,
  `driver_id` bigint(12) unsigned NOT NULL ,
  `current_lat` DOUBLE NOT NULL DEFAULT 0,
  `current_long` DOUBLE NOT NULL DEFAULT 0,
  `last_lat` DOUBLE NOT NULL DEFAULT 0,
  `last_long` DOUBLE NOT NULL DEFAULT 0,
  `current_lat_rad` DOUBLE NOT NULL DEFAULT 0,
  `current_long_rad` DOUBLE NOT NULL DEFAULT 0,
  `last_lat_rad` DOUBLE NOT NULL DEFAULT 0,
  `last_long_rad` DOUBLE NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=1;



SELECT * FROM driver_tracking WHERE acos(sin(user_lat) * sin(current_lat) + cos(user_lat) * cos(current_lat) * cos(current_long - (user_long))) * 6371 <= 1000

INSERT INTO `driver_tracking`(`driver_id`,`current_lat`,`current_long`,`last_lat`,`last_long`) VALUES(?,?,?,?,?)

INSERT INTO `driver_tracking`(`driver_id`,`current_lat`,`current_long`,`last_lat`,`last_long`) VALUES(1,12.1222,11.232323,12.1111,122.1221212)

UPDATE `driver_tracking` SET `current_lat`=?,`current_long`=?,`last_lat`=?,`last_long`=? WHERE  `driver_id` = ?