-- Create schema
CREATE SCHEMA IF NOT EXISTS superadmin5;


-- Create tables
CREATE TABLE superadmin5.Users (
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

CREATE TABLE superadmin5.Projects (
    ProjectID UUID PRIMARY KEY,
    Name TEXT NOT NULL,
    Description TEXT,
    ManagerID UUID NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (ManagerID) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.ProjectMembers (
    ProjectID UUID NOT NULL,
    UserID UUID NOT NULL,
    Role TEXT NOT NULL,
    PRIMARY KEY (ProjectID, UserID),
    FOREIGN KEY (ProjectID) REFERENCES superadmin5.Projects(ProjectID),
    FOREIGN KEY (UserID) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.Ideas (
    IdeaID UUID PRIMARY KEY,
    Title TEXT NOT NULL,
    Description TEXT NOT NULL,
    SubmittedBy UUID NOT NULL,
    SubmittedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    IsAnonymous BOOLEAN NOT NULL DEFAULT FALSE,
    IsApproved BOOLEAN DEFAULT NULL,
    ApprovedBy UUID,
    ApprovedAt TIMESTAMP,
    FOREIGN KEY (SubmittedBy) REFERENCES superadmin5.Users(UserID),
    FOREIGN KEY (ApprovedBy) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.IdeaComments (
    CommentID UUID PRIMARY KEY,
    IdeaID UUID NOT NULL,
    CommentedBy UUID NOT NULL,
    CommentText TEXT NOT NULL,
    CommentedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (IdeaID) REFERENCES superadmin5.Ideas(IdeaID),
    FOREIGN KEY (CommentedBy) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.IdeaVotes (
    VoteID UUID PRIMARY KEY,
    IdeaID UUID NOT NULL,
    VotedBy UUID NOT NULL,
    VotedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (IdeaID) REFERENCES superadmin5.Ideas(IdeaID),
    FOREIGN KEY (VotedBy) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.Photos (
    PhotoID UUID PRIMARY KEY,
    ThemeID UUID NOT NULL,
    UploadedBy UUID NOT NULL,
    PhotoURL TEXT NOT NULL,
    UploadedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    IsWinner BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (UploadedBy) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.PhotoVotes (
    VoteID UUID PRIMARY KEY,
    PhotoID UUID NOT NULL,
    VotedBy UUID NOT NULL,
    VotedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (PhotoID) REFERENCES superadmin5.Photos(PhotoID),
    FOREIGN KEY (VotedBy) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.WhatIfIdeas (
    IdeaID UUID PRIMARY KEY,
    Description TEXT NOT NULL,
    SubmittedBy UUID NOT NULL,
    SubmittedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    Phase TEXT NOT NULL,
    IsSelected BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (SubmittedBy) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.WhatIfVotes (
    VoteID UUID PRIMARY KEY,
    IdeaID UUID NOT NULL,
    VotedBy UUID NOT NULL,
    VotedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (IdeaID) REFERENCES superadmin5.WhatIfIdeas(IdeaID),
    FOREIGN KEY (VotedBy) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.Quizzes (
    QuizID UUID PRIMARY KEY,
    HostID UUID NOT NULL,
    Passcode TEXT NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (HostID) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.QuizQuestions (
    QuestionID UUID PRIMARY KEY,
    QuizID UUID NOT NULL,
    QuestionText TEXT NOT NULL,
    CorrectOption INT NOT NULL,
    SequenceNumber INT NOT NULL,
    FOREIGN KEY (QuizID) REFERENCES superadmin5.Quizzes(QuizID)
);

CREATE TABLE superadmin5.QuizAnswers (
    AnswerID UUID PRIMARY KEY,
    QuestionID UUID NOT NULL,
    UserID UUID NOT NULL,
    SelectedOption INT NOT NULL,
    AnsweredAt TIMESTAMP NOT NULL DEFAULT NOW(),
    IsCorrect BOOLEAN NOT NULL,
    TimeTaken INT,
    FOREIGN KEY (QuestionID) REFERENCES superadmin5.QuizQuestions(QuestionID),
    FOREIGN KEY (UserID) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.QuizParticipants (
    QuizID UUID NOT NULL,
    UserID UUID NOT NULL,
    JoinTime TIMESTAMP NOT NULL DEFAULT NOW(),
    TotalScore INT NOT NULL DEFAULT 0,
    PRIMARY KEY (QuizID, UserID),
    FOREIGN KEY (QuizID) REFERENCES superadmin5.Quizzes(QuizID),
    FOREIGN KEY (UserID) REFERENCES superadmin5.Users(UserID)
);

CREATE TABLE superadmin5.BirthdayWishes (
    WishID UUID PRIMARY KEY,
    UserID UUID NOT NULL,
    WishedBy UUID NOT NULL,
    Message TEXT,
    Emoji TEXT,
    SentAt TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (UserID) REFERENCES superadmin5.Users(UserID),
    FOREIGN KEY (WishedBy) REFERENCES superadmin5.Users(UserID)
);
CREATE TABLE superadmin5.GuessImages (
    ImageID UUID PRIMARY KEY,
    ImageURL TEXT NOT NULL,
    Title TEXT NOT NULL,
    Description TEXT,
    Category guess_image_category NOT NULL,
    UploadedAt TIMESTAMP DEFAULT NOW(),
    UploadedBy UUID, -- could reference a user ID from GlobalUsers
    IsCompleted BOOLEAN DEFAULT FALSE
);
