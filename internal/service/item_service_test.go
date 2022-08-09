package service_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	itemerrors "github.com/osalomon89/test-crud-api/internal/errors"
	"github.com/osalomon89/test-crud-api/internal/model"
	httpmodel "github.com/osalomon89/test-crud-api/internal/model/http"
	"github.com/osalomon89/test-crud-api/internal/server/httpcontext"
	"github.com/osalomon89/test-crud-api/internal/service"
	"github.com/osalomon89/test-crud-api/internal/utils/test/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_itemService_CreateItem(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	type args struct {
		ctx           context.Context
		crateItemBody httpmodel.CreateItemBody
	}

	type repositoryArgs struct {
		err   error
		times int
	}

	tests := []struct {
		name           string
		args           args
		repositoryArgs repositoryArgs
		want           *model.Item
		wantedErr      error
	}{
		{
			name: "Should create an item successfully",
			args: args{
				ctx: httpcontext.BackgroundFromContext(c),
				crateItemBody: httpmodel.CreateItemBody{
					Code:        "SAM27324354",
					Title:       "Tablet Samsung Galaxy Tab S7",
					Description: "Galaxy Tab S7 with S Pen SM-t733 12.4 pulgadas y 4GB de memoria RAM",
					Price:       150000,
					Stock:       35,
					Photos: []string{
						"https://http2.mlstatic.com/D_NQ_NP_729539-MLA48049063325_102021-O.jpg",
						"https://http2.mlstatic.com/D_NQ_NP_879745-MLA48049070326_102021-O.jpg",
					},
					ItemType:    "SELLER",
					Leader:      true,
					LeaderLevel: "PLATINUM",
				},
			},
			repositoryArgs: repositoryArgs{
				times: 1,
				err:   nil,
			},
			want: &model.Item{
				Model:       gorm.Model{ID: 1},
				Code:        "SAM27324354",
				Title:       "Tablet Samsung Galaxy Tab S7",
				Description: "Galaxy Tab S7 with S Pen SM-t733 12.4 pulgadas y 4GB de memoria RAM",
				Price:       150000,
				Stock:       35,
				Photos: []model.Photo{{
					Path: "https://http2.mlstatic.com/D_NQ_NP_729539-MLA48049063325_102021-O.jpg",
				}, {
					Path: "https://http2.mlstatic.com/D_NQ_NP_879745-MLA48049070326_102021-O.jpg",
				}},
				ItemType:    "SELLER",
				Leader:      true,
				LeaderLevel: "PLATINUM",
				Status:      "ACTIVE",
			},
		},
		{
			name: "Should return an error when the body is not complete",
			args: args{
				ctx: httpcontext.BackgroundFromContext(c),
				crateItemBody: httpmodel.CreateItemBody{
					Code:        "SAM27324354",
					Title:       "Tablet Samsung Galaxy Tab S7",
					Description: "Galaxy Tab S7 with S Pen SM-t733 12.4 pulgadas y 4GB de memoria RAM",
					Price:       150000,
					Stock:       35,
					Photos: []string{
						"https://http2.mlstatic.com/D_NQ_NP_729539-MLA48049063325_102021-O.jpg",
						"https://http2.mlstatic.com/D_NQ_NP_879745-MLA48049070326_102021-O.jpg",
					},
					ItemType:    "SELLER",
					Leader:      true,
					LeaderLevel: "",
				},
			},
			repositoryArgs: repositoryArgs{
				times: 0,
				err:   nil,
			},
			want: nil,
			wantedErr: itemerrors.ItemError{
				Message: "Error in params validation. Leader level is not valid: ",
			},
		},
		{
			name: "Should return an error when the repository returns an error",
			args: args{
				ctx: httpcontext.BackgroundFromContext(c),
				crateItemBody: httpmodel.CreateItemBody{
					Code:        "SAM27324354",
					Title:       "Tablet Samsung Galaxy Tab S7",
					Description: "Galaxy Tab S7 with S Pen SM-t733 12.4 pulgadas y 4GB de memoria RAM",
					Price:       150000,
					Stock:       35,
					Photos: []string{
						"https://http2.mlstatic.com/D_NQ_NP_729539-MLA48049063325_102021-O.jpg",
						"https://http2.mlstatic.com/D_NQ_NP_879745-MLA48049070326_102021-O.jpg",
					},
					ItemType:    "SELLER",
					Leader:      true,
					LeaderLevel: "PLATINUM",
				},
			},
			repositoryArgs: repositoryArgs{
				times: 1,
				err:   errors.New("the repository error"),
			},
			want:      nil,
			wantedErr: fmt.Errorf("error in repository: %w", errors.New("the repository error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			repositoryMock := mocks.NewMockItemRepository(mockCtrl)

			repositoryMock.EXPECT().
				CreateItem(tt.args.ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, item *model.Item) error {
					item.ID = 1
					return tt.repositoryArgs.err
				}).
				Times(tt.repositoryArgs.times)

			srv := service.NewItemService(repositoryMock)
			got, err := srv.CreateItem(tt.args.ctx, tt.args.crateItemBody)
			if tt.wantedErr != nil {
				assert.Equal(t, tt.wantedErr, err, "Error is not the expected.")
				return
			}

			assert.Equal(t, tt.want, got, "Result is not the expected.")
			assert.Nil(t, err, "Unexpected error occurred saving item in DB")
		})
	}
}

func Test_itemService_GetItemByID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	type args struct {
		ctx    context.Context
		itemID uint
	}

	type repositoryArgs struct {
		item  *model.Item
		err   error
		times int
	}

	tests := []struct {
		name           string
		args           args
		repositoryArgs repositoryArgs
		want           *model.Item
		wantedErr      error
	}{
		{
			name: "Should get an item successfully",
			args: args{
				ctx:    httpcontext.BackgroundFromContext(c),
				itemID: 1,
			},
			repositoryArgs: repositoryArgs{
				item: &model.Item{
					Model:       gorm.Model{ID: 1},
					Code:        "SAM27324354",
					Title:       "Tablet Samsung Galaxy Tab S7",
					Description: "Galaxy Tab S7 with S Pen SM-t733 12.4 pulgadas y 4GB de memoria RAM",
					Price:       150000,
					Stock:       35,
					Photos: []model.Photo{{
						Path: "https://http2.mlstatic.com/D_NQ_NP_729539-MLA48049063325_102021-O.jpg",
					}, {
						Path: "https://http2.mlstatic.com/D_NQ_NP_879745-MLA48049070326_102021-O.jpg",
					}},
					ItemType:    "SELLER",
					Leader:      true,
					LeaderLevel: "PLATINUM",
					Status:      "ACTIVE",
				},
				times: 1,
				err:   nil,
			},
			want: &model.Item{
				Model:       gorm.Model{ID: 1},
				Code:        "SAM27324354",
				Title:       "Tablet Samsung Galaxy Tab S7",
				Description: "Galaxy Tab S7 with S Pen SM-t733 12.4 pulgadas y 4GB de memoria RAM",
				Price:       150000,
				Stock:       35,
				Photos: []model.Photo{{
					Path: "https://http2.mlstatic.com/D_NQ_NP_729539-MLA48049063325_102021-O.jpg",
				}, {
					Path: "https://http2.mlstatic.com/D_NQ_NP_879745-MLA48049070326_102021-O.jpg",
				}},
				ItemType:    "SELLER",
				Leader:      true,
				LeaderLevel: "PLATINUM",
				Status:      "ACTIVE",
			},
		},
		{
			name: "Should return an error when the repository returns an error",
			args: args{
				ctx:    httpcontext.BackgroundFromContext(c),
				itemID: 1,
			},
			repositoryArgs: repositoryArgs{
				item:  nil,
				times: 1,
				err:   errors.New("the repository error"),
			},
			wantedErr: fmt.Errorf("error in repository: %w", errors.New("the repository error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			repositoryMock := mocks.NewMockItemRepository(mockCtrl)

			repositoryMock.EXPECT().
				GetItemByID(tt.args.ctx, gomock.Any()).
				Return(tt.repositoryArgs.item, tt.repositoryArgs.err).
				Times(tt.repositoryArgs.times)

			srv := service.NewItemService(repositoryMock)
			got, err := srv.GetItemByID(tt.args.ctx, tt.args.itemID)
			if tt.wantedErr != nil {
				assert.Equal(t, tt.wantedErr, err, "Error is not the expected.")
				return
			}

			assert.Equal(t, tt.want, got, "Result is not the expected.")
			assert.Nil(t, err, "Unexpected error occurred getting item from DB")
		})
	}
}
