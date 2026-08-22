CREATE TABLE IF NOT EXISTS feedbacks (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  repair_order_id BIGINT UNSIGNED NOT NULL,
  scheduled_at DATETIME NOT NULL,
  completed_at DATETIME NULL,
  method VARCHAR(32) NOT NULL DEFAULT 'PHONE',
  complaint_type VARCHAR(32) NOT NULL DEFAULT '',
  complaint TEXT NOT NULL,
  resolution TEXT NOT NULL,
  rework_order_id BIGINT UNSIGNED NULL,
  rework_reason VARCHAR(255) NOT NULL DEFAULT '',
  rework_cost DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_feedbacks_repair_order_id (repair_order_id),
  KEY idx_feedbacks_status (status),
  KEY idx_feedbacks_deleted_at (deleted_at),
  CONSTRAINT fk_feedbacks_order FOREIGN KEY (repair_order_id) REFERENCES repair_orders(id) ON DELETE RESTRICT,
  CONSTRAINT fk_feedbacks_rework_order FOREIGN KEY (rework_order_id) REFERENCES repair_orders(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS feedback_scores (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  feedback_id BIGINT UNSIGNED NOT NULL,
  dimension VARCHAR(64) NOT NULL,
  score INT NOT NULL DEFAULT 1,
  comment VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  KEY idx_feedback_scores_feedback_id (feedback_id),
  CONSTRAINT fk_feedback_scores_feedback FOREIGN KEY (feedback_id) REFERENCES feedbacks(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS feedback_status_histories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  feedback_id BIGINT UNSIGNED NOT NULL,
  from_status VARCHAR(32) NOT NULL,
  to_status VARCHAR(32) NOT NULL,
  reason VARCHAR(255) NOT NULL DEFAULT '',
  changed_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  KEY idx_feedback_status_histories_feedback_id (feedback_id),
  CONSTRAINT fk_feedback_status_histories_feedback FOREIGN KEY (feedback_id) REFERENCES feedbacks(id) ON DELETE RESTRICT,
  CONSTRAINT fk_feedback_status_histories_user FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
