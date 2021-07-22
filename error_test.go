package rerpc

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/akshayjshah/rerpc/internal/assert"
)

func TestErrorCodeOK(t *testing.T) {
	assert.Nil(t, Errorf(CodeOK, "everything is fine"), "error with code ok")
}

func TestErrorFormatting(t *testing.T) {
	assert.Equal(t, Errorf(CodeUnavailable, "").Error(), CodeUnavailable.String(), "no message")
	text := Errorf(CodeUnavailable, "foo").Error()
	assert.True(t, strings.Contains(text, CodeUnavailable.String()), "error text should include code")
	assert.True(t, strings.Contains(text, "foo"), "error text should include message")
}

func TestErrorCode(t *testing.T) {
	err := fmt.Errorf("another: %w", Errorf(CodeUnavailable, "foo"))
	rerr, ok := AsError(err)
	assert.True(t, ok, "extract rerpc error")
	assert.Equal(t, rerr.Code(), CodeUnavailable, "extracted code")
}

func TestCodeOf(t *testing.T) {
	assert.Equal(t, CodeOf(nil), CodeOK, "nil error code")
	assert.Equal(t, CodeOf(Errorf(CodeUnavailable, "foo")), CodeUnavailable, "explicitly-set code")
	assert.Equal(t, CodeOf(errors.New("foo")), CodeUnknown, "fallback code")
}

func TestErrorDetails(t *testing.T) {
	second := durationpb.New(time.Second)
	detail, err := anypb.New(second)
	assert.Nil(t, err, "create anypb.Any")
	rerr := Errorf(CodeUnknown, "details").(*Error)
	assert.Zero(t, rerr.Details(), "fresh error")
	assert.Nil(t, rerr.AddDetail(second), "add detail")
	assert.Equal(t, rerr.Details(), []*anypb.Any{detail}, "retrieve details")
	assert.Nil(t, rerr.SetDetails(second, second), "overwrite details")
	assert.Equal(t, rerr.Details(), []*anypb.Any{detail, detail}, "retrieve details")
}
