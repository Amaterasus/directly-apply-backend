package models

import (
	"os"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/jinzhu/gorm"
	// This is required for using postgres with gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	ID string `gorm:"primaryKey;type:uuid"`
	Name string
	Email string
	PhoneNumber string
	HashedPassword string
	FoundUs string
	SendJobMatches bool
	AgreedToTerms bool
	Jwt string `gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// InitialUserMigration will use GORM to migrate the tables in the database.
func InitialUserMigration() {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DataBase")
	}

	defer db.Close()

	db.AutoMigrate(&User{})
}

func (user *User) Authorise(name, password string) bool {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DataBase")
	}

	defer db.Close()

	db.Where("name = ?", name).Find(&user)
	jwt, _ := GenerateJWT(user.ID)

	user.Jwt = jwt

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	return err == nil 
}

func GenerateJWT(id string) (string, error) {

	secret := os.Getenv("SECRET")

    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = id
    claims["expiration"] = time.Now().Add(time.Minute * 30).Unix()

    tokenString, err := token.SignedString([]byte(secret))

    if err != nil {
        fmt.Println(err)
        return "", err
    }

    return tokenString, nil
}

func DecodeJWT(token string) string {
	
	decodedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC); 
		if !ok {
			return nil, fmt.Errorf("There was an error")
		}
		secret := os.Getenv("SECRET")
		return []byte(secret), nil
	})
	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok && decodedToken.Valid {
		id := fmt.Sprintf("%v", claims["id"])
		return id
	} else {
		return ""
	}
}

func (user *User) GetAllUsers() *[]User {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DataBase")
	}

	defer db.Close()

	users := []User{}

	db.Find(&users)

	return &users
}

func (user *User) FindUserByID(id string) {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DataBase")
	}

	defer db.Close()

	db.Where("id = ?", id).Find(&user)
}

func (user *User) Create(name, email, phoneNumber, password, foundUs string, sendJobMatches, agreedToTerms bool) interface{} {

	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DataBase")
	}

	defer db.Close()

	newUserID := uuid.New()

	finalID := newUserID.String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		panic(err)
	}

	jwt, _ := GenerateJWT(finalID)

	newUser := db.Create(&User{ID: finalID, Name: name, Email: email, PhoneNumber: phoneNumber, HashedPassword: string(hashedPassword), FoundUs: foundUs, Jwt: jwt, SendJobMatches: sendJobMatches, AgreedToTerms: agreedToTerms})

	return newUser.Value
}
