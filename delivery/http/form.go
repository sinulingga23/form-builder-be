package delivery

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinulingga23/form-builder-be/api/usecase"
	"github.com/sinulingga23/form-builder-be/define"
	"github.com/sinulingga23/form-builder-be/payload"
)

type formHttp struct {
	mFormUsecase usecase.IMFormUsecase
}

func NewFormHttp(
	mFormUsecase usecase.IMFormUsecase,
) formHttp {
	return formHttp{mFormUsecase: mFormUsecase}
}

func (delivery *formHttp) ServeHandler(r *gin.RouterGroup) {
	r.POST("/api/v1/forms", delivery.HandleAddForm)
}

func (delivery *formHttp) HandleAddForm(r *gin.Context) {
	data := struct {
		Data payload.CreateMFormRequest `json:"data"`
	}{}

	if errBind := r.Bind(&data); errBind != nil {
		log.Printf("errBind: %v", errBind)
		r.JSON(http.StatusBadRequest, payload.Response{
			StatusCode: http.StatusBadRequest,
			Message:    define.ErrFailedBind.Error(),
		})
		return
	}

	response := delivery.mFormUsecase.AddFrom(r.Request.Context(), data.Data)
	r.JSON(response.StatusCode, response)
	return
}
