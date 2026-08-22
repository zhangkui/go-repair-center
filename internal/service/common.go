package service

import (
 "context"
 "errors"
 "fmt"
 "strings"
 "go-repair-center/internal/repository"
)

var ErrForbidden=errors.New("forbidden")
var ErrInvalidTransition=errors.New("invalid status transition")
var ErrValidation=errors.New("validation failed")

type ResourceService struct { repo repository.CRUDRepository; resource string; required []string }
func NewResourceService(repo repository.CRUDRepository,resource string,required ...string)*ResourceService{return &ResourceService{repo:repo,resource:resource,required:required}}
func(s *ResourceService)Create(ctx context.Context,input repository.Record)(repository.Record,error){if err:=s.validate(input,true);err!=nil{return nil,err};return s.repo.Create(ctx,input)}
func(s *ResourceService)Get(ctx context.Context,id int64)(repository.Record,error){if id<=0{return nil,fmt.Errorf("%w: invalid id",ErrValidation)};return s.repo.Get(ctx,id)}
func(s *ResourceService)List(ctx context.Context,page,pageSize int,keyword,status string)(repository.Page,error){keyword=strings.TrimSpace(keyword);status=strings.TrimSpace(status);return s.repo.List(ctx,page,pageSize,keyword,status)}
func(s *ResourceService)Update(ctx context.Context,id int64,input repository.Record)(repository.Record,error){if id<=0{return nil,fmt.Errorf("%w: invalid id",ErrValidation)};if err:=s.validate(input,false);err!=nil{return nil,err};return s.repo.Update(ctx,id,input)}
func(s *ResourceService)Delete(ctx context.Context,id int64)error{if id<=0{return fmt.Errorf("%w: invalid id",ErrValidation)};return s.repo.Delete(ctx,id)}
func(s *ResourceService)validate(input repository.Record,creating bool)error{if len(input)==0{return fmt.Errorf("%w: empty %s",ErrValidation,s.resource)};if creating{for _,field:=range s.required{value,ok:=input[field];if !ok||strings.TrimSpace(fmt.Sprint(value))==""{return fmt.Errorf("%w: %s is required",ErrValidation,field)}}};return nil}
func ValidateTransition(current,next string,allowed map[string][]string)error{for _,candidate:=range allowed[current]{if candidate==next{return nil}};return fmt.Errorf("%w: %s -> %s",ErrInvalidTransition,current,next)}