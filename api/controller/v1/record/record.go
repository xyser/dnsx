package record

import (
	v1 "dnsx/api/controller/v1"
	"dnsx/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

func Lists(c *gin.Context) {
	name := c.Query("name")
	qtype := c.Query("type")
	value := c.Query("value")

	rrs, err := service.GetRecordList(name, qtype, value)
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	c.JSON(http.StatusOK, v1.NewSucResponse(rrs))
}

func Create(c *gin.Context) {
	var args = CreateArgs{}
	if err := c.Bind(&args); err != nil {
		if errors.As(err, &validator.ValidationErrors{}) {
			c.JSON(http.StatusOK, v1.NewErrMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusOK, v1.ErrFailParams)
		return
	}

	rr, err := service.CreateRecord(args.Name, args.Type, args.Value, args.TTL, args.Priority)
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	c.JSON(http.StatusOK, v1.NewSucResponse(rr))
}
