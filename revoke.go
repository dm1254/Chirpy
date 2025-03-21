package main 
import (
	"net/http"
	"workspace/github.com/dm1254/Chirpy/internal/auth"

)

func (c *ApiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	reqHeaderToken,err:= auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Could not get token from header",err)
	}

	_,err = c.db.UpdateToken(r.Context(), reqHeaderToken)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError,"Could not update revoked and update timestamp",err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return 

}
