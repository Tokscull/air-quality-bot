package notification

const (
	saveQuery = `INSERT INTO "notification" (user_id, location_id, is_active, notify_at) 
				VALUES ($1, $2, $3, $4) RETURNING id`

	findById = `SELECT n.id, n.user_id, n.location_id, n.notify_at, n.is_active 
				FROM "notification" n WHERE n.id = $1`

	findByIdWithLocationQuery = `SELECT n.id, n.notify_at, n.is_active, n.user_id,
       			l.id, l.user_id, l.latitude, l.longitude, l.time_zone
			  	FROM "notification" n 
			  	INNER JOIN "location" l on l.id = n.location_id 
			  	WHERE n.id = $1`

	findAllByUserIdWithLocationQuery = `SELECT n.id, n.notify_at, n.is_active, n.user_id,
    			l.id, l.user_id, l.latitude, l.longitude, l.time_zone
				FROM "notification" n
              	INNER JOIN "location" l on l.id = n.location_id 
              	WHERE n.user_id = $1 ORDER BY n.notify_at`

	findAllByNotifyAtPerPageQuery = `SELECT n.id, n.notify_at, n.is_active,
       			u.id, u.user_name, u.chat_id, u.lang_code, u.is_active,
       			l.id, l.user_id, l.latitude, l.longitude, l.time_zone
			  	FROM "notification" n 
    		  	INNER JOIN "bot_user" u on u.id = n.user_id 
    		  	INNER JOIN "location" l on l.id = n.location_id 
    		  	WHERE n.notify_at = $1 AND n.is_active = true AND u.is_active = true 
    		  	LIMIT $2 OFFSET $3`

	updateLastTimeProcessedAtQuery = `UPDATE "notification" SET last_time_processed_at = $1 WHERE id = $2`

	updateIsActiveByIdQuery = `UPDATE "notification" SET is_active = not is_active WHERE id = $1 RETURNING is_active`

	updateNotifyAtByIdQuery = `UPDATE "notification" SET notify_at = $1 WHERE id = $2`

	deleteByIdWithLocationQuery = `WITH ntf AS 
				(DELETE FROM "notification" WHERE id = $1 RETURNING location_id)
				DELETE FROM "location" WHERE id = (SELECT location_id from ntf)`
)
