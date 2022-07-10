package gpnrpc

import (
	"encoding/gob"
	"fmt"
	"gpnrpc/type"

	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int
	CodecType   _type.Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   _type.GobType,
}

// RPC Server.
type Server struct{}


func NewServer() *Server {
	return &Server{}
}


var DefaultServer = NewServer()


func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var opt Option
	if err := gob.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	f := _type.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	server.serveCodec(f(conn))
}


var invalidRequest = struct{}{}

func (server *Server) serveCodec(cc _type.Codec) {
	mylock := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, mylock)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, mylock, wg)
	}
	wg.Wait()
	_ = cc.Close()
}


type request struct {
	h            *_type.Header
	argv, replyv reflect.Value
}

func (server *Server) readRequestHeader(cc _type.Codec) (*_type.Header, error) {
	var h _type.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc _type.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (server *Server)  sendResponse(cc _type.Codec, h *_type.Header, body interface{}, myLock *sync.Mutex) {
	myLock.Lock()
	defer myLock.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc _type.Codec, req *request, mylock *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("gpnrpc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), mylock)
}

func (server *Server) accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}

func Accept(lis net.Listener) { DefaultServer.accept(lis) }
