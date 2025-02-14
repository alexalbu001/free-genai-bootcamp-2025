# Backend Technical Spec

## Business Logic

A language learning school wants to build a prototype of learning portal which will act as three things:
Inventory of possible vocabulary that can be learned
Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
A unified launchpad to launch different learning apps

## Technical Restrictions:
Use SQLite3 as the database
You can use any language or framework 
Does not require authentication/authorization, assume there is a single user
The backend will be written in GO
The API will be built using GIN and return JSON
MAGE is  task runner for GO

## Directory Structure
```text
backend_go/
├── cmd/
│   └── server/
├── internal/
│   ├── models/     # Data structures and database operations
│   ├── handlers/   # HTTP handlers organized by feature (dashboard, words, groups, etc.)
│   └── service/    # Business logic
├── db/
│   ├── migrations/ # Database migration files
│   └── seeds/      # Seed data files
├── magefile.go
├── go.mod
└── words.db


## Database Schema
Our DB will be a single sqlite database called `words.db` that will be in the root of the project folder of `backend-go`
The following tables:
words
- `id` (Primary Key): Unique identifier for each word
- `parts` (JSON, Required): Word components stored in JSON format

groups — Manages collections of words.
- `id` (Primary Key): Unique identifier for each group
- `name` (String, Required): Name of the group
- `words_count` (Integer, Default: 0): Counter cache for the number of words in the group

word_groups — join-table enabling many-to-many relationship between words and groups.
- `word_id` (Foreign Key): References words.id
- `group_id` (Foreign Key): References groups.id

study_activities — Defines different types of study activities available.
- `id` (Primary Key): Unique identifier for each activity
- `name` (String, Required): Name of the activity (e.g., "Flashcards", "Quiz")
- `url` (String, Required): The full URL of the study activity

study_sessions — Records individual study sessions.
- `id` (Primary Key): Unique identifier for each session
- `group_id` (Foreign Key): References groups.id
- `study_activity_id` (Foreign Key): References study_activities.id
- `created_at` (Timestamp, Default: Current Time): When the session was created

word_review_items — Tracks individual word reviews within study sessions.
- `id` (Primary Key): Unique identifier for each review
- `word_id` (Foreign Key): References words.id
- `study_session_id` (Foreign Key): References study_sessions.id
- `correct` (Boolean, Required): Whether the answer was correct
- `created_at` (Timestamp, Default: Current Time): When the review occurred

## Relationships

word belongs to groups through  word_groups
group belongs to words through word_groups
session belongs to a group
session belongs to a study_activity
session has many word_review_items
word_review_item belongs to a study_session
word_review_item belongs to a word



## API

#### GET /api/dashboard/last_study_session
Example response:

```json
{
    "id": 123,
    "group_id": 456,
    "study_activity_id": 789,
    "created_at": "2024-03-20T15:30:00Z",
    "group_id": 456,
    "group_name": "Basic Verbs"
}
```

#### GET /api/dashboard/study_progress
Example response:

```json
{
  "total_words_studied": 3,
  "total_available_words": 124,
}
```

#### GET /api/dashboard/quick_stats
Example response:

```json
{
  "success_rate": 80.0,
  "total_study_sessions": 4,
  "total_active_groups": 3,
  "study_streak_days": 4
}
```
#### GET /api/study_activities/:id
Example response:

```json
{
  "id": 1,
  "name": "Vocabulary Quiz",
  "thumbnail_url": "https://example.com/thumbnail.jpg",
  "description": "Practice your vocabulary with flashcards"
}
```

#### GET /api/study_activities/:id/study_sessions
Example response:

{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

#### GET /api/study_sessions/:id/word_review_items
Example response:

```json
{
  "id": 123,
  "activity_name": "Vocabulary Quiz",
  "group_name": "Basic Greetings",
  "start_time": "2025-02-08T17:20:23-05:00",
  "end_time": "2025-02-08T17:30:23-05:00",
  "review_items_count": 20
}
```

#### POST /api/study_activities/
Required parameters: group_id, study_activity_id


```json
{
  "group_id": 456,
  "study_activity_id": 789
}
```

#### GET /api/words
Example response:

```json
{
  "items": [
    {
      "french": "bonjour",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 500,
    "items_per_page": 100
  }
}
```

#### GET /api/words/:id
Example response:

```json
{
  "french": "bonjour",
  "english": "hello",
  "stats": {
    "correct_count": 5,
    "wrong_count": 2
  },
  "groups": [
    {
      "id": 1,
      "name": "Basic Greetings"
    }
  ]
}
```

#### GET /api/groups
Example response:

```json
{
  "data": [
    {
      "id": 456,
      "name": "Basic Verbs",
      "words_count": 50,
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 10,
    "items_per_page": 100
  }
}
```

#### GET /api/groups/:id
Example response:

```json
{
  "id": 1,
  "name": "Basic Greetings",
  "stats": {
    "total_word_count": 20
  }
}
```

#### GET /api/groups/:id/words
Example response:

```json
{
  "items": [
    {
      "french": "bonjour",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 20,
    "items_per_page": 100
  }
}
```

#### GET /api/words
Note: This endpoint returns words with pagination.

Example response:

```json
{
  "items": [
    {
      "french": "bonjour",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 500,
    "items_per_page": 100
  }
}
```

#### GET /api/groups/:id/study_sessions
Example response:

```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 5,
    "items_per_page": 100
  }
}
```

#### GET /api/study_sessions
Example response:

```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 100
  }
}
```

#### GET /api/study_sessions/:id
Example response:

```json
{
  "id": 123,
  "activity_name": "Vocabulary Quiz",
  "group_name": "Basic Greetings",
  "start_time": "2025-02-08T17:20:23-05:00",
  "end_time": "2025-02-08T17:30:23-05:00",
  "review_items_count": 20
}
```

#### GET /api/study_sessions/:id/words
Example response:

```json
{
  "items": [
    {
      "french": "bonjour",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 20,
    "items_per_page": 100
  }
}
```

#### POST /api/reset_history
Example response:

```json
{
  "message": "All study history has been reset",
  "success": true
}
```

#### POST /api/full_reset
Example response:

```json
{
  "message": "System reset complete",
  "success": true
}
```

#### POST /api/study_sessions/:id/word/:word_id/review
Example request body:

```json
{
  "correct": true
}
```

Example response:

```json
{
  "success": true,
  "word_id": 1,
  "study_session_id": 123,
  "correct": true,
  "created_at": "2025-02-08T17:33:07-05:00"
}
```

## Mage (Tasks)
Mage is a task runner that will be used to run the scripts to initialise the database and reset the database.
### Initialise Database
This task will initialise the sqlite DB called `words.db` 
### Migrate Database
This task will run a series of migrations sql files on the DB
Migrations live in the migrations folder. The migration files will be run in order of their file name. The file names should looks like this:

0001_init.sql
0002_create_words_table.sql

### Seed Data
This task will import json files and transform them into target data for our database

All seed files live in the seeds folder.

In our task we should have DSL to specific each seed file and its expected group word name.
```Json
[
  {
    "french": "bonjour",
    "english": "hello",
  },
  ...
]
