package converter

import (
	"testing"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func TestToRepoModelPaymentMethod(t *testing.T) {
	tests := []struct {
		name               string
		paymentMethod      model.OrderPaymentMethod
		expectedRepoMethod repoModel.OrderPaymentMethod
	}{
		{
			name:               "Test CARD payment method",
			paymentMethod:      model.OrderPaymentMethod_CARD,
			expectedRepoMethod: repoModel.OrderPaymentMethod_CARD,
		},
		{
			name:               "Test CREDIT_CARD payment method",
			paymentMethod:      model.OrderPaymentMethod_CREDIT_CARD,
			expectedRepoMethod: repoModel.OrderPaymentMethod_CREDIT_CARD,
		},
		{
			name:               "Test SBP payment method",
			paymentMethod:      model.OrderPaymentMethod_SBP,
			expectedRepoMethod: repoModel.OrderPaymentMethod_SBP,
		},
		{
			name:               "Test INVESTOR_MONEY payment method",
			paymentMethod:      model.OrderPaymentMethod_INVESTOR_MONEY,
			expectedRepoMethod: repoModel.OrderPaymentMethod_INVESTOR_MONEY,
		},
		{
			name:               "Test unknown payment method",
			paymentMethod:      model.OrderPaymentMethod(100),
			expectedRepoMethod: repoModel.OrderPaymentMethod_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMethod := ToRepoModelPaymentMethod(tt.paymentMethod)
			if repoMethod != tt.expectedRepoMethod {
				t.Errorf("ToRepoModelPaymentMethod() = %v, want %v", repoMethod, tt.expectedRepoMethod)
			}
		})
	}
}

func TestToModelPaymentMethod(t *testing.T) {
	tests := []struct {
		name                string
		paymentMethod       repoModel.OrderPaymentMethod
		expectedModelMethod model.OrderPaymentMethod
	}{
		{
			name:                "Test CARD payment method",
			paymentMethod:       repoModel.OrderPaymentMethod_CARD,
			expectedModelMethod: model.OrderPaymentMethod_CARD,
		},
		{
			name:                "Test CREDIT_CARD payment method",
			paymentMethod:       repoModel.OrderPaymentMethod_CREDIT_CARD,
			expectedModelMethod: model.OrderPaymentMethod_CREDIT_CARD,
		},
		{
			name:                "Test SBP payment method",
			paymentMethod:       repoModel.OrderPaymentMethod_SBP,
			expectedModelMethod: model.OrderPaymentMethod_SBP,
		},
		{
			name:                "Test INVESTOR_MONEY payment method",
			paymentMethod:       repoModel.OrderPaymentMethod_INVESTOR_MONEY,
			expectedModelMethod: model.OrderPaymentMethod_INVESTOR_MONEY,
		},
		{
			name:                "Test unknown payment method",
			paymentMethod:       repoModel.OrderPaymentMethod(100),
			expectedModelMethod: model.OrderPaymentMethod_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelMethod := ToModelPaymentMethod(tt.paymentMethod)
			if modelMethod != tt.expectedModelMethod {
				t.Errorf("ToModelPaymentMethod() = %v, want %v", modelMethod, tt.expectedModelMethod)
			}
		})
	}
}

func TestToModelOrderStatus(t *testing.T) {
	tests := []struct {
		name                string
		status              repoModel.OrderStatus
		expectedModelStatus model.OrderStatus
	}{
		{
			name:                "Test CANCELED status",
			status:              repoModel.OrderStatus_CANCELED,
			expectedModelStatus: model.OrderStatus_CANCELED,
		},
		{
			name:                "Test PAID status",
			status:              repoModel.OrderStatus_PAID,
			expectedModelStatus: model.OrderStatus_PAID,
		},
		{
			name:                "Test default (PENDING_PAYMENT) status",
			status:              repoModel.OrderStatus_PENDING_PAYMENT,
			expectedModelStatus: model.OrderStatus_PENDING_PAYMENT,
		},
		{
			name:                "Test unknown status",
			status:              repoModel.OrderStatus(100),
			expectedModelStatus: model.OrderStatus_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelStatus := ToModelOrderStatus(tt.status)
			if modelStatus != tt.expectedModelStatus {
				t.Errorf("ToModelOrderStatus() = %v, want %v", modelStatus, tt.expectedModelStatus)
			}
		})
	}
}

func TestToRepoModelOrderStatus(t *testing.T) {
	tests := []struct {
		name               string
		status             model.OrderStatus
		expectedRepoStatus repoModel.OrderStatus
	}{
		{
			name:               "Test PAID status conversion",
			status:             model.OrderStatus_PAID,
			expectedRepoStatus: repoModel.OrderStatus_PAID,
		},
		{
			name:               "Test CANCELED status conversion",
			status:             model.OrderStatus_CANCELED,
			expectedRepoStatus: repoModel.OrderStatus_CANCELED,
		},
		{
			name:               "Test default (PENDING_PAYMENT) status conversion",
			status:             model.OrderStatus_PENDING_PAYMENT,
			expectedRepoStatus: repoModel.OrderStatus_PENDING_PAYMENT,
		},
		{
			name:               "Test unknown status conversion",
			status:             model.OrderStatus(100),
			expectedRepoStatus: repoModel.OrderStatus_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoStatus := ToRepoModelOrderStatus(tt.status)
			if repoStatus != tt.expectedRepoStatus {
				t.Errorf("ToRepoModelOrderStatus() = %v, want %v", repoStatus, tt.expectedRepoStatus)
			}
		})
	}
}
