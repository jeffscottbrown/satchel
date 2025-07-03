package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jeffscottbrown/gogoogle/secrets"
	"github.com/jeffscottbrown/satchel/model"
	"github.com/jeffscottbrown/satchel/repository"
	"github.com/jeffscottbrown/satchel/utils"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"gorm.io/gorm"
)

type oauthConfig struct {
	clientId     string
	clientSecret string
	callbackUrl  string
}

func login(c *gin.Context) {
	req := c.Request
	res := c.Writer
	if _, err := gothic.CompleteUserAuth(res, req); err == nil {
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func logout(c *gin.Context) {
	req := c.Request
	res := c.Writer
	gothic.Logout(res, req)
	slog.Info("User logged out")
	http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
}

func authCallback(c *gin.Context) {
	req := c.Request
	res := c.Writer
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		slog.Error("Error authenticating user", "error", err)
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if !isAllowedDomain(user.Email) {
		gothic.Logout(res, req)
		c.Redirect(http.StatusFound, "/forbidden")
		return
	}
	slog.Info("User authenticated", "email", user.Email)

	gothic.StoreInSession("authenticatedUser", user.Email, req, res)

	_, err = repository.GetEmployeeByEmail(user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Info("Profile not found in database - new profile being created", "email", user.Email)
			newEmployee := &model.Employee{
				Name:      user.Name,
				Email:     user.Email,
				ImageName: user.AvatarURL,
			}
			newEmployee.AddReflection("Temporary Thing #1", "1")
			newEmployee.AddReflection("Temporary Thing #2", "2")
			newEmployee.AddReflection("Temporary Thing #3", "3")
			newEmployee.AddReflection("Temporary Thing #4", "4")
			if err := repository.SaveEmployee(newEmployee); err != nil {
				slog.Error("Error adding employee", "error", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			slog.Info("New employee added", "email", user.Email)

		} else {
			slog.Error("Error querying employee", "error", err)
			return
		}
	}

	http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
}

func IsAuthenticated(req *http.Request) bool {
	_, err := gothic.GetFromSession("authenticatedUser", req)
	return err == nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Problem loading .env file", "error", err)
	}

	gothic.Store = sessions.NewCookieStore([]byte("dev-secret-don't-use-in-prod"))

	slog.Debug("Configuring authentication providers")

	googleConfig := createOauthConfig("google")

	goth.UseProviders(
		google.New(googleConfig.clientId, googleConfig.clientSecret, googleConfig.callbackUrl, "profile", "email"),
	)
}

func createOauthConfig(provider string) *oauthConfig {
	providerUpperCase := strings.ToUpper(provider)
	providerLowerCase := strings.ToLower(provider)

	callbackUrlVarName := providerUpperCase + "_OAUTH_CALLBACK_URL"
	idVarName := providerUpperCase + "_OAUTH_CLIENT_ID"
	secretVarName := providerUpperCase + "_OAUTH_CLIENT_SECRET"

	callbackUrl := retrieveSecretValue(callbackUrlVarName)
	if callbackUrl == "" {
		callbackUrl = "http://localhost:8080/auth/" + providerLowerCase + "/callback"
	}

	return &oauthConfig{
		callbackUrl:  callbackUrl,
		clientId:     utils.RetrieveSecretValue(idVarName),
		clientSecret: retrieveSecretValue(secretVarName),
	}
}

func retrieveSecretValue(secretName string) string {
	clientSecret, err := secrets.RetrieveSecret(secretName)
	if err != nil {
		slog.Warn("Falling back to OS environment variable", slog.String("secretName", secretName), slog.Any("error", err))
		clientSecret = os.Getenv(secretName)
	}

	return clientSecret
}

func ConfigureAuthorizationHandlers(router *gin.Engine) {
	providerAwareGroup := router.Group("/auth/:provider")

	providerAwareGroup.Use(providerAware())
	providerAwareGroup.GET("/callback", authCallback)
	providerAwareGroup.GET("/login", login)

	router.GET("/auth/logout", logout)
}

// gothic tries a number of techniques to retrieve the provider
// from the URL but none of them are compatible with how
// the gin library provides access to the value
// see https://github.com/markbates/goth/blob/260588e82ba14930ae070a80acadcf0f75348c05/gothic/gothic.go#L263
// this middleware will add the provider to the context in a way that gothic can use

func providerAware() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		provider := c.Param("provider")
		c.Request = req.WithContext(context.WithValue(req.Context(), "provider", provider))

		c.Next()
	}
}
func AuthRequired(c *gin.Context) {

	if !IsAuthenticated(c.Request) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()

}

var allowedDomains = []string{"objectcomputing.com"}

func isAllowedDomain(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	for _, allowed := range allowedDomains {
		if domain == allowed && parts[0] != "" {
			return true
		}
	}
	return false
}
