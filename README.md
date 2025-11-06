# Family Planning Score Tool

A comprehensive digital assessment system for family planning services, built with Go (Fiber), PostgreSQL, HTML, and Bootstrap.

## Features

- **9 Assessment Types**: Counselling, OCPs, Injectable PO, Implant Insertion, Implant Removal, IUD Insertion, IUD Removal, Mini Lap, and Vasectomy
- **Mobile-Friendly UI**: Accordion-based interface optimized for phone use
- **Cascading Location Selection**: Region → District → Subcounty → Facility
- **Real-time Scoring**: Automatic calculation of scores per thematic area and overall assessment
- **Performance Levels**: Proficient (>90%), Competent (70-89%), Not Acceptable (<70%)
- **Assessment History**: View past assessments with detailed breakdowns
- **Score Analysis**: See what contributed positively or negatively to scores

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Node.js (optional, for development)

## Setup

### 1. Database Setup

Create a PostgreSQL database:

```bash
createdb fpscore
```

Or using psql:

```sql
CREATE DATABASE fpscore;
```

### 2. Run Database Schema

```bash
psql -d fpscore -f schema.sql
```

### 3. Configure Database Connection

Set your database connection string as an environment variable:

**Windows (PowerShell):**
```powershell
$env:DATABASE_URL="postgres://postgres:postgres@localhost/fpscore?sslmode=disable"
```

**Linux/Mac:**
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost/fpscore?sslmode=disable"
```

The application will use the `DATABASE_URL` environment variable, or default to `postgres://postgres:postgres@localhost/fpscore?sslmode=disable` if not set.

**Optional:** Set the port:
```bash
export PORT=3000  # Defaults to 3000
```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Seed Assessment Questions

Populate all assessment questions and thematic areas:

```bash
psql -d fpscore -f seed-questions.sql
```

### 6. Run the Application

**Option 1: Using go run (recommended for development)**
```bash
go run ./cmd
```

**Option 2: Build and run**
```bash
go build -o bin/fpscore.exe ./cmd
./bin/fpscore.exe
```

**Option 3: Using Make (if available on Windows with Make installed)**
```bash
make build
./bin/fpscore.exe
```

Or simply:
```bash
make run
```

**Note:** Make sure you're running from the project root directory, not from inside the `cmd` folder.

The application will be available at `http://localhost:3000`

## Database Structure

### Geographic Hierarchy
- `regions` - Top-level geographic division
- `districts` - Second-level, linked to regions
- `subcounties` - Third-level, linked to districts
- `facilities` - Health facilities, linked to subcounties

### Assessment Structure
- `assessment_types` - Types of assessments (Counselling, OCPs, etc.)
- `thematic_areas` - Sections within each assessment type
- `questions` - Individual questions with score weights (2, 5, or 10 points)
- `assessments` - Main assessment records
- `assessment_responses` - Individual question responses
- `thematic_area_scores` - Pre-calculated scores per thematic area

## API Endpoints

### Geographic Hierarchy
- `GET /api/regions` - Get all regions
- `GET /api/districts/:regionId` - Get districts in a region
- `GET /api/subcounties/:districtId` - Get subcounties in a district
- `GET /api/facilities/:subcountyId` - Get facilities in a subcounty

### Assessment Types
- `GET /api/assessment-types` - Get all assessment types
- `GET /api/assessment-types/:typeId/thematic-areas` - Get thematic areas for an assessment type
- `GET /api/thematic-areas/:thematicAreaId/questions` - Get questions for a thematic area

### Assessments
- `POST /api/assessments` - Create a new assessment
- `GET /api/assessments` - Get all assessments
- `GET /api/assessments/:id` - Get a specific assessment
- `GET /api/assessments/:id/summary` - Get assessment summary with good/bad contributions

## Usage

1. **Select Facility**: Choose Region → District → Subcounty → Facility
2. **Choose Assessment Type**: Click on the assessment type card
3. **Fill Assessment Form**:
   - Enter assessor and client names (optional)
   - Answer questions in each thematic area accordion
   - Use Yes/No/NA buttons for each question
   - Monitor progress in the sidebar
4. **Submit**: Review and submit the assessment
5. **View Results**: See detailed score breakdown, performance level, and contributions

## Scoring System

- **Critical Steps** (Bold): 10 points
- **Important Steps** (Asterisk *): 5 points
- **Recommended Steps** (Normal): 2 points

Final score = (Achieved Score / Total Possible Score) × 100%

Performance Levels:
- **Proficient**: >90%
- **Competent**: 70-89%
- **Not Acceptable**: <70%

## Adding Assessment Data

To add assessment questions, insert data into the database:

```sql
-- Example: Add a thematic area
INSERT INTO thematic_areas (assessment_type_id, name, display_order)
VALUES (1, 'Maintains privacy and confidentiality', 1);

-- Example: Add a question
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order)
VALUES (1, 'Greets and employs a client-centred style of communication', 2, false, false, 1);
```

## Development

### Project Structure

```
fpscore/
├── cmd/
│   └── main.go         # Application entry point
├── config/
│   └── config.go       # Configuration management
├── database/
│   └── database.go     # Database connection
├── handlers/
│   └── handlers.go     # API handlers
├── models/
│   └── models.go       # Data models
├── schema.sql          # Database schema
├── seed-questions.sql  # All assessment questions
├── seed-data.sql       # Sample geographic data
├── go.mod              # Go dependencies
├── Makefile            # Build commands
├── public/             # Frontend files
│   ├── index.html
│   ├── assessment.html
│   ├── view-assessment.html
│   ├── app.js
│   ├── assessment.js
│   └── view-assessment.js
└── README.md
```

## License

This project is for internal use.

## Support

For issues or questions, please contact the development team.

