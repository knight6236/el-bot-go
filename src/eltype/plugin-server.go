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
	ConnMap             map[string][]*websocket.Conn
	MsgQueue            map[string]chan []byte
	AliveMap            map[string]bool
	Upgrader            websocket.Upgrader
}

func NewPluginServer(pluginReader *PluginReader) (*PluginServer, error) {
	server := new(PluginServer)
	server.pluginReader = pluginReader
	server.Addr = flag.String("addr", "127.0.0.1:9999", "http service address")
	server.Upgrader = websocket.Upgrader{}
	server.MsgQueue = make(map[string]chan []byte)
	server.AliveMap = make(map[string]bool)
	server.ConnMap = make(map[string][]*websocket.Conn)
	server.WillBeSentMessage = make(chan Message, 64)
	server.WillBeSentOperation = make(chan Operation, 64)
	server.WillBeSentControl = make(chan Control, 64)
	server.ReceivedEvent = make(chan Event, 64)
	pluginReader.randKeySet["test"] = true
	for key, _ := range pluginReader.randKeySet {
		server.MsgQueue[key] = make(chan []byte, 64)
		// server.AliveMap[key] = false
	}
	return server, nil
}

func (server *PluginServer) fetchEvent(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	pluginName := server.pluginReader.PluginMap[key].Name
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	if !server.pluginReader.randKeySet[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"Wrong key\"}"))
		return
	}
	server.mapMute.RLock()
	if !server.AliveMap[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"no heartbeat connection established\"}"))
		server.mapMute.RUnlock()
		return
	}
	server.ConnMap[key] = append(server.ConnMap[key], c)
	server.mapMute.RUnlock()
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
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	if !server.pluginReader.randKeySet[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"Wrong key\"}"))
		return
	}
	server.mapMute.RLock()
	if !server.AliveMap[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"no heartbeat connection established\"}"))
		server.mapMute.RUnlock()
		return
	}
	server.ConnMap[key] = append(server.ConnMap[key], c)
	server.mapMute.RUnlock()
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
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	if !server.pluginReader.randKeySet[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"Wrong key\"}"))
		return
	}
	server.mapMute.RLock()
	if !server.AliveMap[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"no heartbeat connection established\"}"))
		server.mapMute.RUnlock()
		return
	}
	server.ConnMap[key] = append(server.ConnMap[key], c)
	server.mapMute.RUnlock()
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
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	if !server.pluginReader.randKeySet[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"Wrong key\"}"))
		return
	}
	server.mapMute.RLock()
	if !server.AliveMap[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"no heartbeat connection established\"}"))
		server.mapMute.RUnlock()
		return
	}
	server.ConnMap[key] = append(server.ConnMap[key], c)
	server.mapMute.RUnlock()
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

func (server *PluginServer) receiveHeartbeat(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	pluginName := server.pluginReader.PluginMap[key].Name
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	if !server.pluginReader.randKeySet[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"Wrong key\"}"))
		return
	}
	server.mapMute.RLock()
	server.ConnMap[key] = append(server.ConnMap[key], c)
	server.mapMute.RUnlock()
	defer c.Close()
	for {
		err = c.WriteMessage(1, []byte("Alive"))
		if err != nil {
			log.Printf("[Error] send heartbeat to plugin {%s}: %s\n", pluginName, err.Error())
			break
		}
		log.Printf("[Info] send heartbeat to plugin {%s}: Success\n", pluginName)
		time.Sleep(time.Duration(15) * time.Second)
	}
}

func (server *PluginServer) sendHeartbeat(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	pluginName := server.pluginReader.PluginMap[key].Name
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	if !server.pluginReader.randKeySet[key] {
		c.WriteMessage(1, []byte("{\"code\":1, \"msg\":\"Wrong key\"}"))
		return
	}
	defer c.Close()
	server.mapMute.Lock()
	server.AliveMap[key] = true
	server.mapMute.Unlock()
	log.Printf("[Info] receive heartbeat from plugin {%s}: Success\n", pluginName)
	c.SetReadDeadline(time.Now().Add(time.Duration(20) * time.Second))
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Printf("[Error] receive heartbeat from plugin {%s}: %s\n", pluginName, err.Error())
			server.mapMute.Lock()
			server.AliveMap[key] = false
			for _, conn := range server.ConnMap[key] {
				_ = conn.Close()
			}
			server.mapMute.Unlock()
			break
		}
		log.Printf("[Info] receive heartbeat from plugin {%s}: Success\n", pluginName)
		c.SetReadDeadline(time.Now().Add(time.Duration(20) * time.Second))
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
	http.HandleFunc("/receiveHeartbeat", server.receiveHeartbeat)
	http.HandleFunc("/sendHeartbeat", server.sendHeartbeat)
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
				err = Exec("/bin/bash", "-c", fmt.Sprintf("%s %s %s", plugin.Path, *server.Addr, key))
			}
		case PluginTypeJava:
			if runtime.GOOS == "windows" {
				err = Exec("java", "-jar", fmt.Sprintf("%s %s%s", plugin.Path, *server.Addr, key))
			} else {
				err = Exec("/bin/bash", "-c", fmt.Sprintf("java -jar %s %s %s", plugin.Path, *server.Addr, key))
			}
		case PluginTypePython:
			if runtime.GOOS == "windows" {
				err = Exec("python", plugin.Path, fmt.Sprintf("%s %s %s", plugin.Path, *server.Addr, key))
			} else {
				err = Exec("/bin/bash", "-c", fmt.Sprintf("%s %s %s %s", PythonCommand, plugin.Path, *server.Addr, key))
			}
		case PluginTypeJavaScript:
			if runtime.GOOS == "windows" {
				err = Exec("node", plugin.Path, key)
			} else {
				err = Exec("/bin/bash", "-c", fmt.Sprintf("node %s %s %s", plugin.Path, *server.Addr, key))
			}
		}
		if err != nil {
			log.Printf("[Error] Exec plugin {%s}: %s\n", plugin.Name, err.Error())
		} else {
			log.Printf("[Info] Exec plugin {%s}: Success\n", plugin.Name)
		}
	}
}

func (server *PluginServer) listenEvent() {
	for true {
		select {
		case event := <-server.ReceivedEvent:
			event.CompleteType()
			bytes := server.castEventToJSONStr(event)
			fmt.Println(string(bytes))
			for key, ch := range server.MsgQueue {
				server.mapMute.RLock()
				if !server.AliveMap[key] {
					server.mapMute.RUnlock()
					continue
				}
				server.mapMute.RUnlock()
				ch <- bytes
			}
		}
	}
}

func (server *PluginServer) castEventToJSONStr(event Event) []byte {
	ret := make(map[string]interface{})
	ret["type"] = event.Type

	if len(event.Message.DetailList) != 0 {
		switch event.InnerType {
		case EventTypeGroupMessage:
			temp := make(map[string]interface{})
			temp["id"] = CastStringToInt64(event.PreDefVarMap["el-sender-group-id"])
			temp["name"] = event.PreDefVarMap["el-sender-group-name"]
			ret["senderGroup"] = temp

			temp = make(map[string]interface{})
			temp["id"] = CastStringToInt64(event.PreDefVarMap["el-sender-user-id"])
			temp["name"] = event.PreDefVarMap["el-sender-user-name"]
			ret["senderUser"] = temp
		case EventTypeFriendMessage:
			temp := make(map[string]interface{})
			temp["id"] = CastStringToInt64(event.PreDefVarMap["el-sender-user-id"])
			temp["name"] = event.PreDefVarMap["el-sender-user-name"]
			ret["senderUser"] = temp
		default:
			// temp := make(map[string]interface{})
			// temp["id"] = CastStringToInt64(event.PreDefVarMap["el-sender-group-id"])
			// temp["name"] = event.PreDefVarMap["el-sender-group-name"]
			// ret["senderGroup"] = temp

			// temp = make(map[string]interface{})
			// temp["id"] = CastStringToInt64(event.PreDefVarMap["el-sender-user-id"])
			// temp["name"] = event.PreDefVarMap["el-sender-user-name"]
			// ret["senderUser"] = temp
		}
	}

	messageMap := make(map[string]interface{})

	messageMap["at"] = event.Message.At
	messageMap["messageID"] = event.MessageID
	messageMap["detail"] = make([]map[string]interface{}, 0)

	for _, detail := range event.Message.DetailList {
		temp := make(map[string]interface{})
		temp["type"] = detail.Type
		switch detail.InnerType {
		case MessageTypePlain:
			temp["text"] = detail.Text
		case MessageTypeImage:
			temp["url"] = detail.URL
		case MessageTypeFace:
			temp["faceID"] = detail.FaceID
			temp["faceName"] = detail.FaceName
		case MessageTypeAt:
			temp["target"] = detail.UserID
		case MessageTypeAtAll:
		default:
			continue
		}
		messageMap["detail"] = append(messageMap["detail"].([]map[string]interface{}), temp)
	}

	ret["operation"] = make([]map[string]interface{}, 0)

	for _, operation := range event.OperationList {
		temp := make(map[string]interface{})
		temp["type"] = operation.Type
		switch operation.InnerType {
		case OperationTypeGroupMuteAll, OperationTypeGroupUnMuteAll:
			temp["groupID"] = CastStringToInt64(operation.GroupID)
			temp["groupName"] = operation.GroupName
			temp["operatorID"] = CastStringToInt64(operation.OperatorID)
			temp["operatorName"] = operation.OperatorName
		case OperationTypeMemberJoin:
			temp["groupID"] = CastStringToInt64(operation.GroupID)
			temp["groupName"] = operation.GroupName
			temp["userID"] = CastStringToInt64(operation.UserID)
			temp["userName"] = operation.UserName
		case OperationTypeMemberMute, OperationTypeMemberUnMute, OperationTypeMemberLeaveByKick:
			temp["groupID"] = CastStringToInt64(operation.GroupID)
			temp["groupName"] = operation.GroupName
			temp["operatorID"] = CastStringToInt64(operation.OperatorID)
			temp["operatorName"] = operation.OperatorName
			temp["userID"] = CastStringToInt64(operation.UserID)
			temp["userName"] = operation.UserName
		case OperationTypeMemberLeaveByQuit:
			temp["groupID"] = CastStringToInt64(operation.GroupID)
			temp["groupName"] = operation.GroupName
			temp["userID"] = CastStringToInt64(operation.UserID)
			temp["userName"] = operation.UserName
		default:
			continue
		}
		ret["operation"] = append(ret["operation"].([]map[string]interface{}), temp)
	}
	ret["message"] = messageMap
	bytes, _ := json.Marshal(ret)
	return bytes
}
