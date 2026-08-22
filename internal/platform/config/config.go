package config

import (
 "fmt"
 "os"
 "strconv"
 "time"
 "github.com/joho/godotenv"
)

type Config struct {
 AppName, AppEnv, AppPort string
 ReadTimeout, WriteTimeout time.Duration
 RequestMaxSize int64
 DBHost, DBPort, DBUser, DBPassword, DBName string
 DBMaxOpenConns, DBMaxIdleConns int
 DBConnMaxLifetime time.Duration
 RedisHost, RedisPort, RedisPassword string
 RedisDB, RedisPoolSize int
 JWTSecret string
 JWTAccessTTL, JWTRefreshTTL time.Duration
 LoginMaxAttempts int
 LoginLockoutDuration, IdempotencyTTL time.Duration
 QuotationValidDays, WarrantyDefaultDays, FeedbackGenerateDays int
 QuotationApprovalThreshold float64
 AdminUsername, AdminPassword string
}

func Load() (*Config,error) {
 _=godotenv.Load()
 c:=&Config{AppName:env("APP_NAME","go-repair-center"),AppEnv:env("APP_ENV","development"),AppPort:env("APP_PORT","8080"),DBHost:env("DB_HOST","mysql"),DBPort:env("DB_PORT","3306"),DBUser:env("DB_USER","repair_user"),DBPassword:env("DB_PASSWORD","repair_pass_2026"),DBName:env("DB_NAME","go_repair_center"),RedisHost:env("REDIS_HOST","redis"),RedisPort:env("REDIS_PORT","6379"),RedisPassword:env("REDIS_PASSWORD",""),JWTSecret:env("JWT_SECRET","your-256-bit-secret-change-in-production"),AdminUsername:env("ADMIN_USERNAME","admin"),AdminPassword:env("ADMIN_PASSWORD","Admin123!")}
 var err error
 if c.ReadTimeout,err=duration("APP_READ_TIMEOUT","30s");err!=nil{return nil,err}; if c.WriteTimeout,err=duration("APP_WRITE_TIMEOUT","30s");err!=nil{return nil,err}
 if c.JWTAccessTTL,err=duration("JWT_ACCESS_TTL","15m");err!=nil{return nil,err}; if c.JWTRefreshTTL,err=duration("JWT_REFRESH_TTL","168h");err!=nil{return nil,err}
 if c.LoginLockoutDuration,err=duration("LOGIN_LOCKOUT_DURATION","15m");err!=nil{return nil,err}; if c.IdempotencyTTL,err=duration("IDEMPOTENCY_TTL","24h");err!=nil{return nil,err}; if c.DBConnMaxLifetime,err=duration("DB_CONN_MAX_LIFETIME","300s");err!=nil{return nil,err}
 c.RequestMaxSize=int64(integer("APP_REQUEST_MAX_SIZE",10485760)); c.DBMaxOpenConns=integer("DB_MAX_OPEN_CONNS",25); c.DBMaxIdleConns=integer("DB_MAX_IDLE_CONNS",10); c.RedisDB=integer("REDIS_DB",0); c.RedisPoolSize=integer("REDIS_POOL_SIZE",20); c.LoginMaxAttempts=integer("LOGIN_MAX_ATTEMPTS",5); c.QuotationValidDays=integer("QUOTATION_VALID_DAYS",7); c.WarrantyDefaultDays=integer("WARRANTY_DEFAULT_DAYS",90); c.FeedbackGenerateDays=integer("FEEDBACK_AUTO_GENERATE_DAYS",3); c.QuotationApprovalThreshold,_=strconv.ParseFloat(env("QUOTATION_AUTO_APPROVAL_THRESHOLD","1000"),64)
 return c,nil
}
func (c *Config) MySQLDSN() string{return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&multiStatements=true",c.DBUser,c.DBPassword,c.DBHost,c.DBPort,c.DBName)}
func (c *Config) RedisAddress()string{return c.RedisHost+":"+c.RedisPort}
func env(k,d string)string{if v:=os.Getenv(k);v!=""{return v};return d}
func integer(k string,d int)int{v,e:=strconv.Atoi(env(k,strconv.Itoa(d)));if e!=nil{return d};return v}
func duration(k,d string)(time.Duration,error){v,e:=time.ParseDuration(env(k,d));if e!=nil{return 0,fmt.Errorf("invalid %s: %w",k,e)};return v,nil}