-- Create schema
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
