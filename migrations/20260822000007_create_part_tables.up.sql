CREATE TABLE IF NOT EXISTS parts (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(100) NOT NULL,
  specification VARCHAR(255) NOT NULL DEFAULT '',
  unit VARCHAR(32) NOT NULL DEFAULT '',
  unit_price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  stock_quantity INT NOT NULL DEFAULT 0,
  safety_stock INT NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_parts_code (code),
  KEY idx_parts_status (status),
  KEY idx_parts_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS part_usage_records (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  repair_execution_id BIGINT UNSIGNED NOT NULL,
  part_id BIGINT UNSIGNED NOT NULL,
  quantity INT NOT NULL DEFAULT 0,
  unit_price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  used_by BIGINT UNSIGNED NOT NULL,
  idempotency_key VARCHAR(64) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_part_usage_idempotency_key (idempotency_key),
  KEY idx_part_usage_records_execution_id (repair_execution_id),
  KEY idx_part_usage_records_part_id (part_id),
  CONSTRAINT fk_part_usage_records_execution FOREIGN KEY (repair_execution_id) REFERENCES repair_executions(id) ON DELETE RESTRICT,
  CONSTRAINT fk_part_usage_records_part FOREIGN KEY (part_id) REFERENCES parts(id) ON DELETE RESTRICT,
  CONSTRAINT fk_part_usage_records_user FOREIGN KEY (used_by) REFERENCES users(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS part_stock_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  part_id BIGINT UNSIGNED NOT NULL,
  change_type VARCHAR(32) NOT NULL,
  quantity_before INT NOT NULL DEFAULT 0,
  change_quantity INT NOT NULL DEFAULT 0,
  quantity_after INT NOT NULL DEFAULT 0,
  reference_type VARCHAR(64) NOT NULL DEFAULT '',
  reference_id BIGINT UNSIGNED NULL,
  operator_id BIGINT UNSIGNED NOT NULL,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  KEY idx_part_stock_logs_part_id (part_id),
  KEY idx_part_stock_logs_created_at (created_at),
  CONSTRAINT fk_part_stock_logs_part FOREIGN KEY (part_id) REFERENCES parts(id) ON DELETE RESTRICT,
  CONSTRAINT fk_part_stock_logs_operator FOREIGN KEY (operator_id) REFERENCES users(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE quotation_items
  ADD CONSTRAINT fk_quotation_items_part
  FOREIGN KEY (part_id) REFERENCES parts(id) ON DELETE RESTRICT;
