package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
)

// User represents a platform user
type User struct {
    ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    Username     string         `gorm:"uniqueIndex;size:64;not null" json:"username"`
    Email        string         `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string         `gorm:"not null" json:"-"`
    CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
    Status       string         `gorm:"not null;default:active" json:"status"`
    Country      *string        `json:"country"`
    Profile      datatypes.JSON `gorm:"type:jsonb" json:"profile"`
}

type Developer struct {
    ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    Name        string         `gorm:"uniqueIndex;not null" json:"name"`
    Country     *string        `json:"country"`
    Website     *string        `json:"website"`
    ContactEmail *string       `json:"contact_email"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
    Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

type Game struct {
    ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    DeveloperID uuid.UUID      `gorm:"type:uuid;not null" json:"developer_id"`
    Title       string         `gorm:"not null" json:"title"`
    Description string         `json:"description"`
    Genre       string         `json:"genre"`
    Price       float64        `json:"price"`
    ReleaseDate *time.Time     `json:"release_date"`
    AvgRating   float64        `json:"avg_rating"`
    SalesCount  int64          `json:"sales_count"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
    Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

type GameLicense struct {
    ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    GameID         uuid.UUID  `gorm:"type:uuid;not null" json:"game_id"`
    Key            string     `gorm:"uniqueIndex;not null" json:"key"`
    AssignedToUser *uuid.UUID `gorm:"type:uuid" json:"assigned_to_user"`
    AssignedAt     *time.Time `json:"assigned_at"`
    State          string     `json:"state"`
    Notes          *string    `json:"notes"`
}

type Transaction struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
    Amount    float64        `json:"amount"`
    Currency  string         `json:"currency"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    Provider  *string        `json:"provider"`
    Status    string         `json:"status"`
    Details   datatypes.JSON `gorm:"type:jsonb" json:"details"`
}

type Purchase struct {
    ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
    GameID        uuid.UUID  `gorm:"type:uuid;not null" json:"game_id"`
    TransactionID *uuid.UUID `gorm:"type:uuid" json:"transaction_id"`
    PricePaid     float64    `json:"price_paid"`
    PurchasedAt   time.Time  `json:"purchased_at"`
    Method        string     `json:"method"`
    LicenseID     *uuid.UUID `gorm:"type:uuid" json:"license_id"`
}

type Library struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    GameID    uuid.UUID `gorm:"type:uuid;not null" json:"game_id"`
    LicenseID *uuid.UUID `gorm:"type:uuid" json:"license_id"`
    AddedAt   time.Time `json:"added_at"`
    Active    bool      `json:"active"`
}

type Review struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    GameID      uuid.UUID `gorm:"type:uuid;not null" json:"game_id"`
    Rating      int       `json:"rating"`
    Title       *string   `json:"title"`
    Body        *string   `json:"body"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   *time.Time `json:"updated_at"`
    HelpfulCount int      `json:"helpful_count"`
}

type Achievement struct {
    ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    GameID       uuid.UUID `gorm:"type:uuid;not null" json:"game_id"`
    Code         string    `json:"code"`
    Title        string    `json:"title"`
    Description  *string   `json:"description"`
    RewardPoints int       `json:"reward_points"`
    CreatedAt    time.Time `json:"created_at"`
}

type UserAchievement struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    AchievementID uuid.UUID `gorm:"type:uuid;not null" json:"achievement_id"`
    AchievedAt    time.Time `json:"achieved_at"`
    Progress      datatypes.JSON `gorm:"type:jsonb" json:"progress"`
}

type BatchError struct {
    ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    BatchName string         `json:"batch_name"`
    RowData   datatypes.JSON `gorm:"type:jsonb" json:"row_data"`
    Error     string         `json:"error"`
    CreatedAt time.Time      `json:"created_at"`
}

type AuditLog struct {
    ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    TableName   string         `json:"table_name"`
    Operation   string         `json:"operation"`
    RowID       *uuid.UUID     `gorm:"type:uuid" json:"row_id"`
    PerformedBy *uuid.UUID     `gorm:"type:uuid" json:"performed_by"`
    PerformedAt time.Time      `json:"performed_at"`
    OldData     datatypes.JSON `gorm:"type:jsonb" json:"old_data"`
    NewData     datatypes.JSON `gorm:"type:jsonb" json:"new_data"`
    Query       *string        `json:"query"`
}