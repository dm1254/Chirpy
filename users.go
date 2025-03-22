package main 
import(

	"net/http"
	"encoding/json"
	"time"
	"log"
	"workspace/github.com/dm1254/Chirpy/internal/database"
	"workspace/github.com/dm1254/Chirpy/internal/auth"
	"github.com/google/uuid"
)


func (c *ApiConfig) handlerUsers(w http.ResponseWriter, r *http.Request){
	type requestData struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type responseData struct{
		ID uuid.UUID`json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		ChirpyRed bool `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err := decoder.Decode(&reqData)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}
	usersData,err := c.db.CreateUsers(r.Context(), database.CreateUsersParams{
		Email: reqData.Email,
		HashedPassword: reqData.Password,
	})
	if err != nil{
		log.Printf("Could not create user: %s",err)
		return 
	}
	respondWithJson(w, http.StatusCreated, responseData{
		ID: usersData.ID,
		CreatedAt: usersData.CreatedAt,
		UpdatedAt: usersData.UpdatedAt,
		Email: usersData.Email,
		ChirpyRed: usersData.IsChirpyRed,
	})
	
	
}

func(c *ApiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
	type requestData struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type responseData struct{
		ID uuid.UUID`json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		ChirpyRed bool `json:"is_chirpy_red"`	
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err := decoder.Decode(&reqData)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError,"Couldn't decode parameters",err)
		return 
	}

	getUser,err := c.db.GetUserPass(r.Context(),reqData.Email)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error getting user",err)	
		return
	}	
	
	JWTtoken,err := auth.MakeJWT(getUser.ID, c.JWTSecret , time.Hour)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Could not create JWT",err)
		return
	}
	
	HashedrequestPassword, err := auth.HashPassword(reqData.Password)
	if err != nil{
		log.Printf("Error returning hashed password: %s",err)
	}
	err = auth.ComparePasswordAndHash(getUser.HashedPassword, HashedrequestPassword)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "Could not generate refresh token", err)
	}
	_,err = c.db.CreateRefreshToken(r.Context(),database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: getUser.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),

	})
	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "Could not store refresh token", err)	
	}
	respondWithJson(w, http.StatusOK, responseData{
		ID: getUser.ID,
		CreatedAt: getUser.CreatedAt,
		UpdatedAt: getUser.UpdatedAt,
		Email: getUser.Email,
		ChirpyRed: getUser.IsChirpyRed,
		Token: JWTtoken, 
		RefreshToken: refreshToken,	
		
	})
	return 
}

func(c *ApiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request){
	type requestData struct{
		Email string `json:"email"`
		Password string `json:"password"`
		
	}
	type responseData struct{
		ID uuid.UUID`json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		ChirpyRed bool `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err := decoder.Decode(&reqData)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}
	token,err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"Could not get token",err)
		return
	}
	_,err = auth.ValidateJWT(token, c.JWTSecret)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"Invalid token or user",err)
		return	
	}
	hashPassword,err := auth.HashPassword(reqData.Password)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "Could not hash password",err)
		return 
	}
	
	UpdatedUser,err := c.db.UpdateUserEmailAndPass(r.Context(),database.UpdateUserEmailAndPassParams{
		Email: reqData.Email,
		HashedPassword: hashPassword,
	})
	respondWithJson(w,http.StatusOK, responseData{
		ID: UpdatedUser.ID,
		CreatedAt: UpdatedUser.CreatedAt,
		UpdatedAt: UpdatedUser.UpdatedAt,
		Email: UpdatedUser.Email,
		ChirpyRed: UpdatedUser.IsChirpyRed,
	})

}
