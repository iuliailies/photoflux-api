package categories

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func (h *handler) HandleListCategories(ctx *gin.Context) {

	var categories []model.Category
	err := h.db.WithContext(ctx).Order("id asc").Find(&categories).Error

	if err != nil {
		common.EmitError(ctx, ListCategoriesError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not list categories: %s", err.Error())))
		return
	}

	resp := public.ListCategoryResponse{
		Data: make([]public.CategoryListItemData, 0, len(categories)),
		Links: public.ListCategoryLinks{
			Self: h.apiPaths.Categories + "/",
		},
	}
	for _, c := range categories {
		resp.Data = append(resp.Data, CategoryToPublicListItem(c, h.apiPaths))
	}
	ctx.JSON(http.StatusOK, &resp)
}
