package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/libsv/go-dpp"
	dppMocks "github.com/libsv/go-dpp/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPaymentHandler_CreatedPayment(t *testing.T) {
	tests := map[string]struct {
		paymentCreateFunc func(context.Context, dpp.PaymentCreateArgs, dpp.Payment) (*dpp.PaymentACK, error)
		reqBody           dpp.Payment
		paymentID         string
		expResponse       dpp.PaymentACK
		expStatusCode     int
		expErr            error
	}{
		"successful post": {
			paymentCreateFunc: func(ctx context.Context, args dpp.PaymentCreateArgs, req dpp.Payment) (*dpp.PaymentACK, error) {
				return &dpp.PaymentACK{
					Memo: fmt.Sprintf("payment %s", args.PaymentID),
				}, nil
			},
			paymentID: "abc123",
			reqBody:   dpp.Payment{},
			expResponse: dpp.PaymentACK{
				Memo: "payment abc123",
			},
			expStatusCode: http.StatusCreated,
		},
		"error response returns 422": {
			paymentCreateFunc: func(ctx context.Context, args dpp.PaymentCreateArgs, req dpp.Payment) (*dpp.PaymentACK, error) {
				return &dpp.PaymentACK{
					Memo:  "failed",
					Error: 1,
				}, nil
			},
			paymentID: "abc123",
			reqBody:   dpp.Payment{},
			expResponse: dpp.PaymentACK{
				Error: 1,
				Memo:  "failed",
			},
			expStatusCode: http.StatusUnprocessableEntity,
		},
		"payment create service error is handled": {
			paymentCreateFunc: func(ctx context.Context, args dpp.PaymentCreateArgs, req dpp.Payment) (*dpp.PaymentACK, error) {
				return nil, errors.New("ohnonono")
			},
			paymentID:     "abc123",
			reqBody:       dpp.Payment{},
			expStatusCode: http.StatusInternalServerError,
			expErr:        errors.New("ohnonono"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			h := NewPaymentHandler(&dppMocks.PaymentServiceMock{
				PaymentCreateFunc: test.paymentCreateFunc,
			})
			g := e.Group("/")
			e.HideBanner = true
			h.RegisterRoutes(g)

			body, err := json.Marshal(test.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/v1/payment/:paymentID")
			ctx.SetParamNames("paymentID")
			ctx.SetParamValues(test.paymentID)

			err = h.createPayment(ctx)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, test.expErr, err.Error())
				return
			}

			response := rec.Result()
			defer response.Body.Close()
			assert.Equal(t, test.expStatusCode, response.StatusCode)

			var ack dpp.PaymentACK
			assert.NoError(t, json.NewDecoder(response.Body).Decode(&ack))

			assert.Equal(t, test.expResponse, ack)
		})
	}
}
