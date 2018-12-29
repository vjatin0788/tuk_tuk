CREATE TABLE `tb_ride_details` (
  `id` bigint(12) unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` bigint(12) unsigned NOT NULL ,
  `driver_id` bigint(12) unsigned NOT NULL ,
  `source_lat` DOUBLE NOT NULL DEFAULT 0,
  `source_long` DOUBLE NOT NULL DEFAULT 0,
  `destination_lat` DOUBLE NOT NULL DEFAULT 0,
  `destination_long` DOUBLE NOT NULL DEFAULT 0,
  `status` tinyint(4) unsigned NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=1;