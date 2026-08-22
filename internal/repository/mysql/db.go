package mysql

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "strings"

 "go-repair-center/internal/repository"
)

var ErrNotFound = errors.New("record not found")
var ErrConflict = errors.New("record conflict")

type DB struct { SQL *sql.DB }
func NewDB(db *sql.DB) *DB { return &DB{SQL: db} }
func (d *DB) WithinTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
 tx, err := d.SQL.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
 if err != nil { return err }
 if err = fn(tx); err != nil { _ = tx.Rollback(); return err }
 return tx.Commit()
}

type CRUD struct { db *sql.DB; table string; searchable []string; writable map[string]bool }
func NewCRUD(db *sql.DB, table string, searchable, writable []string) *CRUD {
 allowed:=make(map[string]bool,len(writable)); for _,column:=range writable { allowed[column]=true }
 return &CRUD{db:db,table:table,searchable:searchable,writable:allowed}
}
func (r *CRUD) Create(ctx context.Context, input repository.Record) (repository.Record,error) {
 columns,args,marks:=r.prepare(input)
 if len(columns)==0{return nil,errors.New("no writable fields")}
 query:=fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",r.table,strings.Join(columns,","),strings.Join(marks,","))
 result,err:=r.db.ExecContext(ctx,query,args...);if err!=nil{return nil,err};id,err:=result.LastInsertId();if err!=nil{return nil,err};return r.Get(ctx,id)
}
func (r *CRUD) Get(ctx context.Context,id int64)(repository.Record,error){
 rows,err:=r.db.QueryContext(ctx,fmt.Sprintf("SELECT * FROM %s WHERE id=? AND deleted_at IS NULL LIMIT 1",r.table),id);if err!=nil{return nil,err};defer rows.Close();if !rows.Next(){return nil,ErrNotFound};return scanRecord(rows)
}
func (r *CRUD) List(ctx context.Context,page,pageSize int,keyword,status string)(repository.Page,error){
 if page<1{page=1};if pageSize<1||pageSize>200{pageSize=20};where:="deleted_at IS NULL";args:=[]any{}
 if status!=""&&r.writable["status"]{where+=" AND status=?";args=append(args,status)}
 if keyword!=""&&len(r.searchable)>0{parts:=make([]string,0,len(r.searchable));for _,column:=range r.searchable{parts=append(parts,column+" LIKE ?");args=append(args,"%"+keyword+"%")};where+=" AND ("+strings.Join(parts," OR ")+")"}
 var total int64;if err:=r.db.QueryRowContext(ctx,fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s",r.table,where),args...).Scan(&total);err!=nil{return repository.Page{},err}
 queryArgs:=append(append([]any{},args...),pageSize,(page-1)*pageSize);rows,err:=r.db.QueryContext(ctx,fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY id DESC LIMIT ? OFFSET ?",r.table,where),queryArgs...);if err!=nil{return repository.Page{},err};defer rows.Close();items:=[]repository.Record{};for rows.Next(){item,e:=scanRecord(rows);if e!=nil{return repository.Page{},e};items=append(items,item)};pages:=int((total+int64(pageSize)-1)/int64(pageSize));return repository.Page{Items:items,Page:page,PageSize:pageSize,Total:total,TotalPages:pages},rows.Err()
}
func (r *CRUD) Update(ctx context.Context,id int64,input repository.Record)(repository.Record,error){
 columns,args,_:=r.prepare(input);if len(columns)==0{return r.Get(ctx,id)};sets:=make([]string,len(columns));for i,column:=range columns{sets[i]=column+"=?"};args=append(args,id);result,err:=r.db.ExecContext(ctx,fmt.Sprintf("UPDATE %s SET %s,updated_at=NOW() WHERE id=? AND deleted_at IS NULL",r.table,strings.Join(sets,",")),args...);if err!=nil{return nil,err};affected,_:=result.RowsAffected();if affected==0{return nil,ErrNotFound};return r.Get(ctx,id)
}
func (r *CRUD) Delete(ctx context.Context,id int64)error{result,err:=r.db.ExecContext(ctx,fmt.Sprintf("UPDATE %s SET deleted_at=NOW(),updated_at=NOW() WHERE id=? AND deleted_at IS NULL",r.table),id);if err!=nil{return err};affected,_:=result.RowsAffected();if affected==0{return ErrNotFound};return nil}
func (r *CRUD) prepare(input repository.Record)([]string,[]any,[]string){columns:=[]string{};args:=[]any{};marks:=[]string{};for key,value:=range input{if r.writable[key]{columns=append(columns,key);args=append(args,value);marks=append(marks,"?")}};return columns,args,marks}
func scanRecord(rows *sql.Rows)(repository.Record,error){columns,err:=rows.Columns();if err!=nil{return nil,err};values:=make([]any,len(columns));pointers:=make([]any,len(columns));for i:=range values{pointers[i]=&values[i]};if err=rows.Scan(pointers...);err!=nil{return nil,err};record:=repository.Record{};for i,column:=range columns{value:=values[i];if raw,ok:=value.([]byte);ok{record[column]=string(raw)}else{record[column]=value}};return record,nil}