package models

import (
	"backend/internal/db"
	"log"
	"time"
)

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created-at"`
	AvatarURL string    `json:"avatar-url"`
	Status    string    `json:"status"`
}

type Room struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Discription string    `json:"discription"`
	IsPrivate   bool      `json:"is-private"`
	CreatedAt   time.Time `json:"created-at"`
	CreatedBy   string    `json:"created-by"`
}

type Message struct {
	ID        uint64     `json:"id"`
	RoomID    uint64     `json:"room-id"`
	UserID    uint64     `json:"user-id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created-at"`
	EditedAt  *time.Time `json:"edited-at"`
}

type Membership struct {
	ID        uint64    `json:"id"`
	RoomID    uint64    `json:"room-id"`
	UserID    uint64    `json:"user-id"`
	Role      string    `json:"role"`
	JoindedAt time.Time `json:"joined-at"`
}

type Call struct {
	ID        uint64     `json:"id"`
	RoomID    uint64     `json:"room-id"`
	StartedAt time.Time  `json:"started-at"`
	EndedAt   *time.Time `json:"ended-at"`
	IsActive  bool       `json:"is-active"`
}

type FileAttachment struct {
	ID         uint64    `json:"id"`
	MessageID  uint64    `json:"message-id"`
	URL        string    `json:"url"`
	FileType   string    `json:"file-type"`
	UploadedAt time.Time `json:"uploaded-at"`
}

func NewUser(username string, email string, password string, avatar ...string) *User {
	query := `
		INSERT INTO users 
		(username, email, password, created_at, avatar_url, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURN id;
	`
	var id uint64
	var avatar_url string

	if len(avatar) > 0 {
		avatar_url = avatar[0]
	} else {
		avatar_url = ""
	}

	err := db.DB.QueryRow(query, username, email, password, time.Now(), avatar_url, "offline").Scan(&id)
	if err != nil {
		log.Fatalf("User creation query failed: %v", err)
	}

	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		AvatarURL: "",
		Status:    "offline",
	}
}

func LoadUser(username string) *User {
	query := `
		SELECT *
		FROM users
		WHERE username = %1;
	`
	loaded_user := &User{}
	err := db.DB.QueryRow(query, username).Scan(&loaded_user.ID, &loaded_user.Username, &loaded_user.Email, &loaded_user.Password, &loaded_user.AvatarURL, &loaded_user.CreatedAt, &loaded_user.Status)
	if err != nil {
		log.Fatalf("Unable to load user: %v", err)
	}

	return loaded_user

}

func LoadAllUsers() []*User {
	query := `
		SELECT *
		FROM users
	`
	var allUsers []*User
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Fatalf("Unable to load all users: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		user := &User{}

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AvatarURL, &user.CreatedAt, &user.Status)
		if err != nil {
			log.Fatalf("Unable to scan user: %v", err)
		}

		allUsers = append(allUsers, user)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error while row iteration: %v", err)
	}

	return allUsers
}

// User Methods
func (u *User) SetAvatar(avatar string) {
	u.AvatarURL = avatar
}

// Room

func NewRoom(name string, disc string, author string, isprivate bool) *Room {
	query := `
		INSERT INTO users 
		(name, discription, is_private, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURN id;
	`
	var roomid uint64

	err := db.DB.QueryRow(query, name, disc, isprivate, time.Now(), author).Scan(&roomid)
	if err != nil {
		log.Fatalf("Room creation query failed: %v", err)
	}

	return &Room{
		ID:          roomid,
		Name:        name,
		Discription: disc,
		IsPrivate:   isprivate,
		CreatedAt:   time.Now(),
		CreatedBy:   author,
	}

}

func LoadRoom(name string) *Room {
	query := `
		SELECT *
		FROM rooms
		WHERE name = %1;
	`
	loaded_room := &Room{}
	err := db.DB.QueryRow(query, name).Scan(&loaded_room.ID, &loaded_room.Name, &loaded_room.Discription, &loaded_room.IsPrivate, &loaded_room.CreatedAt, &loaded_room.CreatedBy)
	if err != nil {
		log.Fatalf("Unable to load user: %v", err)
	}

	return loaded_room

}

func LoadAllRooms() []*Room {
	query := `
		SELECT *
		FROM rooms
	`
	var allRooms []*Room
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Fatalf("Unable to load all rooms: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		room := &Room{}

		err := rows.Scan(&room.ID, &room.Name, &room.Discription, &room.IsPrivate, &room.CreatedAt, &room.CreatedBy)
		if err != nil {
			log.Fatalf("Unable to scan room: %v", err)
		}

		allRooms = append(allRooms, room)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error while row iteration: %v", err)
	}

	return allRooms
}

// Room methods
func (r *Room) SetIsPrivate(res bool) {
	r.IsPrivate = res
}

func (r *Room) LoadAllMembers() []*User {
	query_1 := `
		SELECT user_id
		FROM memberships
		WHERE room_id = %1
	`
	query_2 := `
		SELECT *
		FROM users
		WHERE id = %1
	`
	var members []*User
	var id uint64

	rows, err := db.DB.Query(query_1, r.ID)
	if err != nil {
		log.Fatalf("Unable to query memberships: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		user := &User{}

		err := rows.Scan(&id)
		if err != nil {
			log.Fatalf("Unable to scan row: %v", err)
		}

		er := db.DB.QueryRow(query_2, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.AvatarURL, &user.Status)
		if er != nil {
			log.Fatalf("Unable to fetch member: %v", er)
		}

		members = append(members, user)
	}

	return members

}

func (r *Room) LoadOnlineMembers() []*User {
	// TODO
	var online []*User

	return online
}

// Message

func NewMessage(roomID uint64, userID uint64, content string) *Message {
	query := `
		INSERT INTO messages
		(room_id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURN id
	`
	var id uint64
	err := db.DB.QueryRow(query, roomID, userID, content, time.Now()).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to create new message: %v", err)
	}

	return &Message{
		ID:        id,
		RoomID:    roomID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		EditedAt:  nil,
	}
}

func LoadAllMessages() []*Message {
	query := `
		SELECT *
		FROM messages
	`
	var allMessages []*Message
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Fatalf("Unable to load all messages: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		message := &Message{}

		err := rows.Scan(&message.ID, &message.RoomID, &message.UserID, &message.Content, &message.CreatedAt, &message.EditedAt)
		if err != nil {
			log.Fatalf("Unable to scan room: %v", err)
		}

		allMessages = append(allMessages, message)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error while row iteration: %v", err)
	}

	return allMessages
}

// Membership

func NewMembership(roomID uint64, userID uint64, role string) *Membership {
	query := `
		INSERT INTO memberships
		(room_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)
		RETURN id
	`
	var id uint64

	err := db.DB.QueryRow(query, roomID, userID, role, time.Now()).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to create membership: %v", err)
	}

	return &Membership{
		ID:        id,
		RoomID:    roomID,
		UserID:    userID,
		Role:      role,
		JoindedAt: time.Now(),
	}
}
