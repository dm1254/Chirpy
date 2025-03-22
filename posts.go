package main 
import(	
	"net/http"
	"encoding/json"
	"time"
	"log"
	"sort"
	"github.com/google/uuid"
	"workspace/github.com/dm1254/Chirpy/internal/database"
	"workspace/github.com/dm1254/Chirpy/internal/auth"
	"errors"
)


func(c *ApiConfig) handlerChirps (w http.ResponseWriter, r *http.Request){
	type requestData struct{
		Body string `json:"body"`
	}
	type responseData struct{
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`

	}		
	decoder := json.NewDecoder(r.Body)
	reqData := requestData{}
	err := decoder.Decode(&reqData)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}
	getToken,err := auth.GetBearerToken(r.Header)

	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "Error getting token from header",err)
		return 
	}
	validUser,err := auth.ValidateJWT(getToken,c.JWTSecret)	
	if err != nil{
		respondWithError(w,http.StatusUnauthorized, "Couldnt validate user", err)
		return
	}
	if len(reqData.Body) <= 140{
		params := database.CreatePostsParams{
			Body: reqData.Body,
			UserID: validUser,
		}
		postsData,err := c.db.CreatePosts(r.Context(), params)
		if err != nil{
			log.Printf("Couldn't create post: %s", err)
			return 
		}
		respondWithJson(w,http.StatusCreated, responseData{
			ID: postsData.ID,
			CreatedAt: postsData.CreatedAt,
			UpdatedAt: postsData.UpdatedAt,
			Body: postsData.Body,
			UserID: postsData.UserID,


		})
		return
	}
	err = errors.New("validation failed: post too long")
	respondWithError(w,http.StatusBadRequest,"Post cannot exceed 140 characters" ,err )
	return 
}

func (c *ApiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request){
	type responseData struct{
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`

	}
	QueryParamSort := r.URL.Query().Get("sort")
	if QueryParamSort != ""{	
		posts,err := c.db.GetPosts(r.Context())
		if err != nil{
			log.Printf("Error creating posts: %s",err)
		}
		allPosts := []responseData{}
		for _, post := range posts{
			postsData := responseData{
				ID: post.ID,
				CreatedAt: post.CreatedAt,
				UpdatedAt:post.UpdatedAt,
				Body: post.Body,
				UserID:post.UserID,

			}
			allPosts = append(allPosts, postsData)
		}
		if QueryParamSort == "desc"{
			sort.Slice(allPosts, func(i,j int) bool {return allPosts[i].CreatedAt.After(allPosts[j].CreatedAt)})
			respondWithJson(w, http.StatusOK, allPosts)
			return 
		}else if QueryParamSort == "asc"{
			respondWithJson(w,http.StatusOK,allPosts)
			return
		
		}
		return 
	}
	QueryParam := r.URL.Query().Get("author_id")
	if QueryParam != ""{
		userId,err := uuid.Parse(QueryParam)
		if err != nil{
			errors.New("Could not parse string to uuid")
			return 
		}
		postsByAuthor,err := c.db.GetPostsByAuthor(r.Context(),userId)
		if err != nil{
			respondWithError(w, http.StatusNotFound, "User not found",err)
			return
		}
		allPosts := []responseData{}
		for _, post := range postsByAuthor{
			postData := responseData{
				ID: post.ID,
				CreatedAt: post.CreatedAt,
				UpdatedAt: post.UpdatedAt,
				Body: post.Body,
				UserID: post.UserID,
			}
			allPosts = append(allPosts, postData)

		}
		respondWithJson(w, http.StatusOK, allPosts)
		return 
	}
	posts,err := c.db.GetPosts(r.Context())
	if err != nil{
		log.Printf("Error creating posts: %s",err)
	}
	allPosts := []responseData{}
	for _, post := range posts{
		postsData := responseData{
			ID: post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt:post.UpdatedAt,
			Body: post.Body,
			UserID:post.UserID,

		}
		allPosts = append(allPosts, postsData)
	}
	respondWithJson(w,http.StatusOK, allPosts)
	return 

}

func(c *ApiConfig) handlerGetIdChirp(w http.ResponseWriter, r *http.Request){
	postURLId := r.PathValue("chirpsID")
	postID,err := uuid.Parse(postURLId)
	log.Printf("postID:%s\n",postURLId)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "could not convert post id to uuid", err)
		return 
	}
	type responseData struct{
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`

	}
	post,err := c.db.GetSinglePost(r.Context(),postID)
	if err != nil{
		respondWithError(w, http.StatusNotFound,"Post does not exist", err)
		return 
	}
	respondWithJson(w, http.StatusOK, responseData{
		ID: post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body: post.Body,
		UserID: post.UserID,

	})
	return
}

func(c *ApiConfig) handlerDeletePost(w http.ResponseWriter, r *http.Request){
	token,err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Token invalid",err)
		return
	}
	userID,err := auth.ValidateJWT(token,c.JWTSecret)
	if err != nil{
		respondWithError(w, http.StatusForbidden, "Invalid user",err)
		return
	}
	postURLId := r.PathValue("chirpsID")
	postID,err := uuid.Parse(postURLId)
	log.Printf("postID:%s\n",postURLId)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "could not convert post id to uuid", err)
		return 
	}

	post,err := c.db.GetSinglePost(r.Context(),postID)
	if post.UserID != userID{
		respondWithError(w, http.StatusForbidden, "Forbidden action",err)
		return
	}
	err = c.db.DeletePost(r.Context(),postID)
	if err != nil{
		respondWithError(w, http.StatusForbidden, "Forbidden action",err)	
		return
	}
	w.WriteHeader(204)
	return

}
