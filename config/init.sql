CREATE DATABASE IF NOT EXISTS carbuyer_assistance DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE carbuyer_assistance;

CREATE TABLE IF NOT EXISTS users (
  user_id VARCHAR(32) PRIMARY KEY,
  password_hash VARCHAR(255) NOT NULL,
  is_admin TINYINT(1) NOT NULL DEFAULT 0,
  status ENUM('active','disabled') NOT NULL DEFAULT 'active',
  full_name VARCHAR(100) NULL,
  phone VARCHAR(20) NULL,
  budget_min INT NULL,
  budget_max INT NULL,
  preferred_type VARCHAR(100) NULL,
  last_login_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT chk_users_budget_range CHECK (
    budget_min IS NULL OR budget_max IS NULL OR budget_min <= budget_max
  ),
  UNIQUE KEY uk_users_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS consults (
  consult_id CHAR(36) PRIMARY KEY,
  user_id VARCHAR(32) NULL,
  title VARCHAR(128),
  budget_range VARCHAR(255),
  preferred_type VARCHAR(100),
  use_case TEXT,
  fuel_type VARCHAR(50),
  brand_preference VARCHAR(255),
  llm_model VARCHAR(64),
  llm_prompt TEXT,
  llm_response TEXT,
  recommendations JSON,
  status ENUM('created','processing','completed','failed') NOT NULL DEFAULT 'completed',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_consults_user_id (user_id),
  INDEX idx_consults_created_at (created_at),
  CONSTRAINT fk_consult_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS feedbacks (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id VARCHAR(32) NOT NULL,
  consult_id CHAR(36) NULL,
  content TEXT NOT NULL,
  rating TINYINT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_feedback_user_id (user_id),
  INDEX idx_feedback_consult_id (consult_id),
  CONSTRAINT fk_feedback_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_feedback_consult FOREIGN KEY (consult_id) REFERENCES consults(consult_id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS gifts (
  gift_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description TEXT,
  points_cost INT NOT NULL,
  stock INT NOT NULL DEFAULT 0,
  status ENUM('active','inactive') NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_gifts_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS gift_redemptions (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  gift_id BIGINT NOT NULL,
  user_id VARCHAR(32) NOT NULL,
  quantity INT NOT NULL DEFAULT 1,
  points_spent INT NOT NULL,
  status ENUM('pending','completed','cancelled') NOT NULL DEFAULT 'completed',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_redemptions_user (user_id),
  INDEX idx_redemptions_gift (gift_id),
  CONSTRAINT fk_redemption_gift FOREIGN KEY (gift_id) REFERENCES gifts(gift_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_redemption_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS score_transactions (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id VARCHAR(32) NOT NULL,
  amount INT NOT NULL,
  type ENUM('consult','feedback','gift_redeem','manual_adjust','login','register') NOT NULL,
  ref_id VARCHAR(64) NULL,
  description VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_transactions_user (user_id),
  INDEX idx_transactions_type (type),
  CONSTRAINT fk_transaction_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE OR REPLACE VIEW v_user_points AS
SELECT user_id, COALESCE(SUM(amount), 0) AS points_balance
FROM score_transactions
GROUP BY user_id;

DROP PROCEDURE IF EXISTS sp_redeem_gift;
DELIMITER $$
CREATE PROCEDURE sp_redeem_gift(IN p_user_id VARCHAR(32), IN p_gift_id BIGINT, IN p_quantity INT)
BEGIN
  DECLARE v_cost INT;
  DECLARE v_stock INT;
  DECLARE v_points INT;

  IF p_quantity IS NULL OR p_quantity <= 0 THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid quantity';
  END IF;

  START TRANSACTION;

    SELECT points_cost, stock INTO v_cost, v_stock
    FROM gifts
    WHERE gift_id = p_gift_id AND status = 'active'
    FOR UPDATE;

    IF v_cost IS NULL THEN
      ROLLBACK;
      SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Gift not found or inactive';
    END IF;

    IF v_stock < p_quantity THEN
      ROLLBACK;
      SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient gift stock';
    END IF;

    SELECT COALESCE(SUM(amount), 0) INTO v_points
    FROM score_transactions
    WHERE user_id = p_user_id
    FOR UPDATE;

    IF v_points < (v_cost * p_quantity) THEN
      ROLLBACK;
      SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient user points';
    END IF;

    UPDATE gifts SET stock = stock - p_quantity WHERE gift_id = p_gift_id;

    INSERT INTO gift_redemptions (gift_id, user_id, quantity, points_spent, status)
    VALUES (p_gift_id, p_user_id, p_quantity, v_cost * p_quantity, 'completed');

    INSERT INTO score_transactions (user_id, amount, type, ref_id, description)
    VALUES (p_user_id, - (v_cost * p_quantity), 'gift_redeem', CAST(p_gift_id AS CHAR), 'Redeem gift');

  COMMIT;
END$$
DELIMITER ;

INSERT INTO users (user_id, password_hash, is_admin, status)
VALUES ('admin', '$2y$10$O0W7zWwDapRrzV0J9GS9LuZfxwPjLQyWKB1pXrS7JhVxZCU.ZJ5rS', 1, 'active')
ON DUPLICATE KEY UPDATE user_id = user_id;

INSERT INTO gifts (name, description, points_cost, stock, status)
VALUES
  ('车载充电器', '双口快充车载充电器', 200, 100, 'active'),
  ('后备箱收纳箱', '折叠可收纳后备箱整理箱', 500, 50, 'active'),
  ('车载香薰', '长效清香车载香薰', 150, 200, 'active')
ON DUPLICATE KEY UPDATE name = VALUES(name);