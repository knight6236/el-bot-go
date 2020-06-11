package eltype

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// PluginServer 插件通信服务器
type PluginServer struct {
	mapMute             sync.RWMutex
	pluginReader        *PluginReader
	Addr                *string
	ReceivedEvent       chan Event
	WillBeSentMessage   chan Message
	WillBeSentOperation chan Operation
	WillBeSentControl   chan Control
	MsgQueue            map[string]chan []byte
	Upgrader            websocket.Upgrader
}

func NewPluginServer(pluginReader *PluginReader) (*PluginServer, error) {
	server := new(PluginServer)
	server.pluginReader = pluginReader
	server.Addr = flag.String("addr", "127.0.0.1:9999", "http service address")
	server.Upgrader = websocket.Upgrader{}
	server.MsgQueue = make(map[string]chan []byte)
	server.WillBeSentMessage = make(chan Message, 64)
	server.WillBeSentOperation = make(chan Operation, 64)
	server.WillBeSentControl = make(chan Control, 64)
	server.ReceivedEvent = make(chan Event, 64)
	for key, _ := range pluginReader.randKeySet {
		server.MsgQueue[key] = make(chan []byte, 64)
	}
	server.MsgQueue["key"] = make(chan []byte, 64)
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
	for {
		message := <-server.MsgQueue[key]
		err = c.WriteMessage(1, message)
		if err != nil {
			log.Printf("[Error] send event to plugin {%s}: %s\n", pluginName, err.Error())
			break
		}
		log.Printf("[Info] send event to plugin {%s}: Success\n", pluginName)
	}
}

func (server *PluginServer) sendMessage(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	pluginName := server.pluginReader.PluginMap[key].Name
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
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("[Error] receive message from plugin {%s}: %s\n", pluginName, err.Error())
			break
		}
		log.Printf("[Info] receive message from plugin {%s}: Success\n", pluginName)
		var message Message
		json.Unmarshal(msg, &message)
		server.WillBeSentMessage <- message
	}
}

func (server *PluginServer) sendOperation(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	pluginName := server.pluginReader.PluginMap[key].Name
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
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("[Error] receive operation from plugin {%s}: %s\n", pluginName, err.Error())
			break
		}
		log.Printf("[Info] receive operation from plugin {%s}: Success\n", pluginName)
		var operation Operation
		json.Unmarshal(msg, &operation)
		server.WillBeSentOperation <- operation
	}
}

func (server *PluginServer) sendControl(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	pluginName := server.pluginReader.PluginMap[key].Name
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
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("[Error] receive control from plugin {%s}: %s\n", pluginName, err.Error())
			break
		}
		log.Printf("[Info] receive control from plugin {%s}: Success\n", pluginName)
		var control Control
		json.Unmarshal(msg, &control)
		server.WillBeSentControl <- control
	}
}

func (server *PluginServer) Start() {
	go server.startServer()
	go server.listenEvent()
	time.Sleep(time.Duration(2) * time.Second)
	go server.startPlugin()
}

func (server *PluginServer) startServer() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/fetchEvent", server.fetchEvent)
	http.HandleFunc("/sendMessage", server.sendMessage)
	http.HandleFunc("/sendOperation", server.sendOperation)
	http.HandleFunc("/sendControl", server.sendControl)
	log.Println(http.ListenAndServe(*server.Addr, nil))
}

func (server *PluginServer) startPlugin() {
	for key, plugin := range server.pluginReader.PluginMap {
		if !plugin.IsProcMsg {
			continue
		}
		var err error
		switch plugin.Type {
		case PluginTypeBinary:
			if runtime.GOOS == "windows" {
				err = Exec(plugin.Path, key)
			} else {
				err = Exec("/bin/bash", "-c", fmt.Sprintf("%s %s", plugin.Path, key))
			}
		case PluginTypeJava:
			if runtime.GOOS == "windows" {
				err = Exec("java", "-jar", fmt.Sprintf("%s %s", plugin.Path, key))
			} else {
				err = Exec("/bin/bash", "-c", fmt.Sprintf("java -jar %s %s", plugin.Path, key))
			}
		case PluginTypePython:
			if runtime.GOOS == "windows" {
				err = Exec("python", plugin.Path, fmt.Sprintf("%s %s", plugin.Path, key))
			} else {
				err = Exec("/bin/bash", "-c", fmt.Sprintf("%s %s %s", PythonCommand, plugin.Path, key))
			}
		case PluginTypeJavaScript:
			if runtime.GOOS == "windows" {
				err = Exec("node", plugin.Path, fmt.Sprintf("%s %s", plugin.Path, key))
			} else {
				err = Exec("/bin/bash", "-c", fmt.Sprintf("node %s %s", plugin.Path, key))
			}
		}
		if err != nil {
			log.Printf("[Error] Exec plugin {%s}: %s\n", plugin.Path, err.Error())
		} else {
			log.Printf("[Info] Exec plugin {%s}: Success\n", plugin.Path)
		}
	}
}

func (server *PluginServer) listenEvent() {
	for true {
		select {
		case event := <-server.ReceivedEvent:
			bytes, err := json.Marshal(event)
			if err != nil {

			} else {
				for _, ch := range server.MsgQueue {
					ch <- bytes
				}
			}
		}
	}
}