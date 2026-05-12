-- MySQL dump 10.13  Distrib 8.0.19, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: testdb
-- ------------------------------------------------------
-- Server version	8.4.9

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `slug` varchar(100) NOT NULL,
  `type` enum('product','news') NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug` (`slug`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'Laptop','laptop','product','2026-05-03 04:35:06'),(2,'Smartphone','smartphone','product','2026-05-03 04:35:06'),(3,'Tablet','tablet','product','2026-05-03 04:35:06'),(4,'Aksesoris','aksesoris','product','2026-05-03 04:35:06'),(5,'Gaming','gaming','product','2026-05-03 04:35:06'),(6,'AI','ai','news','2026-05-03 04:35:06'),(7,'Teknologi','teknologi','news','2026-05-03 04:35:06'),(8,'Startup','startup','news','2026-05-03 04:35:06'),(9,'Gadget','gadget','news','2026-05-03 04:35:06'),(10,'Security','security','news','2026-05-03 04:35:06'),(11,'Audio','audio','product','2026-05-03 04:35:06'),(12,'Wearable','wearable','product','2026-05-03 04:35:06'),(13,'Networking','networking','product','2026-05-03 04:35:06'),(14,'Cloud','cloud','news','2026-05-03 04:35:06'),(15,'Programming','programming','news','2026-05-03 04:35:06'),(16,'Monitor','monitor','product','2026-05-03 04:35:06'),(17,'Storage','storage','product','2026-05-03 04:35:06'),(18,'Camera','camera','product','2026-05-03 04:35:06'),(19,'Automotive Tech','automotive-tech','news','2026-05-03 04:35:06'),(20,'Blockchain','blockchain','news','2026-05-03 04:35:06');
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `news`
--

DROP TABLE IF EXISTS `news`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `news` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `category_id` bigint unsigned NOT NULL,
  `title` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL,
  `content` longtext,
  `author` varchar(100) NOT NULL,
  `status` enum('draft','published') DEFAULT 'draft',
  `published_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug` (`slug`),
  KEY `fk_news_category` (`category_id`),
  KEY `idx_title` (`title`),
  KEY `idx_status_published` (`status`,`published_at`),
  CONSTRAINT `fk_news_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `news`
--

LOCK TABLES `news` WRITE;
/*!40000 ALTER TABLE `news` DISABLE KEYS */;
INSERT INTO `news` VALUES (1,6,'AI Mengubah Dunia','ai-mengubah-dunia','AI berkembang pesat','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(2,7,'Teknologi 5G Global','teknologi-5g','5G semakin luas','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(3,8,'Startup Unicorn Baru','startup-unicorn','Startup mencapai valuasi tinggi','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(4,9,'Gadget Terbaru 2026','gadget-2026','Banyak inovasi gadget','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(5,10,'Cyber Security Trend','cyber-security','Ancaman meningkat','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(6,6,'AI di Industri','ai-industri','Implementasi AI','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(7,7,'Cloud Computing','cloud-computing','Cloud makin populer','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(8,8,'Pendanaan Startup','funding-startup','Investor aktif','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(9,9,'Review Smartphone','review-smartphone','Perbandingan flagship','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(10,10,'Data Breach Besar','data-breach','Kebocoran data global','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(11,14,'AWS vs GCP','aws-vs-gcp','Perbandingan cloud','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(12,15,'Belajar Golang','belajar-golang','Tips coding','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(13,6,'AI Chatbot','ai-chatbot','Chatbot semakin pintar','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(14,7,'Quantum Computing','quantum','Masa depan computing','Admin','draft',NULL,'2026-05-03 04:35:27','2026-05-03 04:35:27'),(15,8,'IPO Startup','ipo-startup','Perusahaan go public','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(16,9,'Laptop Terbaru','laptop-terbaru','Laptop 2026','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(17,10,'Phishing Attack','phishing','Ancaman email','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(18,14,'Cloud Security','cloud-security','Keamanan cloud','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(19,19,'Mobil Listrik','mobil-listrik','EV berkembang','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27'),(20,20,'Blockchain Future','blockchain-future','Crypto dan blockchain','Admin','published','2026-05-03 04:35:27','2026-05-03 04:35:27','2026-05-03 04:35:27');
/*!40000 ALTER TABLE `news` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `outbox_events`
--

DROP TABLE IF EXISTS `outbox_events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `outbox_events` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `aggregate_type` enum('product','news') NOT NULL,
  `aggregate_id` bigint unsigned NOT NULL,
  `event_type` varchar(50) NOT NULL,
  `payload` json NOT NULL,
  `status` enum('pending','sent') DEFAULT 'pending',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `sent_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `outbox_events`
--

LOCK TABLES `outbox_events` WRITE;
/*!40000 ALTER TABLE `outbox_events` DISABLE KEYS */;
INSERT INTO `outbox_events` VALUES (1,'product',1,'created','{\"title\": \"Laptop ASUS ROG\"}','pending','2026-05-03 04:35:36',NULL),(2,'product',2,'created','{\"title\": \"MacBook Pro\"}','pending','2026-05-03 04:35:36',NULL),(3,'product',3,'created','{\"title\": \"iPhone 15\"}','pending','2026-05-03 04:35:36',NULL),(4,'product',4,'updated','{\"price\": 18000000}','pending','2026-05-03 04:35:36',NULL),(5,'product',5,'created','{\"title\": \"iPad Air\"}','pending','2026-05-03 04:35:36',NULL),(6,'news',1,'published','{\"title\": \"AI Dunia\"}','pending','2026-05-03 04:35:36',NULL),(7,'news',2,'published','{\"title\": \"5G\"}','pending','2026-05-03 04:35:36',NULL),(8,'news',3,'created','{\"title\": \"Startup\"}','pending','2026-05-03 04:35:36',NULL),(9,'news',4,'updated','{\"title\": \"Gadget\"}','pending','2026-05-03 04:35:36',NULL),(10,'news',5,'published','{\"title\": \"Security\"}','pending','2026-05-03 04:35:36',NULL),(11,'product',6,'created','{\"title\": \"Mouse\"}','pending','2026-05-03 04:35:36',NULL),(12,'product',7,'created','{\"title\": \"PS5\"}','pending','2026-05-03 04:35:36',NULL),(13,'product',8,'created','{\"title\": \"Sony XM5\"}','pending','2026-05-03 04:35:36',NULL),(14,'product',9,'updated','{\"price\": 8000000}','pending','2026-05-03 04:35:36',NULL),(15,'product',10,'created','{\"title\": \"Router\"}','pending','2026-05-03 04:35:36',NULL),(16,'news',6,'published','{\"title\": \"AI Industri\"}','pending','2026-05-03 04:35:36',NULL),(17,'news',7,'published','{\"title\": \"Cloud\"}','pending','2026-05-03 04:35:36',NULL),(18,'news',8,'created','{\"title\": \"Funding\"}','pending','2026-05-03 04:35:36',NULL),(19,'news',9,'updated','{\"title\": \"Review\"}','pending','2026-05-03 04:35:36',NULL),(20,'news',10,'published','{\"title\": \"Data Breach\"}','pending','2026-05-03 04:35:36',NULL);
/*!40000 ALTER TABLE `outbox_events` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `category_id` bigint unsigned NOT NULL,
  `title` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL,
  `description` text,
  `price` decimal(15,2) NOT NULL DEFAULT '0.00',
  `status` enum('active','inactive') DEFAULT 'active',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug` (`slug`),
  KEY `idx_title` (`title`),
  KEY `idx_category` (`category_id`),
  CONSTRAINT `fk_products_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
INSERT INTO `products` VALUES (1,1,'Laptop ASUS ROG','laptop-asus-rog','Gaming laptop high performance',25000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(2,1,'Laptop MacBook Pro','macbook-pro','Apple laptop M3 chip',32000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(3,2,'iPhone 15 Pro','iphone-15-pro','Flagship smartphone Apple',21000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(4,2,'Samsung Galaxy S24','galaxy-s24','Android flagship Samsung',18000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(5,3,'iPad Air','ipad-air','Tablet Apple ringan',12000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(6,4,'Mouse Logitech G Pro','logitech-g-pro','Mouse gaming ringan',1500000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(7,5,'PS5 Console','ps5-console','Next gen console',9000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(8,11,'Sony WH-1000XM5','sony-xm5','Noise cancelling headphone',5000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(9,12,'Apple Watch Series 9','apple-watch-s9','Smartwatch Apple',8000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(10,13,'TP-Link AX3000','tplink-ax3000','Router wifi 6',1200000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(11,16,'LG Ultrawide 34\"','lg-ultrawide','Monitor ultrawide',7000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(12,17,'Samsung SSD 1TB','samsung-ssd-1tb','SSD NVMe cepat',1500000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(13,18,'Sony A7 IV','sony-a7iv','Mirrorless camera',35000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(14,1,'Laptop Lenovo Legion','lenovo-legion','Gaming laptop Lenovo',22000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(15,2,'Xiaomi 14','xiaomi-14','Flagship Xiaomi',12000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(16,3,'Samsung Tab S9','tab-s9','Tablet Samsung',14000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(17,4,'Mechanical Keyboard','mech-keyboard','Keyboard RGB',1200000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(18,5,'Xbox Series X','xbox-series-x','Console Microsoft',8500000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(19,11,'JBL Charge 5','jbl-charge-5','Speaker portable',2000000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18'),(20,17,'WD HDD 2TB','wd-hdd-2tb','Harddisk eksternal',900000.00,'active','2026-05-03 04:35:18','2026-05-03 04:35:18');
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'testdb'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-05-03 11:48:17
