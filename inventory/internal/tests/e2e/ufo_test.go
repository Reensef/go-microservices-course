package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
	inventoryProto "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryProto.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "expected successful connection to gRPC server")

		inventoryClient = inventoryProto.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.MongoClearParts(ctx)
		Expect(err).ToNot(HaveOccurred(), "expected to successfully clear parts collection")

		cancel()
	})

	Describe("GetPart", func() {
		It("must return an existing part from mongo", func() {
			part := &repoModel.Part{
				Info: repoModel.PartInfo{
					Name:          gofakeit.ProductName(),
					Description:   gofakeit.Sentence(10),
					Price:         55.55,
					StockQuantity: 123,
					Category:      repoModel.PartCategory_ENGINE,
					Dimensions: &repoModel.PartDimensions{
						Width:  1,
						Height: 2,
						Length: 3,
					},
					Manufacturer: &repoModel.PartManufacturer{
						Name:    "Manufacturer",
						Country: "Country",
						Website: "Website",
					},
					Tags: []string{"Tag1", "Tag2", "Tag3"},
				},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			}

			id, err := env.MongoInsertPart(ctx, part)
			if err != nil {
				Expect(err).ToNot(HaveOccurred())
			}

			resp, err := inventoryClient.GetPart(ctx, &inventoryProto.GetPartRequest{
				Id: id,
			})

			Expect(err).ToNot(HaveOccurred())

			Expect(resp.GetPart().GetId()).To(Equal(id))
			Expect(resp.GetPart().GetName()).To(Equal(part.Info.Name))
			Expect(resp.GetPart().GetDescription()).To(Equal(part.Info.Description))
			Expect(resp.GetPart().GetPrice()).To(Equal(part.Info.Price))
			Expect(resp.GetPart().GetStockQuantity()).To(Equal(part.Info.StockQuantity))
			Expect(resp.GetPart().GetCategory()).To(Equal(inventoryProto.Category_CATEGORY_ENGINE))

			Expect(resp.GetPart().GetDimensions().GetHeight()).
				To(Equal(part.Info.Dimensions.Height))
			Expect(resp.GetPart().GetDimensions().GetWidth()).
				To(Equal(part.Info.Dimensions.Width))
			Expect(resp.GetPart().GetDimensions().GetWeight()).
				To(Equal(part.Info.Dimensions.Weight))
			Expect(resp.GetPart().GetDimensions().GetLength()).
				To(Equal(part.Info.Dimensions.Length))

			Expect(resp.GetPart().GetManufacturer().GetCountry()).
				To(Equal(part.Info.Manufacturer.Country))
			Expect(resp.GetPart().GetManufacturer().GetWebsite()).
				To(Equal(part.Info.Manufacturer.Website))
			Expect(resp.GetPart().GetManufacturer().GetName()).
				To(Equal(part.Info.Manufacturer.Name))

			Expect(resp.GetPart().GetTags()).To(Equal(part.Info.Tags))
		})
	})
})
