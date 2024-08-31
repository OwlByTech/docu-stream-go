package docustream

import (
	"os"
	"testing"
	ds "github.com/owlbytech/docu-stream-go"
	"github.com/stretchr/testify/assert"
)

func applyTestWordApplyService(t *testing.T) {
	c, err := ds.NewWordClient(&ds.ConnectOptions{Url: "localhost:3000"})
	assert.Nil(t, err)
	assert.NotNil(t, c)

	data, err := os.ReadFile("../../docu-stream/test/Template.docx")

	assert.Nil(t, err)
	assert.NotEmpty(t, data)

	logo, err := os.ReadFile("../../docu-stream/test/logo.png")
	assert.Nil(t, err)
	assert.NotEmpty(t, data)

	res, err := c.Apply(&ds.WordApplyReq{
		Docu: data,
		Header: []ds.DocuValue{
			{
				Type:  ds.DocuValueTypeText,
				Key:   "Company Name",
				Value: "OwlByTech",
			},
			{
				Type:  ds.DocuValueTypeImage,
				Key:   "Company Logo",
				Value: &logo,
			},
		},
		Body: []ds.DocuValue{
			{
				Type:  ds.DocuValueTypeText,
				Key:   "Company Name",
				Value: "OwlByTech",
			},
			{
				Type:  ds.DocuValueTypeImage,
				Key:   "Company Logo",
				Value: &logo,
			},
		},
		Footer: []ds.DocuValue{
			{
				Type:  ds.DocuValueTypeText,
				Key:   "Company Name",
				Value: "OwlByTech",
			},
			{
				Type:  ds.DocuValueTypeImage,
				Key:   "Company Logo",
				Value: &logo,
			},
		},
	})

	assert.Nil(t, err)
	assert.NotNil(t, res)

	err = os.WriteFile("./outputs/word_test.docx", res.Docu, 0644)
	assert.Nil(t, err)
}
