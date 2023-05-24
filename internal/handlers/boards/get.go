package boards

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *handler) HandleGetBoard(ctx *gin.Context) {
	ah, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	id := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		common.EmitError(ctx, CreateBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind query parameters: %s", err.Error())))
		return
	}

	if err != nil {
		common.EmitError(ctx, GetBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	collection := h.mongoDb.Database("photoflux").Collection("boards")
	filter := bson.D{{Key: "_id", Value: objID}}
	var board model.Board
	err = collection.FindOne(ctx, filter).Decode((&board))
	if err != nil {
		common.EmitError(ctx, ListBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get boards: %s", err.Error())))
		return
	}

	if board.UserId != ah.User.String() {
		common.EmitError(ctx, ListBoardError(
			http.StatusForbidden,
			fmt.Sprintf("Could not access boards, it belongs to a different user")))
		return
	}

	resp := public.GetBoardResponse{
		Data: BoardToItem(board, h.apiPaths),
	}

	ctx.JSON(http.StatusOK, &resp)
}
