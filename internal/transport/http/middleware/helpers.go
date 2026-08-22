package middleware
import("context";"strconv")
func UserIDString(ctx context.Context)string{return strconv.FormatInt(UserID(ctx),10)}