package docustream

import (
	"bytes"
	"context"
	"io"

	pb "github.com/owlbytech/docu-stream-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Convert struct {
	client pb.ConvertClient
}

func NewConvertClient(c *ConnectOptions) (*Convert, error) {
	conn, err := grpc.NewClient(c.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewConvertClient(conn)

	return &Convert{
		client: client,
	}, nil
}

// TODO: Refactor send and recv for any new grpc service
func (c *Convert) WordToPdf(docu *[]byte) (*[]byte, error) {
	var attachFiles []*[]byte = []*[]byte{docu}
	stream, err := c.client.WordToPdf(context.Background())
	if err != nil {
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

		chunkReq := &pb.WordToPdfReq{
				Docu: &pb.DocuChunk{
					Chunks: chunks,
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

	docuOut := bufRes.Bytes()
	return &docuOut, nil
}
