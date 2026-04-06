
# 1. Start the server (Uses your existing shell script)
run:
	@echo "🚀 Starting Redis Server..."
	./your_program.sh
# 2. PING and setup TCP handshaking
# 2. Test a single PING (Handshake test)
ping:
	@echo "Testing Single PING..."
	@echo -e "PING\r\n" | nc localhost 6379

multi-ping:
	@echo "PING\nPING"