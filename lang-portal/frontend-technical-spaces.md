# Frontend Technical Spec
## pages
### Dashboard `/dashboard`
#### Purpose
The purpose of this page is to provide a quick overview of the user's progress and activity. Will be the default page on load.
#### Components
- Last study session
 shows last activity used
 shows last study session date
 summarizes wrong vs correct from last activity
 has a link to the group

- Study progress
total words studied across all words in database eg 3/124
display a mastery progress bar eg 0%
- Quick stats
success rate
total study sessions
total active groups
study streak
- Start Studying button
goes to study activities page

#### API endpoints
GET /api/dashboard/last_study_session
GET /api/dashboard/study_progress
GET /api/dashboard/quick_stats
POST /api/dashboard/start_studying

### Study Activities `/study-activities`
#### Purpose
The purpose of this page is to provide a list of study activities to the user with a thumbnail and name to either launch or view more information about the activity.

#### Components
- Study Activity Card
    - Thumbnail of the study activity
    - Name of the study activity
    - Description of the study activity
    - Launch button
    - View more button

#### API endpoints
GET /api/study_activities

### Study activity show `/study-activity/:id`
#### Purpose
The purpose of this page is to provide a detailed view of the study activity including the description, thumbnail, and its past activity
#### Components
    - Thumbnail of the study activity
    - Name of the study activity
    - Description of the study activity
    - Launch button
    - Study activity paginated list
        - id
        - activity name
        -group name
        - start time
        -end time ( inferred by the last word_review_item )
        - number of review items

#### API endpoints
GET /api/study_activities/:id
GET /api/study_activities/:id/study_sessions


### Study activities launch `/study-activity/:id/launch`
#### Purpose
The purpose of this page is to launch the study activity 
#### Components

    - Thumbnail of the study activity
    - Name of the study activity
    -Launch form
        - select field for group
        - start button

## Behavior
After the form is submitted a new tab opens with the study activity url provided in the db
Also the after form is submitted the page will redirect to the study session show page 

#### API endpoints
POST /api/study_activities


### words `/words`
#### Purpose
The purpose of this page is to provide a list of all words in the database
#### Components
- paginated list of words
    -Columns
        - French word
        - English word
        - Correct count
        - Wrong count
    - Pagination with 100 items per page
    - Clicking the french word will take us tho the show word page

#### API endpoints  
GET /api/words

### word show `/words/:id`
#### Purpose
The purpose of this page is to provide a detailed view of the word i
#### Components
- French word
- English word
- Study statistics
    - Correct count
    - Wrong count
- Word groups
    - show an a series of pills eg. tags
    - when a group name is clicked it will take us to the group show page

#### Api endpoints
GET /api/words/:id
