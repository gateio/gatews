# !/usr/bin/env python
# coding: utf-8
import abc
import asyncio
import hashlib
import hmac
import json
import logging
import ssl
import time
import typing

import websockets
from websockets.exceptions import WebSocketException

logger = logging.getLogger(__name__)


class GateWebsocketError(Exception):
    def __init__(self, code, message):
        self.code = code
        self.message = message

    def __str__(self):
        return "code: %d, message: %s" % (self.code, self.message)


class Configuration(object):
    def __init__(
        self,
        app: str = "spot",
        settle: str = "usdt",
        test_net: bool = False,
        host: str = "",
        api_key: str = "",
        api_secret: str = "",
        event_loop=None,
        executor_pool=None,
        default_callback=None,
        ping_interval: int = 5,
        max_retry: int = 10,
        verify: bool = True,
    ):
        """Initialize running configuration

        @param app: Which websocket to connect to, spot or futures, default to spot
        @param settle: If app is futures, which settle currency to use, btc or usdt
        @param test_net: If app is futures, whether use test net
        @param host: Websocket host, inferred from app, settle and test_net if not provided
        @param api_key: APIv4 Key, must not be empty if subscribing to private channels
        @param api_secret: APIv4 Secret, must not be empty if subscribing to private channels
        @param event_loop: Event loop to use. default to asyncio default event loop
        @param executor_pool: Your callback executor pool. Default to asyncio default event loop if callback is
        awaitable, otherwise asyncio default concurrent.futures.Executor executor
        @param default_callback: Default callback function for all channels. If channels specific callback is not
        provided, it will be called instead
        @param ping_interval: Active ping interval to websocket server, default to 5 seconds
        @param max_retry: Connection retry times on connection to server lost. Reconnect will be given up if
        max_retry reached. No upper limit if negative. Default to 10.
        @param verify: enable certificate verification, default to True
        """
        self.app = app
        self.api_key = api_key
        self.api_secret = api_secret
        default_host = "wss://api.gateio.ws/ws/v4/"
        if app == "futures":
            default_host = "wss://fx-ws.gateio.ws/v4/ws/%s" % settle
            if test_net:
                default_host = "wss://fx-ws-testnet.gateio.ws/v4/ws/%s" % settle
        if app == "pilot":
            default_host = "wss://api.gateio.ws/ws/v4/pilot"
        self.host = host or default_host
        self.loop = event_loop
        self.pool = executor_pool
        self.default_callback = default_callback
        self.ping_interval = ping_interval
        self.max_retry = max_retry
        self.verify = verify


class WebSocketRequest(object):
    def __init__(
        self,
        cfg: Configuration,
        channel: str,
        event: str,
        payload: str,
        require_auth: bool,
    ):
        self.channel = channel
        self.event = event
        self.payload = payload
        self.require_auth = require_auth
        self.cfg = cfg

    def __str__(self):
        request = {
            "time": int(time.time()),
            "channel": self.channel,
            "event": self.event,
            "payload": self.payload,
        }
        if self.require_auth:
            if not (self.cfg.api_key and self.cfg.api_secret):
                raise ValueError("configuration does not provide api key or secret")
            message = "channel=%s&event=%s&time=%d" % (
                self.channel,
                self.event,
                request["time"],
            )
            request["auth"] = {
                "method": "api_key",
                "KEY": self.cfg.api_key,
                "SIGN": hmac.new(
                    self.cfg.api_secret.encode("utf8"),
                    message.encode("utf8"),
                    hashlib.sha512,
                ).hexdigest(),
            }
        return json.dumps(request)


class ApiRequest(object):
    def __init__(
        self,
        cfg: Configuration,
        channel: str,
        header: str = "",
        req_id: str = "",
        payload: object = {},
    ):
        self.cfg = cfg
        if not (self.cfg.api_key and self.cfg.api_secret):
            raise ValueError("configuration does not provide api key or secret")
        self.channel = channel
        self.header = header
        self.req_id = req_id
        self.payload = payload

    def gen(self):
        data_time = int(time.time())
        param_json = json.dumps(self.payload)
        message = "%s\n%s\n%s\n%d" % ("api", self.channel, param_json, data_time)

        data_param = {
            "time": data_time,
            "channel": self.channel,
            "event": "api",
            "payload": {
                "req_header": {"X-Gate-Channel-Id": self.header},
                "api_key": self.cfg.api_key,
                "timestamp": f"{data_time}",
                "signature": hmac.new(
                    self.cfg.api_secret.encode("utf8"),
                    message.encode("utf8"),
                    hashlib.sha512,
                ).hexdigest(),
                "req_id": self.req_id,
                "req_param": self.payload,
            },
        }

        return json.dumps(data_param)


class WebSocketResponse(object):
    def __init__(self, body: str):
        self.body = body
        msg = json.loads(body)
        self.channel = msg.get("channel") or (msg.get("header") or {}).get("channel")
        if not self.channel:
            raise ValueError("no channel found from response message: %s" % body)

        self.timestamp = msg.get("time")
        self.event = msg.get("event")
        self.result = (
            msg.get("result")
            or (msg.get("data") or {}).get("result")
            or (msg.get("data") or {}).get("errs")
        )
        self.error = None
        if msg.get("error"):
            self.error = GateWebsocketError(
                msg["error"].get("code"), msg["error"].get("message")
            )

    def __str__(self) -> str:
        return self.body


class Connection(object):
    def __init__(self, cfg: Configuration):
        self.cfg = cfg
        self.channels: typing.Dict[str, typing.Any] = dict()
        self.sending_queue = asyncio.Queue()
        self.sending_history = list()
        self.event_loop: asyncio.AbstractEventLoop = (
            cfg.loop or asyncio.get_event_loop()
        )
        self.main_loop = None

    def register(self, channel, callback=None):
        if callback:
            self.channels[channel] = callback

    def unregister(self, channel):
        self.channels.pop(channel, None)

    def send(self, msg):
        self.sending_queue.put_nowait(msg)

    async def _active_ping(self, conn: websockets.WebSocketClientProtocol):
        while True:
            data = json.dumps(
                {"time": int(time.time()), "channel": "%s.ping" % self.cfg.app}
            )
            await conn.send(data)
            await asyncio.sleep(self.cfg.ping_interval)

    async def _write(self, conn: websockets.WebSocketClientProtocol):
        if self.sending_history:
            for msg in self.sending_history:
                if isinstance(msg, WebSocketRequest):
                    msg = str(msg)
                await conn.send(msg)
        while True:
            msg = await self.sending_queue.get()
            self.sending_history.append(msg)
            if isinstance(msg, WebSocketRequest):
                msg = str(msg)
            await conn.send(msg)

    async def _read(self, conn: websockets.WebSocketClientProtocol):
        async for msg in conn:
            response = WebSocketResponse(msg)
            callback = self.channels.get(response.channel, self.cfg.default_callback)
            if callback is not None:
                if asyncio.iscoroutinefunction(callback):
                    self.event_loop.create_task(callback(self, response))
                else:
                    self.event_loop.run_in_executor(
                        self.cfg.pool, callback, self, response
                    )

    def close(self):
        if self.main_loop:
            self.main_loop.cancel()

    async def run(self):
        stopped = False
        retried = 0
        while not stopped:
            try:
                ctx = None
                if self.cfg.host.startswith("wss://"):
                    ctx = ssl.create_default_context()
                    if not self.cfg.verify:
                        ctx.check_hostname = False
                        ctx.verify_mode = ssl.CERT_NONE
                # compression is not fully supported in server
                conn = await websockets.connect(
                    self.cfg.host, ssl=ctx, compression=None
                )
                if retried > 0:
                    logger.warning(
                        "reconnect succeeded after retrying %d times", retried + 1
                    )
                    retried = 0
            # DNS might be resolved to multiple address, which cause multiple ConnectionRefusedError
            # being combined to one OSError
            except (WebSocketException, ConnectionRefusedError, OSError) as e:
                logger.warning(
                    "failed to connect to server for the %d time, try again later: %s",
                    retried + 1,
                    e,
                )
                retried += 1
                if 0 < self.cfg.max_retry < retried:
                    logger.error(
                        "max reconnect time %d reached, give it up", self.cfg.max_retry
                    )
                    stopped = True
                    continue
                await asyncio.sleep(0.5 * retried)
            else:
                tasks: typing.List[asyncio.Task] = list()
                try:
                    tasks.append(self.event_loop.create_task(self._write(conn)))
                    tasks.append(self.event_loop.create_task(self._read(conn)))
                    tasks.append(self.event_loop.create_task(self._active_ping(conn)))
                    self.main_loop = asyncio.gather(*tasks)
                    await self.main_loop
                except websockets.ConnectionClosed:
                    logger.warning("websocket connection lost, retry to reconnect")
                except asyncio.CancelledError:
                    await conn.close()
                    stopped = True
                finally:
                    # user callback tasks are not our concern
                    for task in tasks:
                        task.cancel()


class BaseChannel(abc.ABC):
    name = ""
    require_auth = False

    def __init__(self, conn: Connection, callback=None):
        self.conn = conn
        self.callback = callback
        self.cfg = self.conn.cfg
        self.conn.register(self.name, callback)

    def subscribe(self, payload={}):
        self.conn.send(
            WebSocketRequest(
                self.cfg, self.name, "subscribe", payload, self.require_auth
            )
        )

    def unsubscribe(self, payload={}):
        self.conn.send(
            WebSocketRequest(
                self.cfg, self.name, "unsubscribe", payload, self.require_auth
            )
        )

    def api_request(self, payload={}, header="", req_id=""):
        self.login(header, req_id)
        self.conn.send(ApiRequest(self.cfg, self.name, header, req_id, payload).gen())

    def login(self, header, req_id):
        channel = "spot.login"
        if self.cfg.app != "spot":
            channel = "futures.login"
        self.conn.send(ApiRequest(self.cfg, channel, header, req_id, {}).gen())
