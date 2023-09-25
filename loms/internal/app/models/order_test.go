package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrder_ShouldBeCancelled(t *testing.T) {
	now := time.Now()
	orderTenMinutesAgo := Order{CreatedAt: now.Add(-time.Minute * 10), Status: OrderStatusAwaitingPayment}

	assert.True(t, orderTenMinutesAgo.ShouldBeCancelled(time.Minute*5))
	assert.False(t, orderTenMinutesAgo.ShouldBeCancelled(time.Minute*11))

	orderTenMinutesAgo.Status = OrderStatusNew
	assert.False(t, orderTenMinutesAgo.ShouldBeCancelled(time.Minute*5))
}
