package eltype

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// PluginServer 插件通信服务器
type PluginServer struct {
	pluginReader *PluginReader
	Addr         *string
	Upgrader     websocket.Upgrader
}

func NewPluginServer(pluginReader *PluginReader) (*PluginServer, error) {
	server := new(PluginServer)
	server.pluginReader = pluginReader
	server.Addr = flag.String("addr", "localhost:8080", "http service address")
	server.Upgrader = websocket.Upgrader{}
	return server, nil
}

func (server *PluginServer) fetchEvent(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	pluginName := r.URL.Query().Get("name")
	if !server.pluginReader.randKeySet[key] || pluginName == "" {
		fmt.Fprint(w, "{\"code\":0, \"msg\":\"Wrong key or name\"}")
		return
	}
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	log.Printf("[Info] The plugin {%s} is mounted", r.URL.Query().Get("name"))
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (server *PluginServer) sendMessage(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	pluginName := r.URL.Query().Get("name")
	if !server.pluginReader.randKeySet[key] || pluginName == "" {
		fmt.Fprint(w, "{\"code\":0, \"msg\":\"Wrong key or name\"}")
		return
	}
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	log.Printf("[Info] The plugin {%s} is mounted", r.URL.Query().Get("name"))
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (server *PluginServer) sendOperation(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	pluginName := r.URL.Query().Get("name")
	if !server.pluginReader.randKeySet[key] || pluginName == "" {
		fmt.Fprint(w, "{\"code\":0, \"msg\":\"Wrong key or name\"}")
		return
	}
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	log.Printf("[Info] The plugin {%s} is mounted", r.URL.Query().Get("name"))
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (server *PluginServer) Start() {
	go server.startServer()
}

func (server *PluginServer) startServer() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/fetchEvent", server.fetchEvent)
	http.HandleFunc("/sendMessage", server.sendMessage)
	http.HandleFunc("/sendOperation", server.sendOperation)
	log.Fatal(http.ListenAndServe(*server.Addr, nil))
}
