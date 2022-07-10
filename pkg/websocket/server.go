package websocket
import "fmt"

import (
	"net/http"
	"strconv"

	"github.com/cloudwebrtc/flutter-webrtc-server/pkg/logger"
	"github.com/gorilla/websocket"
)

type WebSocketServerConfig struct {
	Host           string
	Port           int
	CertFile       string
	KeyFile        string
	HTMLRoot       string
	WebSocketPath  string
	TurnServerPath string
}

func DefaultConfig() WebSocketServerConfig {
	return WebSocketServerConfig{
		Host:           "0.0.0.0",
		Port:           8086,
		HTMLRoot:       "web",
		WebSocketPath:  "/ws",
		TurnServerPath: "/api/turn",
	}
}

type WebSocketServer struct {
	handleWebSocket  func(ws *WebSocketConn, request *http.Request)
	handleTurnServer func(writer http.ResponseWriter, request *http.Request)
	// Websocket upgrader
	upgrader websocket.Upgrader
}
//---------------------------------------------------

func NewWebSocketServer(
	wsHandler func(ws *WebSocketConn, request *http.Request),
	turnServerHandler func(writer http.ResponseWriter, request *http.Request)) *WebSocketServer {
	var server = &WebSocketServer{
		handleWebSocket:  wsHandler,
		handleTurnServer: turnServerHandler,
	}
	server.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return server
}
//---------------------------------------------------

func (server *WebSocketServer) handleWebSocketRequest(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("TK: WHAT?server!1");
	logger.Infof("TK: WHAT?server!1")
	responseHeader := http.Header{}
	//responseHeader.Add("Sec-WebSocket-Protocol", "protoo")
	socket, err := server.upgrader.Upgrade(writer, request, responseHeader)
	if err != nil {
		logger.Panicf("%v", err)
	}
	wsTransport := NewWebSocketConn(socket)
	server.handleWebSocket(wsTransport, request)
	wsTransport.ReadMessage()
}

func (server *WebSocketServer) handleTurnServerRequest(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("TK: WHAT?server!2---------", request ,"------------");
	fmt.Println("TK: WHAT?server!2");
	logger.Infof("TK: WHAT?server!2")

	server.handleTurnServer(writer, request)

    //TKADD
    //setupHeader(writer, request)
    //if(request.Method == "OPTIONS"){
    //    writer.WriteHeader(http.StatusOK)
	//    server.handleTurnServer(writer, request)
    //}else{
	//    server.handleTurnServer(writer, request)
    //}
}
//TKADD
//func setupHeader(rw http.ResponseWriter, req *http.Request) {
//        rw.Header().Set("Content-Type", "application/json")
//        rw.Header().Set("Access-Control-Allow-Origin", "*")
//        rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//        rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//}
// Bind .
func (server *WebSocketServer) Bind(cfg WebSocketServerConfig) {
	//func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
	//func Handle(pattern string, handler Handler)
	fmt.Println("TK2: ", http.Dir(cfg.HTMLRoot), " is the stunPort");
	// Websocket handle func
	http.HandleFunc(cfg.WebSocketPath, server.handleWebSocketRequest)
	http.HandleFunc(cfg.TurnServerPath, server.handleTurnServerRequest)

    //TK:TEMP
    //http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	//    fmt.Println("TK: WHAT?server!2---------", req ,"------------");
    //    http.FileServer(http.Dir(cfg.HTMLRoot));
	//})

	//TK: To use the operating system's file system implementation
	http.Handle("/", http.FileServer(http.Dir(cfg.HTMLRoot)))

	logger.Infof("Flutter WebRTC Server listening on: %s:%d", cfg.Host, cfg.Port)
	//TK: panic(http.ListenAndServe(cfg.Host+":"+strconv.Itoa(cfg.Port), nil))

	panic(http.ListenAndServe(cfg.Host+":"+strconv.Itoa(cfg.Port), nil))
	//panic(http.ListenAndServeTLS(cfg.Host+":"+strconv.Itoa(cfg.Port), cfg.CertFile, cfg.KeyFile, nil))
}
