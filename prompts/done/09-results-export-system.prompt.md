# Prompt 9: Results & Export System

Create the final results page showing calculated rankings, Final Priority Scores, and export functionality to CSV/JSON formats. Add Jira-compatible story export with Fibonacci complexity assignments.

## Requirements
- Calculate and display Final Priority Scores using P-WVC formula
- Create ranked feature list with all scoring details
- Implement multiple export formats (CSV, JSON, Jira)
- Add data visualization for results
- Include session summary and methodology explanation
- Allow results filtering and sorting

## Database Schema Addition
```sql
-- Final priority calculations
CREATE TABLE priority_calculations (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    feature_id INTEGER REFERENCES features(id) ON DELETE CASCADE,
    w_value DECIMAL(10,6) NOT NULL,        -- Win-count weight for value
    w_complexity DECIMAL(10,6) NOT NULL,   -- Win-count weight for complexity  
    s_value INTEGER NOT NULL,              -- Fibonacci score for value
    s_complexity INTEGER NOT NULL,         -- Fibonacci score for complexity
    weighted_value DECIMAL(10,6) NOT NULL, -- SValue × WValue
    weighted_complexity DECIMAL(10,6) NOT NULL, -- SComplexity × WComplexity
    final_priority_score DECIMAL(10,6) NOT NULL, -- Weighted Value ÷ Weighted Complexity
    rank INTEGER NOT NULL,
    calculated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(project_id, feature_id)
);
```

## Components to Create

### ResultsRanking Component
- Ranked list of features by Final Priority Score
- Expandable rows showing calculation details
- Sort and filter capabilities
- Export options

### CalculationDetails Component
- Show P-WVC calculation breakdown for each feature
- Win-count weights from pairwise comparisons
- Fibonacci scores from consensus
- Final Priority Score calculation
- Mathematical formula explanation

### ExportOptions Component
- Multiple export format buttons
- CSV download for spreadsheet analysis
- JSON export for API integration
- Jira-compatible format with story points

### ResultsVisualization Component
- Bar chart of Final Priority Scores
- Scatter plot of Value vs Complexity
- Priority quadrant visualization
- Interactive charts with feature details

## API Endpoints to Add
- `POST /api/projects/{id}/calculate-results` - Trigger P-WVC calculation
- `GET /api/projects/{id}/results` - Get calculated rankings
- `GET /api/projects/{id}/results/export?format=csv|json|jira` - Export results

## Export Formats

### CSV Export
```csv
rank,feature_title,description,final_priority_score,s_value,s_complexity,w_value,w_complexity
1,"User Login","Authentication system",5.25,8,3,0.75,0.88
2,"Dashboard","Main dashboard view",4.12,5,2,0.88,0.63
```

### Jira Export Format
```json
{
  "issues": [
    {
      "summary": "User Login",
      "description": "Authentication system\n\nAcceptance Criteria:\n- Users can login with email/password",
      "storyPoints": 3,
      "priority": "High",
      "customFields": {
        "finalPriorityScore": 5.25,
        "valueScore": 8,
        "complexityScore": 3
      }
    }
  ]
}
```

## Calculation Logic
1. Retrieve win-count weights from pairwise comparison results
2. Get consensus Fibonacci scores for Value and Complexity
3. Calculate: FPS = (SValue × WValue) / (SComplexity × WComplexity)
4. Rank features by descending Final Priority Score
5. Store results in priority_calculations table

## User Experience
- Clear visual hierarchy showing top-priority features
- Expandable sections for calculation transparency  
- Multiple visualization options
- One-click export to various formats
- Session summary with methodology recap
- Ability to restart scoring if needed