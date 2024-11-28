package domain_test

import (
	"app/internal/quote/domain"
	mockDomain "app/internal/quote/domain/mock"
	"app/internal/quote/types"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testUnitQuote struct {
	service       *domain.Quote
	repository    *mockDomain.MockquoteRepository
	taxClient     *mockDomain.MocktaxClient
	catalogClient *mockDomain.MockcatalogClient
}

func newTestUnitQuote(ctrl *gomock.Controller) *testUnitQuote {
	repository := mockDomain.NewMockquoteRepository(ctrl)
	taxClient := mockDomain.NewMocktaxClient(ctrl)
	catalogClient := mockDomain.NewMockcatalogClient(ctrl)

	return &testUnitQuote{
		repository:    repository,
		taxClient:     taxClient,
		catalogClient: catalogClient,
		service:       domain.NewQuote(repository, catalogClient, taxClient),
	}
}

func assertQuoteEqual(t *testing.T, expected, actual *types.Quote) {
	assert.Equal(t, expected.Address, actual.Address)
	assert.Equal(t, expected.Payment, actual.Payment)
	assert.Equal(t, expected.Products, actual.Products)

	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.Amount, actual.Amount)
	assert.Equal(t, expected.TaxAmount, actual.TaxAmount)
	assert.Equal(t, expected.TotalAmount, actual.TotalAmount)
	assert.NotEmpty(t, actual.UUID.String())
}

func TestQuoteLoadDraftByCustomerNotFound(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tc := newTestUnitQuote(ctrl)
	ctx := context.Background()
	customerUUID := uuid.New()

	tc.repository.EXPECT().
		FindByCustomerAndStatus(gomock.Any(), gomock.Eq(customerUUID), gomock.Eq(types.QuoteStatusDraft)).
		Return(nil, types.ErrQuoteNotFound)

	// act
	quote, err := tc.service.LoadDraftByCustomer(ctx, customerUUID)

	// assert
	assert.NoError(t, err)
	assertQuoteEqual(t, types.NewQuote(uuid.New(), customerUUID), quote)
}

func TestQuoteLoadDraftByCustomerExists(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tc := newTestUnitQuote(ctrl)
	ctx := context.Background()
	customerUUID := uuid.New()
	expected := types.NewQuote(uuid.New(), customerUUID)

	tc.repository.EXPECT().
		FindByCustomerAndStatus(gomock.Any(), gomock.Eq(customerUUID), gomock.Eq(types.QuoteStatusDraft)).
		Return(expected, nil)

	// act
	actual, err := tc.service.LoadDraftByCustomer(ctx, customerUUID)

	// assert
	assert.NoError(t, err)
	assertQuoteEqual(t, expected, actual)
}

func TestQuoteAddProduct(t *testing.T) {
}

func TestQuoteAddProductRecalculation(t *testing.T) {
}

func TestQuoteAddProductFailed(t *testing.T) {
}

func TestQuoteUpdateProduct(t *testing.T) {
}

func TestQuoteUpdateProductRecalculation(t *testing.T) {
}

func TestQuoteUpdateProductFailed(t *testing.T) {
}

func TestQuoteProcessByCustomerID(t *testing.T) {
}

func TestQuoteProcessByCustomerIDFailed(t *testing.T) {
}

func TestQuoteRemoveProduct(t *testing.T) {
}

func TestQuoteRemoveProductRecalculation(t *testing.T) {
}

func TestQuoteRemoveProductFailed(t *testing.T) {
}

func TestQuoteSaveAddress(t *testing.T) {
}

func TestQuoteSaveAddressFailed(t *testing.T) {
}

func TestQuoteSavePayment(t *testing.T) {
}

func TestQuoteSavePaymentFailed(t *testing.T) {
}
