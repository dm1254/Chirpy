package main
import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"workspace/github.com/dm1254/Chirpy/internal/auth"
)
func(c *ApiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request){
	type requestData struct{
		Event string `json:"event"`
		Data struct{
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	
	}
	key,err := auth.GetApiKey(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "No auth header", err)
		return
	}
	if key != c.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "User not authorized",err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err = decoder.Decode(&reqData)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError,"Couldn't decode parameters",err)
		return 
	}
	if reqData.Event != "user.upgraded"{
		respondWithError(w,http.StatusNoContent,"user not upgraded",err)
		return
	}
	err = c.db.UpgradeUserToRed(r.Context(), reqData.Data.UserID)
	if err != nil{
		respondWithError(w, http.StatusNotFound, "user not found",err)
	}
	w.WriteHeader(204)



}
