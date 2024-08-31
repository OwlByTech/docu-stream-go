package docustream

import (
	"os"
	"testing"

	ds "github.com/owlbytech/docu-stream-go"
	"github.com/stretchr/testify/assert"
)

func applyTestWordToPdf(t *testing.T) {
	c, err := ds.NewConvertClient(&ds.ConnectOptions{Url: "localhost:4014"})
	assert.Nil(t, err)
	assert.NotNil(t, c)

	word, err := os.ReadFile("./outputs/word_test.docx")
	assert.Nil(t, err)
	assert.NotEmpty(t, word)

	pdf, err := c.WordToPdf(&word)
	assert.Nil(t, err)
	assert.NotNil(t, pdf)

	err = os.WriteFile("./outputs/word_test.pdf", *pdf, 0644)
	assert.Nil(t, err)
}
