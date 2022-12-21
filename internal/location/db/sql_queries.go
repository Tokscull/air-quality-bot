package location

const (
	saveQuery = `INSERT INTO "location" (user_id, latitude, longitude, time_zone) 
				VALUES ($1, $2, $3, $4) RETURNING id`

	findByIdQuery = `SELECT id, user_id, latitude, longitude, time_zone
				FROM "location" WHERE id = $1`

	updateByIdQuery = `UPDATE "location" SET user_id = $1, latitude = $2, longitude = $3, time_zone = $4 WHERE id = $5`

	updateTimeZoneByIdQuery = `UPDATE "location" SET time_zone = $1 WHERE id = $2`

	updateTimeZoneAndNotifyAtByIdQuery = `WITH notification_upd AS 
				(UPDATE "notification" SET notify_at = $1 WHERE location_id = $3)
				UPDATE "location" SET time_zone = $2 WHERE id = $3`
)
