-- noinspection SqlInsertValuesForFile
-- name: GetBySku :one
select count from stocks where sku=$1;

-- name: ChangeReserveOfSkuByAmount :exec
update stocks set count = count+$1 where sku=$2;
