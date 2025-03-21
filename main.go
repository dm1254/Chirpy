package main 
import (
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"workspace/github.com/dm1254/Chirpy/internal/database"
	"net/http"
	"database/sql"
	"sync/atomic"
	"log"
	"os"
	"fmt"
)
type ApiConfig struct{
	fileserverhits atomic.Int32
	db *database.Queries
	Platform string
	JWTSecret string
}

func (c *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		c.fileserverhits.Add(1)
		

		next.ServeHTTP(w,r)
	})	
}
func main(){
	const port = "8080"
	const filepathRoot = "."
	
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
		
	dbconn,err := sql.Open("postgres",dbURL)
	if err != nil{
		log.Fatalf("Error connecting to database")
	}
	platform := os.Getenv("PLATFORM")
	JWTsecret := os.Getenv("JWTSECRET")
	dbQueries := database.New(dbconn)
	c := &ApiConfig{
		fileserverhits: atomic.Int32{},
		db: dbQueries,
		Platform: platform,
		JWTSecret: JWTsecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", c.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))

	})
	
	mux.HandleFunc("GET /admin/metrics", c.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", c.handlerReset)
	mux.HandleFunc("POST /api/validate", c.HandlerValidate)
	mux.HandleFunc("POST /api/users", c.handlerUsers)
	mux.HandleFunc("POST /api/chirps",c.handlerChirps)
	mux.HandleFunc("GET /api/chirps", c.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpsID}",c.handlerGetIdChirp)
	mux.HandleFunc("POST /api/login", c.handlerLogin)
	mux.HandleFunc("POST /api/refresh",c.handleRefresh)
	mux.HandleFunc("POST /api/revoke", c.handleRevoke)
	mux.HandleFunc("PUT /api/users", c.handlerUpdateUser)
	s := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}
	log.Printf("Serving files from %s on port: %s\n",filepathRoot,port)
	log.Fatal(s.ListenAndServe())
}
func (c *ApiConfig) handlerMetrics (w http.ResponseWriter, r *http.Request){
	hits := c.fileserverhits.Load()
	html := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)	
	w.WriteHeader(200)
	w.Write([]byte(html))

}

func(c *ApiConfig) handlerReset (w http.ResponseWriter,r *http.Request){
	c.fileserverhits.Store(0)
	if  c.Platform != "dev"{
		w.WriteHeader(http.StatusForbidden)
		return 
	}		
	err := c.db.Reset(r.Context())
	
	if err != nil{
		log.Printf("Error deleting users: %s",err)
	}
	w.WriteHeader(200)
}


