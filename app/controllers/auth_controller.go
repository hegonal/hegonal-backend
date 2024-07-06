package controllers

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
	"github.com/hegonal/hegonal-backend/platform/database"
	"github.com/mileusna/useragent"
)

func UserSignUp(c *fiber.Ctx) error {
	signUp := &models.SignUp{}

	if err := c.BodyParser(signUp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(signUp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	user := &models.User{}

	emailBeenUsed, err := db.IsEmailUsed(signUp.Email)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if emailBeenUsed {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "This email already been used.",
		})
	}

	ownerAccountBeenCreated, err := db.IsOwnerAccountCreated()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	user.UserID = utils.GenerateId()
	user.Name = signUp.Name
	user.Password = utils.GeneratePassword(signUp.Password)
	user.Email = signUp.Email
	user.CreatedAt = utils.TimeNow()
	user.UpdatedAt = utils.TimeNow()

	if ownerAccountBeenCreated {
		user.Role = models.HegonalUser
	} else {
		user.Role = models.HegonalOwner
	}

	if err := db.CreateNewUser(user); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	session := &models.Session{}

	ua := useragent.Parse(c.Get("User-Agent"))

	session.UserID = user.UserID
	session.ExpiryTime = utils.TimeNow().Add(24 * time.Hour)
	session.Ip = c.IP()
	session.Device = c.Get(ua.OS + " " + ua.Name)
	session.CreatedAt = utils.TimeNow()
	session.UpdatedAt = utils.TimeNow()
	session.Session, err = utils.GenerateSessionString()

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := db.CreateNewSession(session); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	setSessionCookie(c, user.UserID, session.Session)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
		"tokens": fiber.Map{
			"session": session.Session,
		},
	})
}

func UserLogin(c *fiber.Ctx) error {
	login := &models.Login{}

	if err := c.BodyParser(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	user, err := db.GetUser(login.Email)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid email or password",
		})
	} else if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if !utils.ComparePasswords(user.Password, login.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid email or password",
		})
	}

	session := &models.Session{}

	ua := useragent.Parse(c.Get("User-Agent"))

	session.UserID = user.UserID
	session.ExpiryTime = utils.TimeNow().Add(24 * time.Hour)
	session.Ip = c.IP()
	session.Device = ua.OS + " " + ua.Name
	session.CreatedAt = utils.TimeNow()
	session.UpdatedAt = utils.TimeNow()
	session.Session, err = utils.GenerateSessionString()

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := db.CreateNewSession(session); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	setSessionCookie(c, user.UserID, session.Session)

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
		"tokens": fiber.Map{
			"session": session.Session,
		},
	})
}

func UserLogout(c *fiber.Ctx) error {
	session := c.Cookies("session")
	userId := c.Cookies("userID")

	db, err := database.OpenDBConnection()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	err = db.DeleteSession(userId, session)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to delete session",
		})
	}

	utils.ClearCookies(c, "session", "userID")

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Logout successful",
	})
}

func setSessionCookie(c *fiber.Ctx,userID, session string) {
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    session,
		Expires:  utils.TimeNow().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "userID",
		Value:    userID,
		Expires:  utils.TimeNow().Add(24 * time.Hour * 365),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
}