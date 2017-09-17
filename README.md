# brokerage_server
HTTP-based interface to stock brokers

This is currently an early work in progress. My plan is to eventually add support for Websockets, and probably [NATS](http://nats.io/) as well for more efficient communication with multiple clients.

Currently this is only slated to operate with Interactive Brokers. But the server is built around a plugin-based architecture, so support for other brokerages can be added.

The intention here is to provide a component that has the minimum necessary intelligence to enable higher-level applications to do their jobs. It will not include algorithms for scanning the market, picking stocks, or executing trading algorithms, as it is designed to be a support component for these sorts of applications. It may gain simpler features though, such as price adjustment for unfilled limit orders.
