package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/models"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func createToken(uid string) (string, error) {
    // Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uid,                    
		"iss": "CloudIDE",                
		"exp": time.Now().Add(time.Hour).Unix(), 
		"iat": time.Now().Unix(),                 
	})

	tokenString, err := claims.SignedString(secretKey)
    if err != nil {
        return "", err
    }
  // Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}

func Login(c *fiber.Ctx) error {
	user := new(models.User)

    if err := c.BodyParser(user); err != nil {
        fmt.Println("error = ", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to parse request body",
            "error":   err.Error(),
        })
    }
	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Username, email, and password cannot be empty",
		})
	}
	userDb, err := utils.FindUserByEmail(user.Email)
	if(err != nil){
		if(err==mongo.ErrNoDocuments){
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch email",
				"error":   "User does not exist",
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"error":   err.Error(),
			})
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(user.Password))
	if err != nil {
		return fmt.Errorf("password does not match")
	}
	token, err := createToken(userDb.ID.Hex())
	if(err!=nil){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to create access token",
            "error":   err.Error(),
        })
	}

	// Set the JWT as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour*24),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User Logged in successfully",
		"user":    userDb,
	})

}

func Register(c *fiber.Ctx) error {
	user := new(models.User)

    if err := c.BodyParser(user); err != nil {
        fmt.Println("error = ", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to parse request body",
            "error":   err.Error(),
        })
    }
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Username, email, and password cannot be empty",
		})
	}

	_, err := utils.FindUserByEmail(user.Email)
	if(err != mongo.ErrNoDocuments){
		if(err==nil){
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch email",
				"error":   "User already exists with given email",
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "User already exists with email",
			})
		}
		
	}

	_, err = utils.FindUserByUserName(user.Username)
	if(err != mongo.ErrNoDocuments){
		if(err==nil){
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch email",
				"error":   "User already exists with given username",
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "User already exists with email",
			})
		}
		
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if(err!=nil){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not hash password",
		})
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.LastLogin = time.Now()
	res, err := utils.CreateUser(user)
	if(err!=nil){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to create user",
            "error":   err.Error(),
        })
	}
	user.ID = res
	token, err := createToken(user.ID.Hex())
	if(err!=nil){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to create access token",
            "error":   err.Error(),
        })
	}

	// Set the JWT as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour*24),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User registered successfully",
		"user":    user,
	})
}

func Logout(c *fiber.Ctx) error {


	// Set the JWT as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(0),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})


}
