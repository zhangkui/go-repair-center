package service

import(
 "context"
 "errors"
 "fmt"
 "golang.org/x/crypto/bcrypt"
 "go-repair-center/internal/repository"
 "go-repair-center/internal/repository/mysql"
)

type UserService struct{*ResourceService;users *mysql.UserRepository}
func NewUserService(repo *mysql.UserRepository)*UserService{return &UserService{ResourceService:NewResourceService(repo,"user","username","display_name"),users:repo}}
func(s *UserService)CreateUser(ctx context.Context,input repository.Record)(repository.Record,error){password:=fmt.Sprint(input["password"]);if err:=validatePassword(password);err!=nil{return nil,err};hash,err:=bcrypt.GenerateFromPassword([]byte(password),12);if err!=nil{return nil,err};delete(input,"password");input["password_hash"]=string(hash);if _,ok:=input["status"];!ok{input["status"]="ACTIVE"};return s.Create(ctx,input)}
func(s *UserService)ResetPassword(ctx context.Context,userID int64,password string)error{if err:=validatePassword(password);err!=nil{return err};hash,err:=bcrypt.GenerateFromPassword([]byte(password),12);if err!=nil{return err};if _,err=s.Update(ctx,userID,repository.Record{"password_hash":string(hash)});err!=nil{return err};return s.users.RevokeSessions(ctx,userID)}
func(s *UserService)SetEnabled(ctx context.Context,userID int64,enabled bool)(repository.Record,error){status:="DISABLED";if enabled{status="ACTIVE"};user,err:=s.Update(ctx,userID,repository.Record{"status":status});if err==nil&&!enabled{err=s.users.RevokeSessions(ctx,userID)};return user,err}
func(s *UserService)AssignRole(ctx context.Context,userID,roleID int64)error{if userID<=0||roleID<=0{return errors.New("invalid user or role")};return s.users.AssignRole(ctx,userID,roleID)}