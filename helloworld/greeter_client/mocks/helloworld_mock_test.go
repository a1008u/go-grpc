/*

THE PURPOSE
client側のテスト<br />
 - clientで呼び出すサーバのmockを作成
 - mockを利用してテストを作成

CREATED BY
  a1008u

CREATED AT
  2020/05/02
*/
package mocks

import (
	"context"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestMockStartANDExec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Clientテスト用のモックを作成する
	moclProdInfoClient := NewMockGreeterClient(ctrl)

	//req := &pb.Product{Name: name, Description: description, Price: price}

	word := "test1"
	wordresult := "test"

	// unaryのmock
	moclProdInfoClient.
		EXPECT().SayHello(gomock.Any(), &pb.HelloRequest{Name: word}).
		Return(&pb.HelloReply{Message: "Hello " + wordresult}, nil)

	testSayHello(t, moclProdInfoClient)
}

func testSayHello(t *testing.T, client pb.GreeterClient) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: "test1"})
	if err != nil || r.Message != "Hello test" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply : ", r.GetMessage())


}
