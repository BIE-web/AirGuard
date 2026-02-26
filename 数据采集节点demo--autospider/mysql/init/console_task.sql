-- MySQL dump 10.13  Distrib 8.0.28, for Linux (x86_64)
--
-- Host: localhost    Database: console
-- ------------------------------------------------------
-- Server version	8.0.28

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
-- Table structure for table `task`
--

DROP TABLE IF EXISTS `task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `level` tinyint NOT NULL DEFAULT '1',
  `dependency_task_id` varchar(64) NOT NULL DEFAULT '',
  `dependency_status` tinyint NOT NULL DEFAULT '1',
  `spec` varchar(64) NOT NULL,
  `protocol` tinyint NOT NULL,
  `command` varchar(256) NOT NULL,
  `http_method` tinyint NOT NULL DEFAULT '1',
  `timeout` mediumint NOT NULL DEFAULT '0',
  `multi` tinyint NOT NULL DEFAULT '1',
  `retry_times` tinyint NOT NULL DEFAULT '0',
  `retry_interval` smallint NOT NULL DEFAULT '0',
  `notify_status` tinyint NOT NULL DEFAULT '1',
  `notify_type` tinyint NOT NULL DEFAULT '0',
  `notify_receiver_id` varchar(256) NOT NULL DEFAULT '',
  `notify_keyword` varchar(128) NOT NULL DEFAULT '',
  `tag` varchar(32) NOT NULL DEFAULT '',
  `remark` varchar(100) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT '0',
  `created` datetime NOT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_task_protocol` (`protocol`),
  KEY `IDX_task_status` (`status`),
  KEY `IDX_task_level` (`level`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task`
--

LOCK TABLES `task` WRITE;
/*!40000 ALTER TABLE `task` DISABLE KEYS */;
INSERT INTO `task` VALUES (1,'vpn pcap amazon',1,'',1,'0 0 2 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'not (src and dst net 172.19.0.0/24)\' -w /pcap/`date +%Y%m%d%H`_vpn_amazon.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:12:54',NULL),(2,'bro pcap amazon',1,'',1,'0 0 2 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'host 172.19.0.2 and port 11156\' -w /pcap/`date +%Y%m%d%H`_bro_amazon.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:14:09',NULL),(3,'spider amazon',1,'',1,'0 0 2 * * *',2,'cd /code && python3 spider.py amazon',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:15:02',NULL),(4,'vpn pcap reddit',1,'',1,'0 0 3 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'not (src and dst net 172.19.0.0/24)\' -w /pcap/`date +%Y%m%d%H`_vpn_reddit.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:12:54',NULL),(5,'bro pcap reddit',1,'',1,'0 0 3 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'host 172.19.0.2 and port 11156\' -w /pcap/`date +%Y%m%d%H`_bro_reddit.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:14:09',NULL),(6,'spider reddit',1,'',1,'0 0 3 * * *',2,'cd /code && python3 spider.py reddit',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:15:02',NULL),(7,'vpn pcap wiki',1,'',1,'0 0 4 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'not (src and dst net 172.19.0.0/24)\' -w /pcap/`date +%Y%m%d%H`_vpn_wiki.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:12:54',NULL),(8,'bro pcap wiki',1,'',1,'0 0 4 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'host 172.19.0.2 and port 11156\' -w /pcap/`date +%Y%m%d%H`_bro_wiki.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:14:09',NULL),(9,'spider wiki',1,'',1,'0 0 4 * * *',2,'cd /code && python3 spider.py wiki',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:15:02',NULL),(10,'vpn pcap yahoo',1,'',1,'0 0 5 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'not (src and dst net 172.19.0.0/24)\' -w /pcap/`date +%Y%m%d%H`_vpn_yahoo.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:12:54',NULL),(11,'bro pcap yahoo',1,'',1,'0 0 5 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'host 172.19.0.2 and port 11156\' -w /pcap/`date +%Y%m%d%H`_bro_yahoo.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:14:09',NULL),(12,'spider yahoo',1,'',1,'0 0 5 * * *',2,'cd /code && python3 spider.py yahoo',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:15:02',NULL),(13,'vpn pcap youtube',1,'',1,'0 0 6 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'not (src and dst net 172.19.0.0/24)\' -w /pcap/`date +%Y%m%d%H`_vpn_youtube.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:12:54',NULL),(14,'bro pcap youtube',1,'',1,'0 0 6 * * *',2,'timeout 1860 tcpdump -nnn -s0 -W 1 -G 1800 -i eth0 \'host 172.19.0.2 and port 11156\' -w /pcap/`date +%Y%m%d%H`_bro_youtube.pcap',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:14:09',NULL),(15,'spider youtube',1,'',1,'0 0 6 * * *',2,'cd /code && python3 spider.py youtube',1,0,0,3,10,0,1,'','','','',1,'2022-04-15 15:15:02',NULL);
/*!40000 ALTER TABLE `task` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-04-22  9:54:50
