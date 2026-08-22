INSERT IGNORE INTO system_configs(config_key, config_value, value_type, description)
VALUES
('quotation_auto_approval_threshold', '1000', 'number', 'Quotation auto approval threshold'),
('warranty_default_days', '90', 'number', 'Default warranty days'),
('feedback_auto_generate_days', '3', 'number', 'Feedback auto generate days');

INSERT IGNORE INTO fault_codes(code, name, category, description, status)
VALUES
('E01', 'Power failure', 'appliance', 'Power failure', 'ACTIVE'),
('E02', 'Main board failure', 'appliance', 'Main board failure', 'ACTIVE'),
('E03', 'Compressor failure', 'appliance', 'Compressor failure', 'ACTIVE'),
('E04', 'Sensor failure', 'appliance', 'Sensor failure', 'ACTIVE');
