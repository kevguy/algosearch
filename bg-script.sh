trap "exit" INT TERM ERR
trap "kill 0" EXIT

./sync-blocks.sh &
# ./someProcessB &

wait
