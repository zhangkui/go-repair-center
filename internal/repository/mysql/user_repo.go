package mysql

import (
 "context"
 "database/sql"
 "errors"
 "time"
 "go-repair-center/internal/repository"
)

type UserRepository struct { *CRUD; db *sql.DB }
func NewUserRepository(db *sql.DB)*UserRepository{return &UserRepository{CRUD:NewCRUD(db,"users",[]string{"username","display_name","phone","email"},[]string{"username","password_hash","display_name","phone","email","status"}),db:db}}
func (r *UserRepository) FindByUsername(ctx context.Context,username string)(repository.Record,error){rows,err:=r.db.QueryContext(ctx,"SELECT * FROM users WHERE username=? AND deleted_at IS NULL LIMIT 1",username);if err!=nil{return nil,err};defer rows.Close();if !rows.Next(){return nil,ErrNotFound};return scanRecord(rows)}
func (r *UserRepository) Permissions(ctx context.Context,userID int64)([]string,error){rows,err:=r.db.QueryContext(ctx,`SELECT DISTINCT p.code FROM permissions p JOIN role_permissions rp ON rp.permission_id=p.id JOIN user_roles ur ON ur.role_id=rp.role_id JOIN roles r ON r.id=ur.role_id WHERE ur.user_id=? AND p.deleted_at IS NULL AND r.deleted_at IS NULL ORDER BY p.code`,userID);if err!=nil{return nil,err};defer rows.Close();items:=[]string{};for rows.Next(){var code string;if err=rows.Scan(&code);err!=nil{return nil,err};items=append(items,code)};return items,rows.Err()}
func (r *UserRepository) AssignRole(ctx context.Context,userID,roleID int64)error{_,err:=r.db.ExecContext(ctx,"INSERT IGNORE INTO user_roles(user_id,role_id) VALUES(?,?)",userID,roleID);return err}
func (r *UserRepository) RevokeSessions(ctx context.Context,userID int64)error{_,err:=r.db.ExecContext(ctx,"UPDATE refresh_tokens SET revoked_at=NOW() WHERE user_id=? AND revoked_at IS NULL",userID);return err}
func (r *UserRepository) SaveRefreshToken(ctx context.Context,userID int64,hash string,expires time.Time)error{_,err:=r.db.ExecContext(ctx,"INSERT INTO refresh_tokens(user_id,token_hash,expires_at) VALUES(?,?,?)",userID,hash,expires);return err}
func (r *UserRepository) ConsumeRefreshToken(ctx context.Context,hash string)(int64,error){tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return 0,err};defer tx.Rollback();var id,userID int64;var expires time.Time;err=tx.QueryRowContext(ctx,"SELECT id,user_id,expires_at FROM refresh_tokens WHERE token_hash=? AND revoked_at IS NULL FOR UPDATE",hash).Scan(&id,&userID,&expires);if err!=nil{return 0,err};if time.Now().After(expires){return 0,errors.New("refresh token expired")};if _,err=tx.ExecContext(ctx,"UPDATE refresh_tokens SET revoked_at=NOW() WHERE id=?",id);err!=nil{return 0,err};if err=tx.Commit();err!=nil{return 0,err};return userID,nil}
func (r *UserRepository) UpdateLastLogin(ctx context.Context,userID int64)error{_,err:=r.db.ExecContext(ctx,"UPDATE users SET last_login_at=NOW() WHERE id=?",userID);return err}
func (r *UserRepository) DB() *sql.DB { return r.db }

