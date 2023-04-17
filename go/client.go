package gatews

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/gorilla/websocket"
)

type WsService struct {
	mu     *sync.Mutex
	Logger *log.Logger
	Ctx    context.Context
	Client *websocket.Conn
	once   *sync.Once
	msgChs *sync.Map // business chan
	calls  *sync.Map
	conf   *ConnConf

	isConnected atomic.Bool
}

// ConnConf default URL is spot websocket
type ConnConf struct {
	subscribeMsg     *sync.Map
	URL              string
	Key              string
	Secret           string
	MaxRetryConn     int
	SkipTlsVerify    bool
	ShowReconnectMsg bool
	PingInterval     string
}

type ConfOptions struct {
	URL              string
	Key              string
	Secret           string
	MaxRetryConn     int
	SkipTlsVerify    bool
	ShowReconnectMsg bool
	PingInterval     string
}

func NewWsService(ctx context.Context, logger *log.Logger, conf *ConnConf) (*WsService, error) {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	defaultConf := getInitConnConf()
	if conf != nil {
		conf = applyOptionConf(defaultConf, conf)
	} else {
		conf = defaultConf
	}

	stop := false
	retry := 0
	var conn *websocket.Conn
	for !stop {
		dialer := websocket.DefaultDialer
		if conf.SkipTlsVerify {
			dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		c, _, err := dialer.Dial(conf.URL, nil)
		if err != nil {
			if retry >= conf.MaxRetryConn {
				log.Printf("max reconnect time %d reached, give it up", conf.MaxRetryConn)
				return nil, err
			}
			retry++
			log.Printf("failed to connect to server for the %d time, try again later", retry)
			time.Sleep(time.Millisecond * (time.Duration(retry) * 500))
			continue
		} else {
			stop = true
			conn = c
		}
	}

	if retry > 0 {
		log.Printf("reconnect succeeded after retrying %d times", retry)
	}

	ws := &WsService{
		mu:     new(sync.Mutex),
		conf:   conf,
		Logger: logger,
		Ctx:    ctx,
		Client: conn,
		calls:  new(sync.Map),
		msgChs: new(sync.Map),
		once:   new(sync.Once),
	}

	ws.isConnected.Store(true)

	go ws.activePing()

	return ws, nil
}

func getInitConnConf() *ConnConf {
	return &ConnConf{
		subscribeMsg:     new(sync.Map),
		MaxRetryConn:     MaxRetryConn,
		Key:              "",
		Secret:           "",
		URL:              BaseUrl,
		SkipTlsVerify:    false,
		ShowReconnectMsg: true,
		PingInterval:     DefaultPingInterval,
	}
}

func applyOptionConf(defaultConf, userConf *ConnConf) *ConnConf {
	if userConf.URL == "" {
		userConf.URL = defaultConf.URL
	}

	if userConf.MaxRetryConn == 0 {
		userConf.MaxRetryConn = defaultConf.MaxRetryConn
	}

	if userConf.PingInterval == "" {
		userConf.PingInterval = defaultConf.PingInterval
	}

	return userConf
}

// NewConnConfFromOption conf from options, recommend using this
func NewConnConfFromOption(op *ConfOptions) *ConnConf {
	if op.URL == "" {
		op.URL = BaseUrl
	}
	if op.MaxRetryConn == 0 {
		op.MaxRetryConn = MaxRetryConn
	}
	return &ConnConf{
		subscribeMsg:     new(sync.Map),
		MaxRetryConn:     op.MaxRetryConn,
		Key:              op.Key,
		Secret:           op.Secret,
		URL:              op.URL,
		SkipTlsVerify:    op.SkipTlsVerify,
		ShowReconnectMsg: op.ShowReconnectMsg,
		PingInterval:     op.PingInterval,
	}
}

func (ws *WsService) GetConnConf() *ConnConf {
	return ws.conf
}

func (ws *WsService) reconnect() error {
	stop := false
	retry := 0
	for !stop {
		c, _, err := websocket.DefaultDialer.Dial(ws.conf.URL, nil)
		if err != nil {
			if retry >= ws.conf.MaxRetryConn {
				ws.Logger.Printf("max reconnect time %d reached, give it up", ws.conf.MaxRetryConn)
				return err
			}
			retry++
			log.Printf("failed to connect to server for the %d time, try again later", retry)
			time.Sleep(time.Millisecond * (time.Duration(retry) * 500))
			continue
		} else {
			stop = true

			ws.mu.Lock()
			ws.Client = c
			ws.isConnected.Store(true)
			ws.mu.Unlock()
		}
	}

	// resubscribe after reconnect
	ws.conf.subscribeMsg.Range(func(key, value interface{}) bool {
		// key is channel, value is []requestHistory
		if _, ok := value.([]requestHistory); ok {
			for _, req := range value.([]requestHistory) {
				if req.op == nil {
					req.op = &SubscribeOptions{
						IsReConnect: true,
					}
				} else {
					req.op.IsReConnect = true
				}
				if err := ws.baseSubscribe(req.Event, req.Channel, req.Payload, req.op); err != nil {
					ws.Logger.Printf("after reconnect, subscribe channel[%s] err:%s", key.(string), err.Error())
				} else {
					if ws.conf.ShowReconnectMsg {
						ws.Logger.Printf("reconnect channel[%s] with payload[%v] success", key.(string), req.Payload)
					}
				}
			}
		}
		return true
	})

	return nil
}

func (ws *WsService) SetKey(key string) {
	ws.conf.Key = key
}

func (ws *WsService) GetKey() string {
	return ws.conf.Key
}

func (ws *WsService) SetSecret(secret string) {
	ws.conf.Secret = secret
}

func (ws *WsService) GetSecret() string {
	return ws.conf.Secret
}

func (ws *WsService) SetMaxRetryConn(max int) {
	ws.conf.MaxRetryConn = max
}

func (ws *WsService) GetMaxRetryConn() int {
	return ws.conf.MaxRetryConn
}

func (ws *WsService) GetChannelMarkets(channel string) []string {
	var markets []string
	set := mapset.NewSet()
	if v, ok := ws.conf.subscribeMsg.Load(channel); ok {
		for _, req := range v.([]requestHistory) {
			if req.Event == Subscribe {
				for _, pl := range req.Payload {
					if strings.Contains(pl, "_") {
						set.Add(pl)
					}
				}
			} else {
				for _, pl := range req.Payload {
					if strings.Contains(pl, "_") {
						set.Remove(pl)
					}
				}
			}
		}

		for _, v := range set.ToSlice() {
			markets = append(markets, v.(string))
		}
		return markets
	}
	return markets
}

func (ws *WsService) GetChannels() []string {
	var channels []string
	ws.calls.Range(func(key, value interface{}) bool {
		channels = append(channels, key.(string))
		return true
	})
	return channels
}

func (ws *WsService) IsConnected() bool {
	return ws.isConnected.Load()
}

func (ws *WsService) Close() {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ws.Client != nil {
		ws.isConnected.Store(false)
		if err := ws.Client.Close(); err != nil {
			ws.Logger.Printf("close err: %s", err.Error())
		}

		ws.Client = nil
	}
}

func (ws *WsService) activePing() {
	du, err := time.ParseDuration(ws.conf.PingInterval)
	if err != nil {
		ws.Logger.Printf("failed to parse ping interval: %s, use default ping interval 10s instead", ws.conf.PingInterval)
		du, err = time.ParseDuration(DefaultPingInterval)
		if err != nil {
			du = time.Second * 10
		}
	}

	ticker := time.NewTicker(du)
	defer ticker.Stop()

	for {
		select {
		case <-ws.Ctx.Done():
			return
		case <-ticker.C:
			subscribeMap := map[string]int{}
			ws.conf.subscribeMsg.Range(func(key, value interface{}) bool {
				splits := strings.Split(key.(string), ".")
				if len(splits) == 2 {
					subscribeMap[splits[0]] = 1
				}
				return true
			})

			for app := range subscribeMap {
				ws.Subscribe(fmt.Sprintf("%s.ping", app), nil)
			}
		}
	}
}
