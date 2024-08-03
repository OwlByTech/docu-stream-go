package docustream

import (
	"bytes"
	"context"
	"io"

	pb "github.com/owlbytech/docu-stream-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectOptions struct {
	Url string
}

type WordApplyReq struct {
	Docu   []byte
	Body   map[string]string
	Header map[string]string
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

func mapToDocuValues(m map[string]string) []*pb.DocuValues {
	var values []*pb.DocuValues
	for k, v := range m {
		values = append(values, &pb.DocuValues{Key: k, Value: v})
	}

	return values
}

func (w *Word) Apply(req *WordApplyReq) (*WordApplyRes, error) {
	stream, err := w.client.Apply(context.Background())

	if err != nil {
		return nil, err
	}

	initRequest := &pb.WordApplyReq{
		Request: &pb.WordApplyReq_Word{
			Word: &pb.DocuWord{
				Body:   mapToDocuValues(req.Body),
				Header: mapToDocuValues(req.Header),
			},
		},
	}

	if err := stream.Send(initRequest); err != nil {
		return nil, err
	}

	buffer := bytes.NewReader(req.Docu)
	chunkSize := 1024
	buf := make([]byte, chunkSize)

	for {
		n, err := buffer.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n == 0 {
			stream.CloseSend()
			break
		}

		chunkReq := &pb.WordApplyReq{
			Request: &pb.WordApplyReq_Docu{
				Docu: &pb.DocuChunk{
					Chunk: buf[:n],
				},
			},
		}

		if err := stream.Send(chunkReq); err != nil {
			return nil, err
		}
	}

	var docuRes bytes.Buffer
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		chunk := res.Docu.Chunk
		if len(chunk) <= 0 {
			continue
		}

		if _, err := docuRes.Write(chunk); err != nil {
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
