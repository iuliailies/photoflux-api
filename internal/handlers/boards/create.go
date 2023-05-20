package boards

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func (h *handler) HandleCreateBoard(ctx *gin.Context) {
	ah, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	var req public.CreateBoardRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, CreateBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	collection := h.mongoDb.Database("photoflux").Collection("boards")
	board := model.Board{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    ah.User.String(),
		Data:      req.Data,
	}
	_, err = collection.InsertOne(ctx, board)
	if err != nil {
		common.EmitError(ctx, CreateBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create board: %s", err.Error())))
		return
	}

	resp := public.CreateBoardResponse{
		Data: BoardToItem(board, h.apiPaths),
	}

	ctx.JSON(http.StatusOK, &resp)
}
