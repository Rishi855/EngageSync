CREATE TYPE guess_image_category AS ENUM ('movies', 'songs', 'actressAndActor');

CREATE TABLE kanaka.GuessImages (
    ImageID UUID PRIMARY KEY,
    ImageURL TEXT NOT NULL,
    Title TEXT NOT NULL,
    Description TEXT,
    Category guess_image_category NOT NULL,
    UploadedAt TIMESTAMP DEFAULT NOW(),
    UploadedBy UUID, -- could reference a user ID from GlobalUsers
    IsCompleted BOOLEAN DEFAULT FALSE
);