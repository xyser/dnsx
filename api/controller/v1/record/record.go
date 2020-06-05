package record

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	v1 "github.com/dingdayu/dnsx/api/controller/v1"
	"github.com/dingdayu/dnsx/internal/service"
)

// Lists record list api
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

// Create create record api
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
