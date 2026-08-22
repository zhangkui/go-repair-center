package handler

import(
 "context"
 "encoding/json"
 "errors"
 "net/http"
 "strconv"

 "github.com/go-chi/chi/v5"
 "go-repair-center/internal/repository"
 "go-repair-center/internal/service"
 "go-repair-center/internal/transport/http/middleware"
 "go-repair-center/internal/transport/http/response"
)

type CRUDService interface{Create(context.Context,repository.Record)(repository.Record,error);Get(context.Context,int64)(repository.Record,error);List(context.Context,int,int,string,string)(repository.Page,error);Update(context.Context,int64,repository.Record)(repository.Record,error);Delete(context.Context,int64)error}
type GenericHandler struct{service CRUDService;resource string}
func NewGenericHandler(service CRUDService,resource string)*GenericHandler{return &GenericHandler{service:service,resource:resource}}
func(h *GenericHandler)List(w http.ResponseWriter,r *http.Request){page:=queryInt(r,"page",1);pageSize:=queryInt(r,"page_size",20);result,err:=h.service.List(r.Context(),page,pageSize,r.URL.Query().Get("keyword"),r.URL.Query().Get("status"));if err!=nil{writeError(w,r,err);return};response.Success(w,response.PageData{Items:result.Items,Pagination:response.Pagination{Page:result.Page,PageSize:result.PageSize,Total:result.Total,TotalPages:result.TotalPages}},middleware.RequestID(r.Context()))}
func(h *GenericHandler)Get(w http.ResponseWriter,r *http.Request){id,err:=pathID(r);if err!=nil{writeError(w,r,err);return};item,err:=h.service.Get(r.Context(),id);if err!=nil{writeError(w,r,err);return};response.Success(w,item,middleware.RequestID(r.Context()))}
func(h *GenericHandler)Create(w http.ResponseWriter,r *http.Request){input:=repository.Record{};if err:=decode(r,&input);err!=nil{writeError(w,r,err);return};item,err:=h.service.Create(r.Context(),input);if err!=nil{writeError(w,r,err);return};response.Created(w,item,middleware.RequestID(r.Context()))}
func(h *GenericHandler)Update(w http.ResponseWriter,r *http.Request){id,err:=pathID(r);if err!=nil{writeError(w,r,err);return};input:=repository.Record{};if err=decode(r,&input);err!=nil{writeError(w,r,err);return};item,err:=h.service.Update(r.Context(),id,input);if err!=nil{writeError(w,r,err);return};response.Success(w,item,middleware.RequestID(r.Context()))}
func(h *GenericHandler)Delete(w http.ResponseWriter,r *http.Request){id,err:=pathID(r);if err!=nil{writeError(w,r,err);return};if err=h.service.Delete(r.Context(),id);err!=nil{writeError(w,r,err);return};response.Success(w,map[string]any{"deleted":true,"id":id},middleware.RequestID(r.Context()))}
func(h *GenericHandler)Register(router chi.Router,readPermission,writePermission string){router.With(middleware.Require(readPermission)).Get("/",h.List);router.With(middleware.Require(readPermission)).Get("/{id}",h.Get);router.With(middleware.Require(writePermission)).Post("/",h.Create);router.With(middleware.Require(writePermission)).Put("/{id}",h.Update);router.With(middleware.Require(writePermission)).Delete("/{id}",h.Delete)}
func decode(r *http.Request,target any)error{decoder:=json.NewDecoder(r.Body);decoder.DisallowUnknownFields();return decoder.Decode(target)}
func pathID(r *http.Request)(int64,error){return strconv.ParseInt(chi.URLParam(r,"id"),10,64)}
func queryInt(r *http.Request,key string,fallback int)int{value,err:=strconv.Atoi(r.URL.Query().Get(key));if err!=nil{return fallback};return value}
func writeError(w http.ResponseWriter,r *http.Request,err error){status,code:=500,response.CodeInternal;switch{case errors.Is(err,service.ErrValidation):status,code=400,response.CodeValidation;case errors.Is(err,service.ErrInvalidTransition):status,code=409,response.CodeInvalidTransition;case errors.Is(err,service.ErrForbidden):status,code=403,response.CodeForbidden};response.Error(w,status,code,err.Error(),middleware.RequestID(r.Context()))}
