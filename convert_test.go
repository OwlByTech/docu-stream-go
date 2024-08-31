package docustream

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordToPdf(t *testing.T){
	c, err := NewConvertClient(&ConnectOptions{Url: "localhost:4014"})
	assert.Nil(t, err)
	assert.NotNil(t, c)

	word, err := os.ReadFile("../docu-stream/test/Template.docx")
	assert.Nil(t, err)
	assert.NotEmpty(t, word)

	pdf, err := c.WordToPdf(&word)
	assert.Nil(t, err)
	assert.NotNil(t, pdf)

	err = os.WriteFile("./word_test.pdf", *pdf, 0644)
	assert.Nil(t, err)
}
