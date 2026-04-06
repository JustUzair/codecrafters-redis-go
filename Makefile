
# 1. Start the server (Uses your existing shell script)
run:
	@echo "🚀 Starting Redis Server..."
	./your_program.sh
# 2. PING and setup TCP handshaking
# 2. Test a single PING (Handshake test)
ping:
# 	@echo "Testing Single PING..."
# 	@echo -e "PING\r\n" | nc localhost 6379
	@printf '$$4\r\nPING\r\n' | nc localhost 6379

echo:
	@printf "*2\r\n\$$4\r\nECHO\r\n\$$3\r\nhey\r\n" | nc localhost 6379

get:
	@printf 'redis-cli GET foo'
set:
	@printf 'redis-cli SET foo bar'
set-expiry-millis:
	@echo -e "Setting key for 10000ms"
	@printf 'redis-cli SET foo bar PX 10000'
set-expiry-secs:
	@echo -e "Setting key for 10s"
	@printf 'redis-cli SET foo bar MX 10"
