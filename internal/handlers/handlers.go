package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Kroning/x-test/internal/database"
	"github.com/Kroning/x-test/internal/entities"
	"github.com/Kroning/x-test/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type App struct {
	Db        database.AbstractDB
	Logger    logger.Logger
	JWTSecret []byte
}

var (
	AuthHeaderMissing    = "Authorization header missing"
	TokenIncorrectFormat = "Incorrectly formatted authorization header"
)

func (app *App) CreateCompany(c *gin.Context) {
	app.Logger.Info("CreateCompany")
	var req CompanyRequest
	if err := c.BindJSON(&req); err != nil {
		app.Logger.Info("CreateCompany", "error", err)
		SendRawResponse(c, http.StatusBadRequest, CreateErrorResponse(err.Error(), http.StatusBadRequest))
		c.Abort()
		return
	}

	app.Logger.Info("CreateCompany", "request", req)
	company := entities.NewCompany(uuid.Nil, req.Name, req.Description, req.AmountOfEmployees, req.Registered, req.Type)
	err := company.Validate()
	if err != nil {
		app.Logger.Info("CreateCompany", "error", err)
		SendRawResponse(c, http.StatusBadRequest, CreateErrorResponse(err.Error(), http.StatusBadRequest))
		c.Abort()
		return
	}

	err = app.Db.CreateCompany(company)
	if err != nil {
		app.Logger.Info("CreateCompany", "error", err)
		SendRawResponse(c, http.StatusInternalServerError, CreateErrorResponse(err.Error(), http.StatusInternalServerError))
		c.Abort()
		return
	}

	status := ResponseStatus{ErrorCode: 200, Message: "Company created"}
	resp := CompanyResponse{Status: status, Company: *company}
	respBody, err := json.Marshal(resp)
	if err != nil {
		respBody = JSONErrorResponse()
	}

	app.Logger.Info("CreateCompany", "created", company.Id)
	SendRawResponse(c, http.StatusOK, respBody)
}

func (app *App) PatchCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		app.Logger.Info("PatchCompany", "error", err)
		SendRawResponse(c, http.StatusBadRequest, CreateErrorResponse(err.Error(), http.StatusBadRequest))
		c.Abort()
		return
	}

	app.Logger.Info("PatchCompany", "id", c.Param("id"))

	var req CompanyRequest
	if err := c.BindJSON(&req); err != nil {
		app.Logger.Info("PatchCompany", "error", err)
		SendRawResponse(c, http.StatusBadRequest, CreateErrorResponse(err.Error(), http.StatusBadRequest))
		c.Abort()
		return
	}

	app.Logger.Info("PatchCompany", "request", req)
	company := entities.NewCompany(id, req.Name, req.Description, req.AmountOfEmployees, req.Registered, req.Type)
	err = company.Validate()
	if err != nil {
		app.Logger.Info("PatchCompany", "error", err)
		SendRawResponse(c, http.StatusBadRequest, CreateErrorResponse(err.Error(), http.StatusBadRequest))
		c.Abort()
		return
	}

	err = app.Db.PatchCompany(company)
	if err != nil {
		app.Logger.Info("PatchCompany", "error", err)
		SendRawResponse(c, http.StatusInternalServerError, CreateErrorResponse(err.Error(), http.StatusInternalServerError))
		c.Abort()
		return
	}

	status := ResponseStatus{ErrorCode: 200, Message: "Company patched"}
	resp := CompanyResponse{Status: status, Company: *company}
	respBody, err := json.Marshal(resp)
	if err != nil {
		respBody = JSONErrorResponse()
	}

	app.Logger.Info("PatchCompany", "changed", company.Id)
	SendRawResponse(c, http.StatusOK, respBody)
}

func (app *App) DeleteCompany(c *gin.Context) {
	idString := c.Param("id")
	_, err := uuid.Parse(idString)
	if err != nil {
		app.Logger.Info("DeleteCompany", "error", err)
		SendRawResponse(c, http.StatusBadRequest, CreateErrorResponse(err.Error(), http.StatusBadRequest))
		c.Abort()
		return
	}

	app.Logger.Info("DeleteCompany", "id", idString)

	err = app.Db.DeleteCompany(idString)
	if err != nil {
		app.Logger.Info("DeleteCompany", "error", err)
		SendRawResponse(c, http.StatusInternalServerError, CreateErrorResponse(err.Error(), http.StatusInternalServerError))
		c.Abort()
		return
	}

	status := ResponseStatusWithID{ErrorCode: 200, Message: "Company deleted", Id: idString}
	respBody, err := json.Marshal(status)
	if err != nil {
		respBody = JSONErrorResponse()
	}

	app.Logger.Info("DeleteCompany", "deleted", idString)
	SendRawResponse(c, http.StatusOK, respBody)
}

func (app *App) GetCompany(c *gin.Context) {
	app.Logger.Info("GetCompany", "id", c.Param("id"))
	idString := c.Param("id")
	_, err := uuid.Parse(idString)
	if err != nil {
		app.Logger.Info("GetCompany", "error", err)
		SendRawResponse(c, http.StatusBadRequest, CreateErrorResponse(err.Error(), http.StatusBadRequest))
		c.Abort()
		return
	}

	app.Logger.Info("GetCompany", "id", idString)

	company, err := app.Db.GetCompany(idString)
	if err != nil {
		app.Logger.Info("GetCompany", "error", err)
		SendRawResponse(c, http.StatusInternalServerError, CreateErrorResponse(err.Error(), http.StatusInternalServerError))
		c.Abort()
		return
	}

	status := ResponseStatus{ErrorCode: 200, Message: "Company found"}
	resp := CompanyResponse{Status: status, Company: *company}
	respBody, err := json.Marshal(resp)
	if err != nil {
		respBody = JSONErrorResponse()
	}

	app.Logger.Info("GetCompany", "found", company.Id)
	SendRawResponse(c, http.StatusOK, respBody)
}

func (app *App) CheckAuth(ctx *gin.Context) {
	app.Logger.Debug("CheckAuth")

	// Get JWT
	fullToken := ctx.Request.Header.Get("Authorization")
	if fullToken == "" {
		app.Logger.Info("CheckAuth", "error", AuthHeaderMissing)
		SendRawResponse(ctx, http.StatusUnauthorized, CreateErrorResponse(AuthHeaderMissing, http.StatusUnauthorized))
		ctx.Abort()
		return
	}
	splitToken := strings.Split(fullToken, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		app.Logger.Info("CheckAuth", "error", TokenIncorrectFormat, "fullToken", fullToken, "token", splitToken[0])
		SendRawResponse(ctx, http.StatusUnauthorized, CreateErrorResponse(TokenIncorrectFormat, http.StatusUnauthorized))
		ctx.Abort()
		return
	}
	jwtToken := splitToken[1]

	// Check JWT
	token, err := app.verifyJWT(jwtToken)
	if err != nil {
		app.Logger.Info("CheckAuth", "error", err)
		SendRawResponse(ctx, http.StatusUnauthorized, CreateErrorResponse(err.Error(), http.StatusUnauthorized))
		ctx.Abort()
		return
	}
	app.Logger.Info("CheckAuth token is good", "Claims", token.Claims)
}

func (app *App) verifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return app.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// We also can test some claims like expiration
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		app.Logger.Info("verifyJWT", "claims", claims)
	} else {
		return nil, fmt.Errorf("invalida token")
	}

	return token, nil
}

func SendRawResponse(ctx *gin.Context, httpCode int, respBody []byte) {
	ctx.Data(httpCode, "application/json; charset=utf-8", respBody)
}

func CreateErrorResponse(message string, errCode int) []byte {
	var status ResponseStatus
	status.Message = message
	status.ErrorCode = errCode

	respBody, err := json.Marshal(status)
	if err != nil {
		respBody = JSONErrorResponse()
	}

	return respBody
}

func JSONErrorResponse() []byte {
	return []byte(`{status:{"http_code": 500, "error_code":"", "message":"JSON marshal error"}, data: {}}`)
}
