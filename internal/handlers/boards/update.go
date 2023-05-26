package boards

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *handler) HandleUpdateBoard(ctx *gin.Context) {
	_, ok := common.GetAuthHeader(ctx)
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

	var req public.UpdateBoardRequest
	err = ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, CreateBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	collection := h.mongoDb.Database("photoflux").Collection("boards")
	boardAttr := model.BoardUpdateAttr{
		UpdatedAt: time.Now(),
		Data:      req.Data,
	}
	update := bson.D{{Key: "$set", Value: boardAttr}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		common.EmitError(ctx, CreateBoardError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not update board: %s", err.Error())))
		return
	}

	// TODO: return whole object
	resp := public.UpdateBoardResponse{
		Data: req.Data,
	}

	ctx.JSON(http.StatusOK, &resp)
}
