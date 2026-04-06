# 1. Start the server (Uses your existing shell script)
run:
	@echo "🚀 Starting Redis Server..."
	./your_program.sh

# 2. Basic Tests
ping:
	redis-cli PING

echo:
	redis-cli ECHO hey

# 3. String & Expiry Tests
set:
	redis-cli SET foo bar

get:
	redis-cli GET foo

set-expiry-millis:
	@echo "🕒 Setting key 'foo' with 3000ms expiry..."
	redis-cli SET foo bar PX 3000

set-expiry-secs:
	@echo "🕒 Setting key 'foo' with 10s expiry (using your MX flag)..."
	redis-cli SET foo bar MX 10

# 4. List Tests (RPUSH)
rpush:
	@echo "📥 Testing RPUSH..."
	redis-cli RPUSH my_list "item1"
	redis-cli RPUSH my_list "item2"
rpush-multi:
	@echo "📥 Testing RPUSH..."
	redis-cli RPUSH my_list "item1"
	redis-cli RPUSH my_list "item2" "item3" "item4"

# 5. Combined Stress Test (The "Full Circuit")
test-all: set set-expiry-millis rpush-test
	@echo "✅ All manual tests triggered."