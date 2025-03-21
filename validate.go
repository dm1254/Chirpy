package main
import(
	"net/http"
	"encoding/json"
	"strings"
	"fmt"
)
type requestData struct{
	Body string `json:"body"`
}

type responseData struct{
	Body string `json:"body"`
}

func (c *ApiConfig) HandlerValidate (w http.ResponseWriter,r *http.Request){
	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err := decoder.Decode(&reqData)
	if err != nil{
		w.WriteHeader(500)
	}

	if len(reqData.Body) <= 140{
		check_body := checkProfan(reqData.Body)
		fmt.Println(check_body)
		if check_body != reqData.Body{
			response := map[string]interface{}{
				"cleaned_body": check_body,
			}
			
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(response)	
		}else{
			response := map[string]interface{}{
				"cleaned_body": reqData.Body, 
			}
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(response)	
		}
	}		
	

	if len(reqData.Body) > 140{
		response := map[string]interface{}{
			"error": "Chirp is too long",
		}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response)
	
	}

}

func checkProfan(body string) string{
	split_body := strings.Fields(body)
	
	for i,word := range split_body {
		if word == "Kerfuffle" || word == "kerfuffle"{
			split_body[i] = "****"
		}
		if word == "sharbert" || word == "Sharbert"{
			split_body[i] = "****"			
		}
		if word == "fornax" ||word == "Fornax"{
			split_body[i] = "****"
		}
	}
	joined_body := strings.Join(split_body, " ")

	return joined_body




}
