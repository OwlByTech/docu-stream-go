package docustream

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strconv"

	pb "github.com/owlbytech/docu-stream-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectOptions struct {
	Url string
}

type DocuValueType pb.DocuValueType

const DocuValueTypeText DocuValueType = 0
const DocuValueTypeImage DocuValueType = 1

type DocuValueRaw = interface{}

type DocuValue struct {
	Type  DocuValueType
	Key   string
	Value DocuValueRaw
}

type WordApplyReq struct {
	Docu   []byte
	Body   []DocuValue
	Header []DocuValue
}

type WordApplyRes struct {
	Docu []byte
}

type Word struct {
	client pb.WordClient
}

func NewWordClient(c *ConnectOptions) (*Word, error) {
	conn, err := grpc.NewClient(c.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewWordClient(conn)

	return &Word{
		client: client,
	}, nil
}

func processData(attachFiles *[]*[]byte, m []DocuValue) ([]*pb.DocuValue, error) {
	var values []*pb.DocuValue

	for _, v := range m {
		kind := reflect.ValueOf(v.Value).Kind()

		switch kind {
		case reflect.String:
			values = append(values, &pb.DocuValue{Key: v.Key, Value: v.Value.(string), Type: pb.DocuValueType_TEXT})
		case reflect.Ptr:
			attachFile, ok := v.Value.(*[]byte)
			if !ok {
				return nil, fmt.Errorf("Slice must be of an pointer of slice of byte types")
			}

			*attachFiles = append(*attachFiles, attachFile)
			values = append(values, &pb.DocuValue{Key: v.Key, Value: strconv.Itoa(len(*attachFiles)-1), Type: pb.DocuValueType_IMAGE})
		default:
			return nil, fmt.Errorf("Docu stream only support string and pointer slice of bytes")
		}
	}

	return values, nil
}

func (w *Word) Apply(req *WordApplyReq) (*WordApplyRes, error) {
	var attachFiles []*[]byte = []*[]byte{&req.Docu}

	bodyValues, err := processData(&attachFiles, req.Body)
	if err != nil {
		return nil, err
	}

	headerValues, err := processData(&attachFiles, req.Header)
	if err != nil {
		return nil, err
	}

	initRequest := &pb.WordApplyReq{
		Request: &pb.WordApplyReq_Word{
			Word: &pb.DocuWord{
				Body:   bodyValues,
				Header: headerValues,
			},
		},
	}

	stream, err := w.client.Apply(context.Background())
	if err != nil {
		return nil, err
	}

	if err := stream.Send(initRequest); err != nil {
		return nil, err
	}

	var attachFilesReaders []*bytes.Reader

	for _, file := range attachFiles {
		attachFilesReaders = append(attachFilesReaders, bytes.NewReader(*file))
	}

	for {
		allFileDone := true

		chunks := make([][]byte, len(attachFiles))
		for k, fileReader := range attachFilesReaders {
			chunkSize := 1024
			buf := make([]byte, chunkSize)

			n, err := fileReader.Read(buf)
			if err != nil && err != io.EOF {
				return nil, err
			}

			if n > 0 {
				allFileDone = false
				chunks[k] = buf[:n]
			}
		}

		chunkReq := &pb.WordApplyReq{
			Request: &pb.WordApplyReq_Docu{
				Docu: &pb.DocuChunk{
					Chunks: chunks,
				},
			},
		}

		if err := stream.Send(chunkReq); err != nil {
			return nil, err
		}

		if allFileDone {
			break
		}
	}

	stream.CloseSend()

	var docuRes bytes.Buffer
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		chunk := res.Docu.Chunks
		if len(chunk) <= 0 {
			continue
		}

		if _, err := docuRes.Write(chunk[0]); err != nil {
			return nil, err
		}
	}

	bufRes := new(bytes.Buffer)
	_, err = bufRes.ReadFrom(&docuRes)
	if err != nil {
		return nil, err
	}

	return &WordApplyRes{Docu: bufRes.Bytes()}, nil
}
