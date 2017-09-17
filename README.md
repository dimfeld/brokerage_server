# brokerage_server
HTTP-based interface to the Interactive Brokers API

This is currently an early work in progress. My plan is to eventually add support for Websockets, and probably [NATS](http://nats.io/) as well for more efficient communication with multiple clients.

The intention here is to provide a component that has the minimum necessary intelligence to enable higher-level applications to do their jobs. It will not include algorithms for scanning the market, picking stocks, or executing trading algorithms, as it is designed to be a support component for these sorts of applications. It may gain simpler features though, such as price adjustment for unfilled limit orders.
