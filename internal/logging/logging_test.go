package logging

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestApplyToContext(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, logrus.TraceLevel)
	ctx := AddToContext(t.Context(), logger)
	assert.Equal(t, logger, FromContext(ctx))
}

func TestFromContext(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, logrus.TraceLevel)
	ctx := AddToContext(t.Context(), logger)
	assert.Equal(t, logger, FromContext(ctx))
}
