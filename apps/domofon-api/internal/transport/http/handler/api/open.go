package apiRoute

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type openDto struct {
	Code string `json:"code" form:"code" uri:"code" validate:"required"`
}
type resSigninDto struct {
	Success bool `json:"success"`
}

func (h *Route) open(c *gin.Context) {
	var req openDto
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Code != h.config.SecretKey {
		fmt.Printf("code: %s, need: %s\n", req.Code, h.config.SecretKey)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong code"})
		return
	}

	//c.JSON(http.StatusOK, resSigninDto{true})
	//return

	key, err := h.rosdomofon.CreateTemporaryKey(h.config.KeyId)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error on create key"})
		return
	}
	fmt.Printf("key: %s\n", key)

	err = h.rosdomofon.ActivateKey(key)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, resSigninDto{true})

	//var reqData reqSigninDto
	//if err := h.validator.ShouldBindJSON(c, &reqData); err != nil {
	//	httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
	//	return
	//}

	//_, token, err := h.usersService.Sign(reqData.Email, reqData.Password)
	//if err != nil {
	//	err.(*httpError.HTTPError).SendError(c)
	//	return
	//}
	//c.JSON(http.StatusOK, resSigninDto{
	//	Token: "123",
	//})
}
