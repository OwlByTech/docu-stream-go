package docustream

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordApplyService(t *testing.T) {
	c, err := NewWordClient(&ConnectOptions{Url: "localhost:3000"})
	assert.Nil(t, err)
	assert.NotNil(t, c)

	data, err := os.ReadFile("../docu-stream/office/test/Template.docx")

	assert.Nil(t, err)
	assert.NotEmpty(t, data)

	res, err := c.Apply(&WordApplyReq{
		Docu: data,
		Header: map[string]string{
			"Company Name": "OwlByTech",
		},
		Body: map[string]string{
			"Company Name": "OwlByTech",
		},
	})

	assert.Nil(t, err)
	assert.NotNil(t, res)

	err = os.WriteFile("./word_test.docx", res.Docu, 0644)
	assert.Nil(t, err)
}
