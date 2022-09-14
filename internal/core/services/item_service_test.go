package services_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/osalomon89/test-crud-api/internal/core/domain"
	. "github.com/osalomon89/test-crud-api/internal/core/services"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/server/handler/dto"
	"github.com/osalomon89/test-crud-api/internal/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Item service", func() {
	var itemRepositoryMock *mocks.MockItemRepository
	var itemBody *dto.ItemBody

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		itemRepositoryMock = mocks.NewMockItemRepository(mockCtrl)
		itemBody = &dto.ItemBody{
			Code:        "sa4123",
			Title:       "my-title",
			Description: "my-description",
			Price:       50,
			Stock:       150,
			ItemType:    "SELLER",
			Leader:      true,
			LeaderLevel: domain.LeaderLevelPlatinum,
			Photos: []string{
				"https://http2.mlstatic.com/D_NQ_NP_729539-MLA48049063325_102021-O.jpg",
				"https://http2.mlstatic.com/D_NQ_NP_879745-MLA48049070326_102021-O.jpg",
			},
		}
	})

	Describe("Creating items", func() {
		When("the item data is valid", func() {
			Context("and the repository works", func() {
				It("is saved correctly", func() {
					itemRepositoryMock.EXPECT().
						SaveItem(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1)

					itemService, err := NewItemService(itemRepositoryMock)
					Expect(err).NotTo(HaveOccurred())

					item := itemBody.ToItemDomain()

					result, err := itemService.CreateItem(context.TODO(), item)
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Code).To(Equal(itemBody.Code))
				})
			})
			Context("and the information is not complete", func() {
				It("an error is returned", func() {
					itemService, err := NewItemService(itemRepositoryMock)
					Expect(err).NotTo(HaveOccurred())

					itemBody.Photos = []string{}

					item := itemBody.ToItemDomain()

					result, err := itemService.CreateItem(context.TODO(), item)
					Expect(err).NotTo(BeNil())
					Expect(errors.As(err, new(domain.ItemError))).Should(Equal(true))
					Expect(result).To(BeNil())
				})
			})
		})
	})
})
