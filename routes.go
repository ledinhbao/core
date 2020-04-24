package core

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ()

const (
	superUserSecretKey = string("this-is*super+user/secret%key#6p5Zy8qwrnbWPyB@")
)

func initRoutes(e *gin.Engine) {
	e.POST("/superadminregisterfirsttime", superAdminRegisterFirstTime)
}

func superAdminRegisterFirstTime(c *gin.Context) {
	var user User
	var messages []error
	displayError := func(c *gin.Context, code int, message interface{}) {
		c.JSON(code, gin.H{"message": message})
		c.Abort()
	}

	err := c.Bind(&user)
	if err != nil {
		messages = append(messages, fmt.Errorf("Wrong format, %s", err.Error()))
	}

	secretKey := c.PostForm("secret_key")
	if secretKey != superUserSecretKey {
		displayError(c, http.StatusUnauthorized, "SU Secret Key")
		return
	}
	if user.Password != user.PasswordConfirm {
		messages = append(messages, fmt.Errorf("Password don't match"))
	}
	user.SetPassword(user.Password)
	user.Rank = RankSuperAdmin

	// check if there is any super user presents yet.
	var count int
	server.DB().Where("rank=?", RankSuperAdmin).Find(&User{}).Count(&count)
	if count > 0 {
		displayError(c, http.StatusBadRequest, "Super Admin is already registered for the first time. Login to add more.")
		return
	}

	// TODO check for strong password if needed
	messages = append(messages, server.DB().Create(&user).GetErrors()...)
	if len(messages) > 0 {
		displayError(c, http.StatusBadRequest, fmt.Sprintf("%v", messages))
		return
	}

	server.Log.WithFields(logrus.Fields{
		"module": "core",
		"action": "superAdminRegisterFirstTime",
	}).Info("Super Admin created for first time.")

	c.JSON(http.StatusCreated, gin.H{
		"result":  "success",
		"message": "Super User created",
	})
}
