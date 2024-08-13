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

	logo, err := os.ReadFile("../docu-stream/office/test/logo.png")
	assert.Nil(t, err)
	assert.NotEmpty(t, data)

	res, err := c.Apply(&WordApplyReq{
		Docu: data,
		Header: []DocuValue{
			{
				Type:  DocuValueTypeText,
				Key:   "Company Name",
				Value: "OwlByTech",
			},
			{
				Type:  DocuValueTypeImage,
				Key:   "Company Logo",
				Value: &logo,
			},
		},
		Body: []DocuValue{
			{
				Type:  DocuValueTypeText,
				Key:   "Company Name",
				Value: "OwlByTech",
			},
			{
				Type:  DocuValueTypeImage,
				Key:   "Company Logo",
				Value: &logo,
			},
		},
		Footer: []DocuValue{
			{
				Type:  DocuValueTypeText,
				Key:   "Company Name",
				Value: "OwlByTech",
			},
			{
				Type:  DocuValueTypeImage,
				Key:   "Company Logo",
				Value: &logo,
			},
		},
	})

	assert.Nil(t, err)
	assert.NotNil(t, res)

	err = os.WriteFile("./word_test.docx", res.Docu, 0644)
	assert.Nil(t, err)
}
