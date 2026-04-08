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
	redis-cli RPUSH my_list "item1"
	redis-cli RPUSH my_list "item2"
rpush-multi:
	redis-cli RPUSH my_list "item1"
	redis-cli RPUSH my_list "item2" "item3" "item4"
	redis-cli RPUSH list_key "a" "b" "c" "d" "e"
lrange-pos:
	redis-cli RPUSH list_key "a" "b" "c" "d" "e"
	redis-cli LRANGE list_key 0 2
	redis-cli LRANGE list_key 2 4
lrange-neg:
	redis-cli RPUSH banana pear grape pineapple apple strawberry orange raspberry
	redis-cli LRANGE banana -8 -1
lpush:
	redis-cli LPUSH list_key "a" "b" "c"
	redis-cli LRANGE list_key 0 -1
llen:
	redis-cli RPUSH list_key "a" "b" "c" "d"
	redis-cli LLEN list_key
lpop:
	redis-cli RPUSH list_key "one" "two" "three" "four"
	redis-cli LPOP list_key
n_lpop:
	redis-cli RPUSH list_key "a" "b" "c" "d"
	redis-cli LPOP list_key 2
	redis-cli LRANGE list_key 0 -1
blpop:
	redis-cli RPUSH list_key "a" "b" "c" "d"
	redis-cli BLPOP list_key 2
type:
	redis-cli RPUSH list_key "a" "b" "c" "d"	
	redis-cli TYPE list_key
xadd:
	redis-cli XADD stream_key 1526919030474-0 temperature 36 humidity 95
	redis-cli XADD stream_key 0-1 foo bar
xadd-entry-ids:
	redis-cli XADD some_key 1-1 foo bar
	redis-cli XADD some_key 1-1 bar baz
	redis-cli XADD some_key 0-2 bar baz

test-all: set set-expiry-millis rpush rpush-multi lrange-pos lrange-neg lpush llen lpop n_lpop blpop type
	@echo "✅ All manual tests triggered."