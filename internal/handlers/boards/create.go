package boards

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	boardAttr := model.BoardAttr{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    ah.User.String(),
		Data:      req.Data,
	}
	res, err := collection.InsertOne(ctx, boardAttr)
	if err != nil {
		common.EmitError(ctx, CreateBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create board: %s", err.Error())))
		return
	}

	board := model.Board{
		Id: res.InsertedID.(primitive.ObjectID),
		// ??????
		CreatedAt: boardAttr.CreatedAt,
		UpdatedAt: boardAttr.UpdatedAt,
		UserId:    ah.User.String(),
		Data:      req.Data,
	}

	resp := public.CreateBoardResponse{
		Data: BoardToItem(board, h.apiPaths),
	}

	ctx.JSON(http.StatusOK, &resp)
}
