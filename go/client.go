package gatews

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"sync"
	"time"
)

type WsService struct {
	Logger *log.Logger
	Ctx    context.Context
	Client *websocket.Conn
	once   *sync.Once
	buChs  *sync.Map
	calls  *sync.Map
	conf   *ConnConf
}

type ConnConf struct {
	markets      *sync.Map
	Key          string
	Secret       string
	URL          string
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
		buChs:  new(sync.Map),
		once:   new(sync.Once),
	}

	return ws, nil
}

func getInitConnConf() *ConnConf {
	return &ConnConf{
		markets:      new(sync.Map),
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
		markets:      new(sync.Map),
		MaxRetryConn: maxRetry,
		Key:          key,
		Secret:       secret,
		URL:          url,
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
				log.Printf("max reconnect time %d reached, give it up", ws.conf.MaxRetryConn)
				ws.Client.Close()
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
	ws.conf.markets.Range(func(key, value interface{}) bool {
		if err := ws.Subscribe(key.(string), value.([]string)); err != nil {
			log.Printf("after reconnect, subscribe channel[%s] err:%s", key.(string), err.Error())
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
	if v, ok := ws.conf.markets.Load(channel); ok {
		return v.([]string)
	}
	return nil
}

func (ws *WsService) GetChannels() []string {
	var channels []string
	ws.calls.Range(func(key, value interface{}) bool {
		channels = append(channels, key.(string))
		return true
	})
	return channels
}
