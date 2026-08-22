package bootstrap

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"

	"go-repair-center/internal/platform/config"
)

type SeedRole struct {
	Code string
	Name string
}

type SeedPermission struct {
	Code   string
	Name   string
	Module string
	Action string
}

func SeedDefaults(ctx context.Context, db *sql.DB, cfg *config.Config) error {
	if err := seedRoles(ctx, db); err != nil {
		return err
	}
	if err := seedPermissions(ctx, db); err != nil {
		return err
	}
	if err := seedRolePermissions(ctx, db); err != nil {
		return err
	}
	if err := seedConfigs(ctx, db); err != nil {
		return err
	}
	if err := seedFaultCodes(ctx, db); err != nil {
		return err
	}
	if err := seedAdmin(ctx, db, cfg); err != nil {
		return err
	}
	return nil
}

func seedRoles(ctx context.Context, db *sql.DB) error {
	roles := []SeedRole{
		{Code: "ADMIN", Name: "业务管理员"},
		{Code: "OPERATOR", Name: "日常操作人员"},
		{Code: "REVIEWER", Name: "审批复核人员"},
		{Code: "AUDITOR", Name: "只读审计人员"},
	}
	for _, role := range roles {
		if _, err := db.ExecContext(ctx, `
INSERT INTO roles(code,name,description)
VALUES(?,?,?)
ON DUPLICATE KEY UPDATE name = VALUES(name), description = VALUES(description)
`, role.Code, role.Name, role.Name); err != nil {
			return err
		}
	}
	return nil
}

func seedPermissions(ctx context.Context, db *sql.DB) error {
	permissions := []SeedPermission{
		{Code: "system:*", Name: "系统全部权限", Module: "system", Action: "*"},
		{Code: "system:user:manage", Name: "用户与角色管理", Module: "system", Action: "user:manage"},
		{Code: "customer:read", Name: "查看客户", Module: "customer", Action: "read"},
		{Code: "customer:write", Name: "维护客户", Module: "customer", Action: "write"},
		{Code: "device:read", Name: "查看设备", Module: "device", Action: "read"},
		{Code: "device:write", Name: "维护设备", Module: "device", Action: "write"},
		{Code: "repair_order:read", Name: "查看维修单", Module: "repair_order", Action: "read"},
		{Code: "repair_order:write", Name: "维护维修单", Module: "repair_order", Action: "write"},
		{Code: "quotation:read", Name: "查看报价单", Module: "quotation", Action: "read"},
		{Code: "quotation:write", Name: "维护报价单", Module: "quotation", Action: "write"},
		{Code: "repair_execution:read", Name: "查看维修执行", Module: "repair_execution", Action: "read"},
		{Code: "repair_execution:write", Name: "维护维修执行", Module: "repair_execution", Action: "write"},
		{Code: "part:read", Name: "查看配件", Module: "part", Action: "read"},
		{Code: "part:write", Name: "维护配件", Module: "part", Action: "write"},
		{Code: "warranty:read", Name: "查看质保", Module: "warranty", Action: "read"},
		{Code: "warranty:write", Name: "维护质保", Module: "warranty", Action: "write"},
		{Code: "feedback:read", Name: "查看回访", Module: "feedback", Action: "read"},
		{Code: "feedback:write", Name: "维护回访", Module: "feedback", Action: "write"},
		{Code: "audit:read", Name: "查看审计日志", Module: "audit", Action: "read"},
	}
	for _, permission := range permissions {
		if _, err := db.ExecContext(ctx, `
INSERT INTO permissions(code,name,module,action,description)
VALUES(?,?,?,?,?)
ON DUPLICATE KEY UPDATE name = VALUES(name), module = VALUES(module), action = VALUES(action), description = VALUES(description)
`, permission.Code, permission.Name, permission.Module, permission.Action, permission.Name); err != nil {
			return err
		}
	}
	return nil
}

func seedRolePermissions(ctx context.Context, db *sql.DB) error {
	if err := assignAllPermissionsToRole(ctx, db, "ADMIN"); err != nil {
		return err
	}
	if err := assignPermissionsToRole(ctx, db, "OPERATOR", []string{
		"customer:read",
		"customer:write",
		"device:read",
		"device:write",
		"repair_order:read",
		"repair_order:write",
		"quotation:read",
		"part:read",
		"part:write",
		"repair_execution:read",
		"repair_execution:write",
		"feedback:read",
		"feedback:write",
	}); err != nil {
		return err
	}
	if err := assignPermissionsToRole(ctx, db, "REVIEWER", []string{
		"quotation:read",
		"quotation:write",
		"warranty:read",
		"warranty:write",
		"repair_execution:read",
	}); err != nil {
		return err
	}
	if err := assignPermissionsToRole(ctx, db, "AUDITOR", []string{
		"customer:read",
		"device:read",
		"repair_order:read",
		"quotation:read",
		"repair_execution:read",
		"part:read",
		"warranty:read",
		"feedback:read",
		"audit:read",
	}); err != nil {
		return err
	}
	return nil
}

func assignAllPermissionsToRole(ctx context.Context, db *sql.DB, roleCode string) error {
	_, err := db.ExecContext(ctx, `
INSERT IGNORE INTO role_permissions(role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = ?
`, roleCode)
	return err
}

func assignPermissionsToRole(ctx context.Context, db *sql.DB, roleCode string, permissionCodes []string) error {
	if len(permissionCodes) == 0 {
		return nil
	}
	roleID, err := findRoleID(ctx, db, roleCode)
	if err != nil {
		return err
	}
	for _, permissionCode := range permissionCodes {
		permissionID, err := findPermissionID(ctx, db, permissionCode)
		if err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, `
INSERT IGNORE INTO role_permissions(role_id, permission_id)
VALUES(?, ?)
`, roleID, permissionID); err != nil {
			return err
		}
	}
	return nil
}

func seedConfigs(ctx context.Context, db *sql.DB) error {
	configs := map[string]string{
		"quotation_auto_approval_threshold": "1000",
		"warranty_default_days":            "90",
		"feedback_auto_generate_days":      "3",
	}
	for key, value := range configs {
		if _, err := db.ExecContext(ctx, `
INSERT IGNORE INTO system_configs(config_key,config_value,value_type,description)
VALUES(?,?,?,?)
`, key, value, "string", key); err != nil {
			return err
		}
	}
	return nil
}

func seedFaultCodes(ctx context.Context, db *sql.DB) error {
	faultCodes := []SeedPermission{
		{Code: "E01", Name: "电源故障", Module: "家电", Action: "fault"},
		{Code: "E02", Name: "主板故障", Module: "家电", Action: "fault"},
		{Code: "E03", Name: "压缩机故障", Module: "家电", Action: "fault"},
		{Code: "E04", Name: "传感器故障", Module: "家电", Action: "fault"},
	}
	for _, fault := range faultCodes {
		if _, err := db.ExecContext(ctx, `
INSERT INTO fault_codes(code,name,category,description,status)
VALUES(?,?,?,?,?)
ON DUPLICATE KEY UPDATE name = VALUES(name), category = VALUES(category), description = VALUES(description), status = VALUES(status)
`, fault.Code, fault.Name, fault.Module, fault.Name, "ACTIVE"); err != nil {
			return err
		}
	}
	return nil
}

func seedAdmin(ctx context.Context, db *sql.DB, cfg *config.Config) error {
	roleID, err := findRoleID(ctx, db, "ADMIN")
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), 12)
	if err != nil {
		return err
	}

	var userID int64
	err = db.QueryRowContext(ctx, `
SELECT id
FROM users
WHERE username = ? AND deleted_at IS NULL
LIMIT 1
`, cfg.AdminUsername).Scan(&userID)
	if err == sql.ErrNoRows {
		result, createErr := db.ExecContext(ctx, `
INSERT INTO users(username,password_hash,display_name,status)
VALUES(?,?,?,?)
`, cfg.AdminUsername, string(passwordHash), "Local Admin", "ACTIVE")
		if createErr != nil {
			return createErr
		}
		userID, _ = result.LastInsertId()
	} else if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
INSERT IGNORE INTO user_roles(user_id,role_id)
VALUES(?,?)
`, userID, roleID)
	return err
}

func findRoleID(ctx context.Context, db *sql.DB, code string) (int64, error) {
	var roleID int64
	err := db.QueryRowContext(ctx, `
SELECT id
FROM roles
WHERE code = ? AND deleted_at IS NULL
LIMIT 1
`, code).Scan(&roleID)
	return roleID, err
}

func findPermissionID(ctx context.Context, db *sql.DB, code string) (int64, error) {
	var permissionID int64
	err := db.QueryRowContext(ctx, `
SELECT id
FROM permissions
WHERE code = ? AND deleted_at IS NULL
LIMIT 1
`, code).Scan(&permissionID)
	return permissionID, err
}

func SeededAt() time.Time {
	return time.Now()
}
