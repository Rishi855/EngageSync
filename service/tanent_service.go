package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

func AddOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	userDetails, err := ExtractUserFromToken(r)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if userDetails.Role != "Superadmin" {
		http.Error(w, "You don't have permission for use this api", http.StatusBadGateway)
		return
	}
	type RequestPayload struct {
		OrgName      string  `json:"orgName"`
		SchemaName   string  `json:"schemaName"`
		UserName     string  `json:"userName"`
		UserEmail    string  `json:"userEmail"`
		UserPassword string  `json:"userPassword"`
		Role         string  `json:"role"` // Admin or User
		PhotoURL     *string `json:"photoURL,omitempty"`
		BirthDate    *string `json:"birthDate,omitempty"`
		Department   *string `json:"department,omitempty"`
	}

	var creds RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tenantID := uuid.New().String()
	globalUserID := uuid.New().String()
	userID := uuid.New().String()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Insert into TenantRegistry
	_, err = tx.Exec(`INSERT INTO TenantRegistry (TenantID, OrgName, SchemaName) VALUES ($1, $2, $3)`,
		tenantID, creds.OrgName, creds.SchemaName)

	if err != nil {
		http.Error(w, "Failed to insert into TenantRegistry: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert into GlobalUsers
	_, err = tx.Exec(`INSERT INTO GlobalUsers (GlobalUserID, Email, Password, TenantID, Role) VALUES ($1, $2, $3, $4, $5)`,
		globalUserID, creds.UserEmail, creds.UserPassword, tenantID, creds.Role)

	if err != nil {
		http.Error(w, "Failed to insert into GlobalUsers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Generate migration file for the schema
	migrationFilename := fmt.Sprintf("migrations/%s_%d_adding_%s_schema.up.sql", time.Now().Format("20060102150405"), time.Now().Unix(), creds.SchemaName)
	migrationContent := strings.ReplaceAll(schemaTemplate, "kanaka", creds.SchemaName)

	if err = os.WriteFile(migrationFilename, []byte(migrationContent), 0644); err != nil {
		http.Error(w, "Failed to create migration file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Run migration command
	migrateCmd := exec.Command("migrate",
		"-database", "postgres://postgres:root@localhost:5432/engagesyncdb?sslmode=disable",
		"-path", "migrations",
		"up",
	)
	migrateCmd.Stdout = os.Stdout
	migrateCmd.Stderr = os.Stderr
	if err = migrateCmd.Run(); err != nil {
		http.Error(w, "Migration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert into new schema.Users
	userInsertQuery := fmt.Sprintf(`INSERT INTO %s.Users (UserID, TenantID, Name, Email, PasswordHash, PhotoURL, BirthDate, Department, Role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, creds.SchemaName)
	_, err = tx.Exec(userInsertQuery, userID, tenantID, creds.UserName, creds.UserEmail,
		creds.UserPassword, creds.PhotoURL, creds.BirthDate, creds.Department, creds.Role)
	if err != nil {
		http.Error(w, "Failed to insert into schema.Users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Organization and user created successfully"))
}

func GetOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	userDetails, err := ExtractUserFromToken(r)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if userDetails.Role != "Superadmin" {
		http.Error(w, "You don't have permission for use this api", http.StatusBadGateway)
		return
	}
	// Query all organizations
	rows, err := db.Query(`SELECT TenantID, OrgName, SchemaName, CreatedAt, IsActive FROM TenantRegistry`)
	if err != nil {
		http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type OrgInfo struct {
		TenantID   string    `json:"tenantId"`
		OrgName    string    `json:"orgName"`
		SchemaName string    `json:"schemaName"`
		CreatedAt  time.Time `json:"createdAt"`
		IsActive   bool      `json:"isActive"`
	}

	var organizations []OrgInfo

	for rows.Next() {
		var org OrgInfo
		if err := rows.Scan(&org.TenantID, &org.OrgName, &org.SchemaName, &org.CreatedAt, &org.IsActive); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		organizations = append(organizations, org)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(organizations)
}

const schemaTemplate = `-- Create schema
CREATE SCHEMA IF NOT EXISTS kanaka;


-- Create tables
CREATE TABLE kanaka.Users (
    UserID UUID PRIMARY KEY,
    TenantID UUID NOT NULL,
    Name TEXT NOT NULL,
    Email TEXT NOT NULL UNIQUE,
    PasswordHash TEXT NOT NULL,
    PhotoURL TEXT,
    BirthDate DATE,
    Department TEXT,
    Role user_role NOT NULL
);

CREATE TABLE kanaka.Projects (
    ProjectID UUID PRIMARY KEY,
    Name TEXT NOT NULL,
    Description TEXT,
    ManagerID UUID NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (ManagerID) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.ProjectMembers (
    ProjectID UUID NOT NULL,
    UserID UUID NOT NULL,
    Role TEXT NOT NULL,
    PRIMARY KEY (ProjectID, UserID),
    FOREIGN KEY (ProjectID) REFERENCES kanaka.Projects(ProjectID),
    FOREIGN KEY (UserID) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.Ideas (
    IdeaID UUID PRIMARY KEY,
    Title TEXT NOT NULL,
    Description TEXT NOT NULL,
    SubmittedBy UUID NOT NULL,
    SubmittedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    IsAnonymous BOOLEAN NOT NULL DEFAULT FALSE,
    IsApproved BOOLEAN DEFAULT NULL,
    ApprovedBy UUID,
    ApprovedAt TIMESTAMP,
    FOREIGN KEY (SubmittedBy) REFERENCES kanaka.Users(UserID),
    FOREIGN KEY (ApprovedBy) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.IdeaComments (
    CommentID UUID PRIMARY KEY,
    IdeaID UUID NOT NULL,
    CommentedBy UUID NOT NULL,
    CommentText TEXT NOT NULL,
    CommentedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (IdeaID) REFERENCES kanaka.Ideas(IdeaID),
    FOREIGN KEY (CommentedBy) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.IdeaVotes (
    VoteID UUID PRIMARY KEY,
    IdeaID UUID NOT NULL,
    VotedBy UUID NOT NULL,
    VotedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (IdeaID) REFERENCES kanaka.Ideas(IdeaID),
    FOREIGN KEY (VotedBy) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.Photos (
    PhotoID UUID PRIMARY KEY,
    ThemeID UUID NOT NULL,
    UploadedBy UUID NOT NULL,
    PhotoURL TEXT NOT NULL,
    UploadedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    IsWinner BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (UploadedBy) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.PhotoVotes (
    VoteID UUID PRIMARY KEY,
    PhotoID UUID NOT NULL,
    VotedBy UUID NOT NULL,
    VotedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (PhotoID) REFERENCES kanaka.Photos(PhotoID),
    FOREIGN KEY (VotedBy) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.WhatIfIdeas (
    IdeaID UUID PRIMARY KEY,
    Description TEXT NOT NULL,
    SubmittedBy UUID NOT NULL,
    SubmittedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    Phase TEXT NOT NULL,
    IsSelected BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (SubmittedBy) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.WhatIfVotes (
    VoteID UUID PRIMARY KEY,
    IdeaID UUID NOT NULL,
    VotedBy UUID NOT NULL,
    VotedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (IdeaID) REFERENCES kanaka.WhatIfIdeas(IdeaID),
    FOREIGN KEY (VotedBy) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.Quizzes (
    QuizID UUID PRIMARY KEY,
    HostID UUID NOT NULL,
    Passcode TEXT NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (HostID) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.QuizQuestions (
    QuestionID UUID PRIMARY KEY,
    QuizID UUID NOT NULL,
    QuestionText TEXT NOT NULL,
    CorrectOption INT NOT NULL,
    SequenceNumber INT NOT NULL,
    FOREIGN KEY (QuizID) REFERENCES kanaka.Quizzes(QuizID)
);

CREATE TABLE kanaka.QuizAnswers (
    AnswerID UUID PRIMARY KEY,
    QuestionID UUID NOT NULL,
    UserID UUID NOT NULL,
    SelectedOption INT NOT NULL,
    AnsweredAt TIMESTAMP NOT NULL DEFAULT NOW(),
    IsCorrect BOOLEAN NOT NULL,
    TimeTaken INT,
    FOREIGN KEY (QuestionID) REFERENCES kanaka.QuizQuestions(QuestionID),
    FOREIGN KEY (UserID) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.QuizParticipants (
    QuizID UUID NOT NULL,
    UserID UUID NOT NULL,
    JoinTime TIMESTAMP NOT NULL DEFAULT NOW(),
    TotalScore INT NOT NULL DEFAULT 0,
    PRIMARY KEY (QuizID, UserID),
    FOREIGN KEY (QuizID) REFERENCES kanaka.Quizzes(QuizID),
    FOREIGN KEY (UserID) REFERENCES kanaka.Users(UserID)
);

CREATE TABLE kanaka.BirthdayWishes (
    WishID UUID PRIMARY KEY,
    UserID UUID NOT NULL,
    WishedBy UUID NOT NULL,
    Message TEXT,
    Emoji TEXT,
    SentAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (UserID) REFERENCES kanaka.Users(UserID),
    FOREIGN KEY (WishedBy) REFERENCES kanaka.Users(UserID)
);
`
