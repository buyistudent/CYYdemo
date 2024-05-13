
-- name: SaveEsim :execresult
INSERT INTO `t_cp_running_bill_esim` ( `id`, `iccid`, `imsi`, `msisdn`, `json`,  `create_time`, `status`)
VALUES
    ( ?, ?, ?, ?,?,?, ? );


-- name: UpdateEsimEffectiveTime :exec
UPDATE t_cp_running_bill_esim SET effective_time = ?  ,  `status` = ?
WHERE

        id = ?;

