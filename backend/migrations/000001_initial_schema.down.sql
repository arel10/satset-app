-- SatSet Laundry Initial Schema Rollback

DROP TABLE IF EXISTS `idempotency_keys`;
DROP TABLE IF EXISTS `audit_logs`;
DROP TABLE IF EXISTS `ai_messages`;
DROP TABLE IF EXISTS `ai_conversations`;
DROP TABLE IF EXISTS `expenses`;
DROP TABLE IF EXISTS `expense_categories`;
DROP TABLE IF EXISTS `payments`;
DROP TABLE IF EXISTS `order_sequences`;
DROP TABLE IF EXISTS `order_status_history`;
DROP TABLE IF EXISTS `order_items`;
DROP TABLE IF EXISTS `orders`;
DROP TABLE IF EXISTS `services`;
DROP TABLE IF EXISTS `customers`;
DROP TABLE IF EXISTS `settings`;
DROP TABLE IF EXISTS `refresh_tokens`;
DROP TABLE IF EXISTS `users`;
