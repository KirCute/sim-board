package sim_board

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ClientMessage struct {
	ID      uint64          `json:"-"`
	Command string          `json:"cmd"`
	Player  string          `json:"player"`
	Room    string          `json:"room"`
	Data    json.RawMessage `json:"data"`
}

type ServerMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

var roomMap sync.Map

type joinQuitMsg struct {
	player string
	conn   *websocket.Conn
}

type Room struct {
	quitChan   chan *joinQuitMsg
	joinChan   chan *joinQuitMsg
	cmdChan    chan *ClientMessage
	conn       map[string][]*websocket.Conn
	board      map[string]*PublicCard
	hole       map[string]map[DeckCard]int
	decks      []Deck
	placeCnter uint
}

func GetOrCreateRoom(name string) *Room {
	r, ok := roomMap.LoadOrStore(name, &Room{})
	room := r.(*Room)
	if !ok {
		room.cmdChan = make(chan *ClientMessage, 64)
		room.joinChan = make(chan *joinQuitMsg, 8)
		room.quitChan = make(chan *joinQuitMsg, 8)
		room.conn = make(map[string][]*websocket.Conn)
		room.board = make(map[string]*PublicCard)
		room.hole = make(map[string]map[DeckCard]int)
		go room.handleCommand(name)
	}
	return room
}

func GetRoom(name string) (*Room, bool) {
	r, ok := roomMap.Load(name)
	return r.(*Room), ok
}

func RemoveRoom(name string) bool {
	_, ok := roomMap.LoadAndDelete(name)
	return ok
}

func (r *Room) TryRemoveSubscribe(player string, conn *websocket.Conn) {
	r.quitChan <- &joinQuitMsg{player, conn}
}

func (r *Room) PushSubscribe(player string, conn *websocket.Conn) {
	r.joinChan <- &joinQuitMsg{player, conn}
}

func (r *Room) PushRequest(msg *ClientMessage) {
	r.cmdChan <- msg
}

func (r *Room) handleQuit(m *joinQuitMsg) {
	arr, ok := r.conn[m.player]
	if !ok {
		return
	}
	for i, conn := range arr {
		if m.conn == conn {
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}
}

func (r *Room) handleJoin(m *joinQuitMsg) {
	r.conn[m.player] = append(r.conn[m.player], m.conn)
	if _, ok := r.hole[m.player]; !ok {
		r.hole[m.player] = make(map[DeckCard]int)
	}
	r.handleWelcome(m.player)
}

func handle[T any](r *Room, player string, data json.RawMessage, f func(string, T)) {
	var args T
	if err := json.Unmarshal(data, &args); err != nil {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "请求格式错误"})
		return
	}
	f(player, args)
}

func (r *Room) handleCommand(name string) {
	ticker := time.NewTicker(30 * time.Minute)
	defer func() {
		logrus.Infof("room '%s' expired", name)
		ticker.Stop()
		RemoveRoom(name)
		close(r.cmdChan)
	}()
	logrus.Infof("created new room '%s'", name)
	for {
		select {
		case quit := <-r.quitChan:
			r.handleQuit(quit)
		case join := <-r.joinChan:
			r.handleJoin(join)
		case msg := <-r.cmdChan:
			switch msg.Command {
			case "draw":
				handle(r, msg.Player, msg.Data, r.handleDraw)
			case "announce":
				handle(r, msg.Player, msg.Data, r.handleAnnounce)
			case "collect":
				handle(r, msg.Player, msg.Data, r.handleCollect)
			case "all_collect":
				r.handleAllCollect(msg.Player)
			case "discard_board":
				handle(r, msg.Player, msg.Data, r.handleDiscardBoard)
			case "discard_hole":
				handle(r, msg.Player, msg.Data, r.handleDiscardHole)
			case "add_deck":
				handle(r, msg.Player, msg.Data, r.handleAddDeck)
			case "move":
				handle(r, msg.Player, msg.Data, r.handleMove)
			case "reset":
				r.handleReset()
			}
			ticker.Reset(30 * time.Minute)
		case <-ticker.C:
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func initStatic(r *gin.Engine, public fs.FS) {
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/assets/") {
			c.Header("Cache-Control", "public, max-age=15552000")
		}
	})
	assets, err := fs.Sub(public, "assets")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/assets/", http.FS(assets))
	r.StaticFileFS("/favicon.ico", "./favicon.ico", http.FS(public))
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != "GET" && c.Request.Method != "POST" {
			c.Status(405)
			return
		}
		f, e := public.Open("index.html")
		if e != nil {
			c.Status(500)
			return
		}
		c.Header("Content-Type", "text/html")
		c.Status(200)
		_, _ = io.Copy(c.Writer, f)
		c.Writer.Flush()
		c.Writer.WriteHeaderNow()
	})
}

func Run(public fs.FS) {
	r := gin.Default()
	initStatic(r, public)
	r.GET("/ws", func(c *gin.Context) {
		logrus.Infof("handshake from %s", c.ClientIP())
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logrus.Errorf("WebSocket upgrade failed: %+v", err)
			return
		}
		defer func() {
			if e := conn.Close(); e != nil {
				logrus.Errorf("failed to close WebSocket connection from %s: %+v", c.ClientIP(), e)
			} else {
				logrus.Infof("WebSocket connection from %s closed", c.ClientIP())
			}
		}()
		handleClientMessage(c.ClientIP(), conn)
	})
	logrus.Infof("listening on :6700")
	if err := r.Run(":6700"); err != nil {
		panic(err)
	}
}

type playerRoomPair struct {
	room   string
	player string
}

func handleClientMessage(ip string, conn *websocket.Conn) {
	var mid uint64 = 0
	var joined []*playerRoomPair
	defer func() {
		for _, pair := range joined {
			room, ok := GetRoom(pair.room)
			if !ok {
				continue
			}
			room.TryRemoveSubscribe(pair.player, conn)
		}
	}()
	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err) {
				logrus.Errorf("read msg from %s failed: %+v", ip, err)
				_ = conn.WriteJSON(&ServerMessage{Type: "fatal", Data: err.Error()})
			}
			break
		}
		var msg ClientMessage
		if err = json.Unmarshal(m, &msg); err != nil {
			logrus.Errorf("unmarshal msg from %s failed: %+v", ip, err)
			_ = conn.WriteJSON(&ServerMessage{Type: "error", Data: "无效的请求"})
			continue
		}
		if msg.Command == "nop" {
			continue
		}
		msg.ID = mid
		mid++
		logrus.Infof("recv msg player=%s, room=%s, mid=%d, cmd=%s, data=%s", msg.Player, msg.Room, msg.ID, msg.Command, string(msg.Data))
		if msg.Command == "download" {
			var req downloadRequest
			if err = json.Unmarshal(msg.Data, &req); err != nil {
				logrus.Errorf("failed to unmarshal download req(mid=%d) from %s: %+v", msg.ID, ip, err)
				_ = conn.WriteJSON(&ServerMessage{Type: "error", Data: "请求错误"})
				continue
			}
			handleDownload(req, conn)
			continue
		}
		if msg.Room == "" || msg.Player == "" {
			_ = conn.WriteJSON(&ServerMessage{Type: "error", Data: "房间名和玩家名不能为空"})
			continue
		}
		if msg.Command == "join" {
			room := GetOrCreateRoom(msg.Room)
			joined = append(joined, &playerRoomPair{room: msg.Room, player: msg.Room})
			room.PushSubscribe(msg.Player, conn)
			continue
		}
		room, ok := GetRoom(msg.Room)
		if !ok {
			logrus.Errorf("failed to handle command(player=%s, room=%s, mid=%d, cmd=%s): room not exists", msg.Player, msg.Room, msg.ID, msg.Command)
			_ = conn.WriteJSON(&ServerMessage{Type: "error", Data: "房间不存在"})
			continue
		}
		room.PushRequest(&msg)
	}
}

type downloadRequest struct {
	Deck string `json:"deck"`
	Card string `json:"card"`
}

type downloadResponse struct {
	downloadRequest
	HTML string `json:"html"`
}

func handleDownload(req downloadRequest, conn *websocket.Conn) {
	h, ok := GetCardHTML(req.Deck, req.Card)
	if !ok {
		_ = conn.WriteJSON(&ServerMessage{Type: "error", Data: fmt.Sprintf("不存在的模型：%s@%s", req.Deck, req.Card)})
		return
	}
	_ = conn.WriteJSON(&ServerMessage{Type: "download", Data: downloadResponse{
		downloadRequest: req,
		HTML:            h,
	}})
}
