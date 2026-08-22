package response
import("encoding/json";"net/http";"time")
type Envelope struct{Code int `json:"code"`;Message string `json:"message"`;Data any `json:"data,omitempty"`;Timestamp string `json:"timestamp"`;RequestID string `json:"request_id"`}
type Pagination struct{Page int `json:"page"`;PageSize int `json:"page_size"`;Total int64 `json:"total"`;TotalPages int `json:"total_pages"`}
type PageData struct{Items any `json:"items"`;Pagination Pagination `json:"pagination"`}
func JSON(w http.ResponseWriter,status,code int,message string,data any,requestID string){w.Header().Set("Content-Type","application/json; charset=utf-8");w.WriteHeader(status);_=json.NewEncoder(w).Encode(Envelope{code,message,data,time.Now().Format(time.RFC3339),requestID})}
func Success(w http.ResponseWriter,data any,id string){JSON(w,200,0,"success",data,id)}
func Created(w http.ResponseWriter,data any,id string){JSON(w,201,0,"success",data,id)}
func Error(w http.ResponseWriter,status,code int,message,id string){JSON(w,status,code,message,nil,id)}