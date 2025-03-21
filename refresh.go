package main
import(
	"net/http"
	"fmt"
	"time"
	"workspace/github.com/dm1254/Chirpy/internal/auth"
)

func (c *ApiConfig) handleRefresh(w http.ResponseWriter, r *http.Request){
	type responseData struct{
		Token string `json:"token"`

	}
	reqHeaderToken,err := auth.GetBearerToken(r.Header)
	fmt.Println(reqHeaderToken)
	fmt.Println("Authorization header value:", r.Header.Get("Authorization"))
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Could not get token",err)
	}
	user,err := c.db.GetRefreshTokenUser(r.Context(), reqHeaderToken)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Could not get user refresh token",err)
		return	
	}
	accessToken,err := auth.MakeJWT(user.ID,c.JWTSecret,time.Hour)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Token is invalid", err)
		return
	}
	respondWithJson(w,http.StatusOK, responseData{
			Token: accessToken,
	})

}
