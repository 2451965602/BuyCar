-- file: config/sql/init.sql
CREATE DATABASE IF NOT EXISTS carbuyer_assistance DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE carbuyer_assistance;

CREATE TABLE IF NOT EXISTS users (
                                     user_id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                     username VARCHAR(255) UNIQUE NOT NULL,
                                     name VARCHAR(255)  NOT NULL,
                                     telephone VARCHAR(255) UNIQUE NOT NULL,
                                     email VARCHAR(255) UNIQUE NOT NULL,
                                     password_hash VARCHAR(255) NOT NULL,
                                     is_admin bool NOT NULL DEFAULT false,
                                     score INT NOT NULL DEFAULT 0,
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS consults (
                                        consult_id BIGINT PRIMARY KEY,
                                        user_id BIGINT NULL,
                                        title VARCHAR(255),
                                        budget_range VARCHAR(255),
                                        preferred_type VARCHAR(255),
                                        use_case TEXT,
                                        fuel_type VARCHAR(255),
                                        brand_preference VARCHAR(255),
                                        llm_model VARCHAR(255),
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
                                         user_id BIGINT NOT NULL,
                                         consult_id BIGINT NULL,
                                         content TEXT NOT NULL,
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
                                     score_cost INT NOT NULL,
                                     stock INT NOT NULL DEFAULT 0,
                                     status ENUM('active','inactive') NOT NULL DEFAULT 'active',
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                     INDEX idx_gifts_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS gift_purchase (
                                             id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                             gift_id BIGINT NOT NULL,
                                             user_id BIGINT NOT NULL,
                                             number INT NOT NULL DEFAULT 1,
                                             score_cost INT NOT NULL,
                                             status ENUM('pending','completed','cancelled') NOT NULL DEFAULT 'completed',
                                             created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                             INDEX idx_redemptions_user (user_id),
                                             INDEX idx_redemptions_gift (gift_id),
                                             CONSTRAINT fk_redemption_gift FOREIGN KEY (gift_id) REFERENCES gifts(gift_id) ON DELETE RESTRICT ON UPDATE CASCADE,
                                             CONSTRAINT fk_redemption_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS score_transactions (
                                                  id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                                  user_id BIGINT NOT NULL,
                                                  change_amount INT NOT NULL,
                                                  reason ENUM('consult','feedback','gift_redeem','manual_adjust') NOT NULL,
                                                  ref_id BIGINT NULL,
                                                  description VARCHAR(255) NULL,
                                                  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                  INDEX idx_transactions_user (user_id),
                                                  INDEX idx_transactions_type (reason),
                                                  CONSTRAINT fk_transaction_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
