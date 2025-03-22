package auth
import(
	"strings"
	"net/http"
	"errors"
)




func GetApiKey(headers http.Header) (string,error){
	authHeader := headers.Get("Authorization")
	if authHeader == ""{
		return "",ErrNoAuthHeaderIncluded
	}
	splitHeader := strings.Split(authHeader, " ")
	if len(splitHeader) < 2 || splitHeader[0] != "ApiKey"{
		return "", errors.New("malformed authorization header")
	}
	return strings.TrimSpace(splitHeader[1]), nil
}
