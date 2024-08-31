package docustream

import (
	"testing"
)

func TestWordToPdf(t *testing.T) {
	t.Run("Generate docx", func(t *testing.T) {
		applyTestWordApplyService(t)
	})

	t.Run("Convert docx to pdf", func(t *testing.T) {
		applyTestWordToPdf(t)
	})
}
