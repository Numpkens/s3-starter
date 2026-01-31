package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ThumbnailURL *string   `json:"thumbnail_url"`
	VideoURL     *string   `json:"video_url"`
	CreateVideoParams
}

type CreateVideoParams struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
}

func (c Client) GetVideo(id uuid.UUID) (Video, error) {
	query := `
	SELECT id, created_at, updated_at, title, description, thumbnail_url, video_url, user_id
	FROM videos WHERE id = ?`

	var video Video
	err := c.db.QueryRow(query, id).Scan(
		&video.ID, &video.CreatedAt, &video.UpdatedAt, &video.Title,
		&video.Description, &video.ThumbnailURL, &video.VideoURL, &video.UserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Video{}, err
		}
		return Video{}, err
	}
	return video, nil
}

func (c Client) UpdateVideo(video Video) error {
	query := `
	UPDATE videos
	SET
		title = ?,
		description = ?,
		thumbnail_url = ?,
		video_url = ?,
		user_id = ?,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	_, err := c.db.Exec(
		query,
		video.Title,
		video.Description,
		video.ThumbnailURL,
		video.VideoURL,
		video.UserID,
		video.ID,
	)
	return err
}

func (c Client) CreateVideo(params CreateVideoParams) (Video, error) {
	id := uuid.New()
	query := `
    INSERT INTO videos (id, title, description, thumbnail_url, video_url, user_id)
    VALUES (?, ?, ?, ?, ?, ?)`
	_, err := c.db.Exec(query, id.String(), params.Title, params.Description, nil, nil, params.UserID.String())
	if err != nil {
		return Video{}, err
	}
	return c.GetVideo(id)
}

func (c Client) DeleteVideo(id uuid.UUID) error {
	res, err := c.db.Exec("DELETE FROM videos WHERE id = ?", id.String())
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (c Client) GetVideos() ([]Video, error) {
	query := `
    SELECT id, created_at, updated_at, title, description, thumbnail_url, video_url, user_id
    FROM videos
    ORDER BY created_at DESC
    `
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		if err := rows.Scan(
			&v.ID, &v.CreatedAt, &v.UpdatedAt, &v.Title,
			&v.Description, &v.ThumbnailURL, &v.VideoURL, &v.UserID,
		); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}