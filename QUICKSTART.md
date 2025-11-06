# Quick Start Guide

## Step 1: Database Setup

1. **Create PostgreSQL database:**
   ```bash
   createdb fpscore
   ```
   Or using psql:
   ```sql
   CREATE DATABASE fpscore;
   ```

2. **Run the schema:**
   ```bash
   psql -d fpscore -f schema.sql
   ```

3. **(Optional) Add sample data:**
   ```bash
   psql -d fpscore -f seed-data.sql
   ```

## Step 2: Configure Database Connection

Set your database connection string:

**Windows (PowerShell):**
```powershell
$env:DATABASE_URL="postgres://postgres:postgres@localhost/fpscore?sslmode=disable"
```

**Linux/Mac:**
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost/fpscore?sslmode=disable"
```

Or create a `.env` file (you'll need to load it manually or use a library like `godotenv`).

## Step 3: Run the Application

```bash
go run main.go handlers.go
```

The server will start on `http://localhost:3000` (or the port specified in the `PORT` environment variable).

## Step 4: Add Assessment Questions

You need to populate the database with questions from your scoring key documents. Here's the structure:

### Example: Adding Counselling Questions

```sql
-- 1. Add thematic areas for Counselling (assessment_type_id = 1)
INSERT INTO thematic_areas (assessment_type_id, name, display_order) VALUES
    (1, 'Maintains privacy and confidentiality', 1),
    (1, 'Provides comprehensive and correct information on service options', 2),
    (1, 'Explains how the chosen service would be provided', 3),
    (1, 'Provides information about other SRHR services', 4),
    (1, 'Assesses the client''s medical eligibility for the chosen method', 5)
ON CONFLICT DO NOTHING;

-- 2. Add questions for first thematic area
INSERT INTO questions (thematic_area_id, question_text, score_weight, is_critical, is_important, display_order) VALUES
    ((SELECT id FROM thematic_areas WHERE name = 'Maintains privacy and confidentiality' LIMIT 1),
     'Greets and employs a client-centred style of communication when speaking to clients', 2, false, false, 1),
    ((SELECT id FROM thematic_areas WHERE name = 'Maintains privacy and confidentiality' LIMIT 1),
     'Uses language the client is comfortable with', 2, false, false, 2),
    ((SELECT id FROM thematic_areas WHERE name = 'Maintains privacy and confidentiality' LIMIT 1),
     'Follows a structured counselling approach like REDI (Rapport, Explore, Decide and Implement)', 5, false, true, 3),
    ((SELECT id FROM thematic_areas WHERE name = 'Maintains privacy and confidentiality' LIMIT 1),
     'Asks client about the service(s) they are seeking and if they have something specific in mind', 2, false, false, 4)
ON CONFLICT DO NOTHING;
```

### Score Weights:
- **Critical (Bold)**: `score_weight = 10`, `is_critical = true`
- **Important (Asterisk *)**: `score_weight = 5`, `is_important = true`
- **Normal**: `score_weight = 2`, both flags = false

## Step 5: Access the Application

1. Open your browser and go to `http://localhost:3000`
2. Select a facility (Region → District → Subcounty → Facility)
3. Choose an assessment type
4. Fill out the assessment form
5. Submit and view results

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL is running
- Check your connection string format: `postgres://username:password@host:port/database?sslmode=disable`
- Verify the database exists: `psql -l | grep fpscore`

### Port Already in Use
- Change the port: `$env:PORT="8080"` (Windows) or `export PORT=8080` (Linux/Mac)
- Or stop the process using port 3000

### Missing Questions
- The application will work even with empty assessment types
- Add questions using SQL INSERT statements as shown above
- Questions should be organized by thematic areas

## Next Steps

1. **Populate all assessment types** with questions from your scoring key documents
2. **Add your actual facilities** by inserting into regions, districts, subcounties, and facilities tables
3. **Customize the UI** by modifying the HTML/CSS files in the `public` directory
4. **Deploy** to your production server

## API Testing

You can test the API using curl or any HTTP client:

```bash
# Get regions
curl http://localhost:3000/api/regions

# Get assessment types
curl http://localhost:3000/api/assessment-types

# Get assessments
curl http://localhost:3000/api/assessments
```

## Support

Refer to the main README.md for more detailed information.

