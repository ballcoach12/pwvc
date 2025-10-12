# Prompt 6: React Frontend Foundation

Create React application with routing, project setup pages, attendee management interface, and feature input/import functionality. Use Material-UI for consistent design system.

## Requirements
- Set up React application with Vite build tool
- Configure routing with React Router
- Implement Material-UI design system
- Create project setup and management pages
- Build attendee management interface
- Add feature input and CSV import functionality
- Set up API service layer for backend communication

## Project Structure
```
web/
├── src/
│   ├── components/
│   │   ├── Layout/
│   │   ├── ProjectCard/
│   │   ├── AttendeeList/
│   │   ├── FeatureForm/
│   │   └── FileUpload/
│   ├── pages/
│   │   ├── ProjectList/
│   │   ├── ProjectSetup/
│   │   ├── AttendeeManagement/
│   │   └── FeatureManagement/
│   ├── services/
│   │   ├── api.js
│   │   ├── projectService.js
│   │   └── featureService.js
│   ├── hooks/
│   │   ├── useProject.js
│   │   └── useFeatures.js
│   └── utils/
│       └── csvParser.js
```

## Pages to Create

### ProjectList Page
- Display all user's projects
- Create new project button
- Project cards with basic info
- Navigation to project details

### ProjectSetup Page  
- Project creation form (name, description)
- Edit project details
- Navigation to attendee management

### AttendeeManagement Page
- Add/remove attendees
- Set facilitator roles
- Attendee list with roles
- Continue to feature management

### FeatureManagement Page
- Add individual features form
- CSV import functionality
- Feature list with edit/delete actions
- Export features to CSV
- Continue to scoring phases

## Key Components
- Responsive navigation layout
- Form validation and error handling
- Loading states and error messages
- File upload with CSV validation
- Confirmation dialogs for deletions

## API Integration
- Axios-based service layer
- Error handling for API failures
- Loading state management
- Success/error notifications using Material-UI snackbars