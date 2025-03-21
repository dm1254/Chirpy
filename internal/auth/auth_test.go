package auth
import(
	"testing"
	"time"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T){
	password1 := "correctPassword123"
	password2 := "anotherPassword123"
	hash1,_:= HashPassword(password1)
	hash2,_ := HashPassword(password2)
	
	tests := []struct{
		name string
		password string
		hash string 
		wantErr bool
	}{
		{
			name : "Correct Password",
			password: password1,
			hash: hash1,
			wantErr: false,
		},
		{
			name: "Incorrect password",
			password: "wrongPassword",
			hash: hash1,
			wantErr: true,
		},
		{
			name: "Password doesnt match hatch",
			password: password1,
			hash: hash2, 
			wantErr: true,
		},
		{
			name: "Empty password",
			password: "",
			hash: hash1,
			wantErr: true,
		},
		{
			name: "invalid hash",
			password: password1,
			hash: "invalidhash",
			wantErr: true,
		},
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T){
			err := ComparePasswordAndHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr{
				t.Errorf("ComparePasswordAndHash() error = %v, wantErr %v", err, tt.wantErr)
			}
	
		})
	
	}
}

func TestJWTAuth(t *testing.T){
	ID1 := uuid.New()
	tokenSecret := "01234"
	expiresIn := (5 * time.Minute)
	newToken,_ := MakeJWT(ID1, tokenSecret, expiresIn)	
	expiredToken,_:= MakeJWT(ID1, tokenSecret, (-5 * time.Minute))
	tests := []struct{
		name string
		tokenString string
		tokenSecret string
		wantErr bool
		
	}{
		{
			name: "Valid token",
			tokenString:newToken,
			tokenSecret: tokenSecret,
			wantErr: false,
		},
		{
			name: "Invalid token",
			tokenString: newToken,
			tokenSecret: "4321",
			wantErr: true,
		},
		{	
			name: "Expired Token",
			tokenString: expiredToken,
			tokenSecret: tokenSecret,
			wantErr: true,
		},
		{
			name: "Empty token string",
			tokenString: "",
			tokenSecret: tokenSecret,
			wantErr: true,
		},

	}
	for _,tt := range tests{	
		t.Run(tt.name, func(t *testing.T){
			_,err := ValidateJWT(tt.tokenString,tt.tokenSecret)
			if (err != nil) != tt.wantErr{
				t.Errorf("ValidateJWT() error = %v, wantErr: %v", err,tt.wantErr)
			}
		})

	}
	

}
