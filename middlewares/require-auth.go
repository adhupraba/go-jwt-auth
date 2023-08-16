package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/adhupraba/go-jwt-auth/helpers"
	"github.com/adhupraba/go-jwt-auth/initializers"
	"github.com/adhupraba/go-jwt-auth/models"
)

type authedHandler func(http.ResponseWriter, *http.Request, models.User)

func RequireAuth(handler authedHandler) http.HandlerFunc {
	fmt.Println("In middleware")

	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := r.Cookie("Authorization")

		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Not authenticated: %s", err))
			return
		}

		token, err := jwt.Parse(tokenStr.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header)
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unable to parse auth token: %s", err))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid auth token")
			return
		}

		exp, err := claims.GetExpirationTime()

		if err != nil || time.Now().Unix() > exp.Unix() {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Expired auth token")
			return
		}

		sub, err := claims.GetSubject()

		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Subject not found in token")
			return
		}

		userId, err := strconv.Atoi(sub)

		if err != nil || userId == 0 {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Unable to get userId from subject")
			return
		}

		var user models.User
		initializers.DB.Model(models.User{Model: gorm.Model{ID: user.ID}}).First(&user)

		if user.ID == 0 {
			helpers.RespondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		handler(w, r, user)
	}
}
