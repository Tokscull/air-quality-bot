package repository

const (
	saveOrUpdateQuery = `INSERT INTO "bot_user" (id, user_name, chat_id, lang_code, is_active, last_seen_at) 
				VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) 
				DO UPDATE SET user_name = $2, chat_id = $3, lang_code = $4, last_seen_at = $6 RETURNING is_active`

	findAllQuery = `SELECT id, user_name, chat_id, lang_code, is_active, created_at, last_seen_at FROM "bot_user"`

	findByIdQuery = `SELECT id, user_name, chat_id, lang_code, is_active, created_at, last_seen_at 
				FROM "bot_user" WHERE id = $1`
)
