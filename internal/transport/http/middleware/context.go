package middleware
import "context"
type contextKey string
const(userIDKey contextKey="user_id";usernameKey contextKey="username";permissionsKey contextKey="permissions";requestIDKey contextKey="request_id")
func UserID(ctx context.Context)int64{value,_:=ctx.Value(userIDKey).(int64);return value}
func Username(ctx context.Context)string{value,_:=ctx.Value(usernameKey).(string);return value}
func Permissions(ctx context.Context)[]string{value,_:=ctx.Value(permissionsKey).([]string);return value}
func RequestID(ctx context.Context)string{value,_:=ctx.Value(requestIDKey).(string);return value}