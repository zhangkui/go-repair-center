package service

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "errors"
 "fmt"
 "strconv"
 "strings"
 "time"

 redisv9 "github.com/redis/go-redis/v9"
 "golang.org/x/crypto/bcrypt"
 "go-repair-center/internal/platform/jwt"
 "go-repair-center/internal/repository"
 "go-repair-center/internal/repository/mysql"
)

var ErrInvalidCredentials=errors.New("invalid username or password")
var ErrAccountDisabled=errors.New("account disabled")
var ErrLoginLocked=errors.New("login temporarily locked")

type AuthService struct{users *mysql.UserRepository;redis *redisv9.Client;jwt *jwt.Manager;maxAttempts int;lockout time.Duration}
type TokenPair struct{AccessToken string `json:"access_token"`;RefreshToken string `json:"refresh_token"`;TokenType string `json:"token_type"`;ExpiresIn int64 `json:"expires_in"`}
func NewAuthService(users *mysql.UserRepository,redis *redisv9.Client,jwtManager *jwt.Manager,maxAttempts int,lockout time.Duration)*AuthService{return &AuthService{users:users,redis:redis,jwt:jwtManager,maxAttempts:maxAttempts,lockout:lockout}}
func(s *AuthService)Register(ctx context.Context,username,password,displayName string)(repository.Record,error){username=strings.TrimSpace(username);if len(username)<3{return nil,fmt.Errorf("%w: username too short",ErrValidation)};if err:=validatePassword(password);err!=nil{return nil,err};if _,err:=s.users.FindByUsername(ctx,username);err==nil{return nil,errors.New("username already exists")};hash,err:=bcrypt.GenerateFromPassword([]byte(password),12);if err!=nil{return nil,err};user,err:=s.users.Create(ctx,repository.Record{"username":username,"password_hash":string(hash),"display_name":displayName,"status":"ACTIVE"});if err!=nil{return nil,err};id,_:=recordInt64(user,"id");var roleID int64;if err=s.users.DB().QueryRowContext(ctx,"SELECT id FROM roles WHERE code='OPERATOR' AND deleted_at IS NULL").Scan(&roleID);err==nil{err=s.users.AssignRole(ctx,id,roleID)};return user,err}
func(s *AuthService)Login(ctx context.Context,username,password,ip string)(TokenPair,error){key:="login:fail:"+ip;attempts,_:=s.redis.Get(ctx,key).Int();if attempts>=s.maxAttempts{return TokenPair{},ErrLoginLocked};user,err:=s.users.FindByUsername(ctx,username);if err!=nil{s.recordFailure(ctx,key);return TokenPair{},ErrInvalidCredentials};if fmt.Sprint(user["status"])!="ACTIVE"{return TokenPair{},ErrAccountDisabled};hash:=fmt.Sprint(user["password_hash"]);if bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))!=nil{s.recordFailure(ctx,key);return TokenPair{},ErrInvalidCredentials};_ = s.redis.Del(ctx,key).Err();userID,_:=recordInt64(user,"id");permissions,err:=s.users.Permissions(ctx,userID);if err!=nil{return TokenPair{},err};pair,err:=s.issue(ctx,userID,username,permissions);if err==nil{_ = s.users.UpdateLastLogin(ctx,userID)};return pair,err}
func(s *AuthService)Refresh(ctx context.Context,raw string)(TokenPair,error){claims,err:=s.jwt.Parse(raw,"refresh");if err!=nil{return TokenPair{},ErrInvalidCredentials};hash:=tokenHash(raw);userID,err:=s.users.ConsumeRefreshToken(ctx,hash);if err!=nil||userID!=claims.UserID{return TokenPair{},ErrInvalidCredentials};user,err:=s.users.Get(ctx,userID);if err!=nil{return TokenPair{},err};if fmt.Sprint(user["status"])!="ACTIVE"{return TokenPair{},ErrAccountDisabled};permissions,err:=s.users.Permissions(ctx,userID);if err!=nil{return TokenPair{},err};return s.issue(ctx,userID,fmt.Sprint(user["username"]),permissions)}
func(s *AuthService)Logout(ctx context.Context,raw string)error{_,err:=s.users.ConsumeRefreshToken(ctx,tokenHash(raw));if err!=nil{return nil};return nil}
func(s *AuthService)Me(ctx context.Context,userID int64)(repository.Record,error){user,err:=s.users.Get(ctx,userID);if err!=nil{return nil,err};permissions,err:=s.users.Permissions(ctx,userID);if err!=nil{return nil,err};user["permissions"]=permissions;delete(user,"password_hash");return user,nil}
func(s *AuthService)ChangePassword(ctx context.Context,userID int64,oldPassword,newPassword string)error{user,err:=s.users.Get(ctx,userID);if err!=nil{return err};if bcrypt.CompareHashAndPassword([]byte(fmt.Sprint(user["password_hash"])),[]byte(oldPassword))!=nil{return ErrInvalidCredentials};if err=validatePassword(newPassword);err!=nil{return err};hash,err:=bcrypt.GenerateFromPassword([]byte(newPassword),12);if err!=nil{return err};_,err=s.users.Update(ctx,userID,repository.Record{"password_hash":string(hash)});if err==nil{err=s.users.RevokeSessions(ctx,userID)};return err}
func(s *AuthService)issue(ctx context.Context,userID int64,username string,permissions []string)(TokenPair,error){access,accessExp,err:=s.jwt.Issue(userID,username,permissions,"access");if err!=nil{return TokenPair{},err};refresh,refreshExp,err:=s.jwt.Issue(userID,username,nil,"refresh");if err!=nil{return TokenPair{},err};if err=s.users.SaveRefreshToken(ctx,userID,tokenHash(refresh),refreshExp);err!=nil{return TokenPair{},err};return TokenPair{access,refresh,"Bearer",int64(time.Until(accessExp).Seconds())},nil}
func(s *AuthService)recordFailure(ctx context.Context,key string){count,_:=s.redis.Incr(ctx,key).Result();if count==1{_ = s.redis.Expire(ctx,key,s.lockout).Err()}}
func tokenHash(value string)string{sum:=sha256.Sum256([]byte(value));return hex.EncodeToString(sum[:])}
func validatePassword(value string)error{if len(value)<8{return fmt.Errorf("%w: password must contain at least 8 characters",ErrValidation)};var upper,digit,special bool;for _,r:=range value{switch{case r>='A'&&r<='Z':upper=true;case r>='0'&&r<='9':digit=true;case !((r>='a'&&r<='z')||(r>='A'&&r<='Z')):special=true}};if !upper||!digit||!special{return fmt.Errorf("%w: password complexity not met",ErrValidation)};return nil}
func recordInt64(record repository.Record,key string)(int64,error){switch value:=record[key].(type){case int64:return value,nil;case int:return int64(value),nil;case []byte:return strconv.ParseInt(string(value),10,64);case string:return strconv.ParseInt(value,10,64);default:return 0,fmt.Errorf("invalid %s",key)}}