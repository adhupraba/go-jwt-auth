package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/adhupraba/go-jwt-auth/helpers"
	"github.com/adhupraba/go-jwt-auth/initializers"
	"github.com/adhupraba/go-jwt-auth/models"
)

type signupBody struct {
	Email    string
	Password string
}

type signinBody struct {
	Email    string
	Password string
}

func Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body signupBody

	if err := helpers.BodyParser(r.Body, &body); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Unable to parse body")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	res := initializers.DB.Create(&user)

	if res.Error != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	helpers.RespondWithJson(w, http.StatusCreated, struct{}{})
}

func Signin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body signinBody
	err := helpers.BodyParser(r.Body, &body)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Unable to parse body")
		return
	}

	var user models.User
	initializers.DB.Model(models.User{Email: body.Email}).First(&user)

	if user.ID == 0 {
		helpers.RespondWithError(w, http.StatusNotFound, "User does not exist")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Incorrect email or password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", user.ID),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to create jwt token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenStr,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Now().Add(time.Hour * 1),
	})

	helpers.RespondWithJson(w, http.StatusOK, struct{}{})
}

func Validate(w http.ResponseWriter, r *http.Request, user models.User) {
	helpers.RespondWithJson(w, http.StatusOK, user)
}
