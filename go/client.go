package gatews

import (
	"context"
	mapset "github.com/deckarep/golang-set"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type WsService struct {
	Logger *log.Logger
	Ctx    context.Context
	Client *websocket.Conn
	once   *sync.Once
	msgChs *sync.Map // business chan
	calls  *sync.Map
	conf   *ConnConf
}

// ConnConf default URL is spot websocket
type ConnConf struct {
	subscribeMsg *sync.Map
	URL          string
	Key          string
	Secret       string
	MaxRetryConn int
}

type ConfOptions struct {
	URL          string
	Key          string
	Secret       string
	MaxRetryConn int
}

func NewWsService(ctx context.Context, logger *log.Logger, conf *ConnConf) (*WsService, error) {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var cfg *ConnConf
	if conf != nil {
		cfg = conf
	} else {
		cfg = getInitConnConf()
	}

	stop := false
	retry := 0
	var conn *websocket.Conn
	for !stop {
		c, _, err := websocket.DefaultDialer.Dial(cfg.URL, nil)
		if err != nil {
			if retry >= cfg.MaxRetryConn {
				log.Printf("max reconnect time %d reached, give it up", cfg.MaxRetryConn)
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
		conf:   cfg,
		Logger: logger,
		Ctx:    ctx,
		Client: conn,
		calls:  new(sync.Map),
		msgChs: new(sync.Map),
		once:   new(sync.Once),
	}

	return ws, nil
}

func getInitConnConf() *ConnConf {
	return &ConnConf{
		subscribeMsg: new(sync.Map),
		MaxRetryConn: MaxRetryConn,
		Key:          "",
		Secret:       "",
		URL:          BaseUrl,
	}
}

func NewConnConf(url, key, secret string, maxRetry int) *ConnConf {
	if url == "" {
		url = BaseUrl
	}
	if maxRetry == 0 {
		maxRetry = MaxRetryConn
	}
	return &ConnConf{
		subscribeMsg: new(sync.Map),
		MaxRetryConn: maxRetry,
		Key:          key,
		Secret:       secret,
		URL:          url,
	}
}

// NewConnConfFromOption conf from options, recommend to use this
func NewConnConfFromOption(op *ConfOptions) *ConnConf {
	if op.URL == "" {
		op.URL = BaseUrl
	}
	if op.MaxRetryConn == 0 {
		op.MaxRetryConn = MaxRetryConn
	}
	return &ConnConf{
		subscribeMsg: new(sync.Map),
		MaxRetryConn: op.MaxRetryConn,
		Key:          op.Key,
		Secret:       op.Secret,
		URL:          op.URL,
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
			ws.Client = c
		}
	}

	// resubscribe after reconnect
	ws.conf.subscribeMsg.Range(func(key, value interface{}) bool {
		// key is channel, value is []requestHistory
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
				ws.Logger.Printf("reconnect channel[%s] with payload[%v] success", key.(string), req.Payload)
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

func (ws *WsService) GetConnection() *websocket.Conn {
	return ws.Client
}
