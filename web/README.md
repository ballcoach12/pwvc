# PairWise Web Frontend

React-based frontend for the PairWise feature prioritization application.

## Features

- **Project Management**: Create and manage PairWise projects
- **Attendee Management**: Add team members and assign facilitator roles
- **Feature Management**: Add features manually or import from CSV
- **Material-UI Design**: Modern, responsive interface
- **Real-time Validation**: Form validation and error handling
- **CSV Import/Export**: Bulk feature management capabilities

## Technology Stack

- **React 18** - UI framework
- **Vite** - Build tool and dev server
- **Material-UI (MUI)** - Component library and design system
- **React Router** - Client-side routing
- **Axios** - HTTP client for API communication
- **Emotion** - CSS-in-JS styling

## Getting Started

### Prerequisites

- Node.js 16 or later
- Go backend server running on localhost:8080

### Installation

```bash
# Navigate to web directory
cd web/

# Install dependencies
npm install

# Start development server
npm run dev
```

The application will be available at `http://localhost:3000`.

### Development

```bash
# Start dev server with hot reload
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run linting
npm run lint
```

## Project Structure

```
src/
├── components/           # Reusable UI components
│   ├── Layout/          # App layout and navigation
│   ├── ProjectCard/     # Project display card
│   ├── AttendeeList/    # Attendee management list
│   ├── FeatureForm/     # Feature creation/editing form
│   ├── FileUpload/      # CSV file upload component
│   └── NotificationProvider/ # Toast notifications
├── pages/               # Page components
│   ├── ProjectList/     # Projects overview page
│   ├── ProjectSetup/    # Project creation/editing
│   ├── AttendeeManagement/ # Team member management
│   └── FeatureManagement/  # Feature input and import
├── services/            # API communication layer
│   ├── api.js          # Axios instance and interceptors
│   ├── projectService.js # Project-related API calls
│   └── featureService.js  # Feature-related API calls
├── hooks/               # Custom React hooks
│   ├── useProject.js    # Project state management
│   └── useFeatures.js   # Feature state management
├── utils/               # Utility functions
│   ├── csvParser.js     # CSV parsing and validation
│   └── helpers.js       # General helper functions
├── App.jsx             # Main app component
└── main.jsx            # React app entry point
```

## Pages and Flow

### 1. Project List

- View all projects
- Create new project
- Edit/delete existing projects
- Navigate to project management

### 2. Project Setup

- Create/edit project details
- Set project name and description
- Navigate to attendee management

### 3. Attendee Management

- Add team members with email and role
- Assign facilitator permissions
- Remove attendees
- Navigate to feature management

### 4. Feature Management

- Add features manually via form
- Import features from CSV file
- View and edit feature list
- Export features to CSV
- Navigate to pairwise comparison (coming soon)

## API Integration

The frontend communicates with the Go backend API through:

- **Proxy Configuration**: Vite dev server proxies `/api/*` to `localhost:8080`
- **Axios Interceptors**: Handle authentication and error responses
- **Service Layer**: Organized API calls by domain (projects, features, etc.)
- **Error Handling**: Centralized error handling with user-friendly messages

## CSV Import Format

Features can be imported via CSV with the following format:

```csv
name,description
User Authentication,Login and registration system
Dashboard Analytics,Real-time analytics dashboard
Mobile App,iOS and Android mobile application
```

- **Required**: `name` column
- **Optional**: `description` column
- **Validation**: Names must be 3-100 characters, descriptions max 500 characters

## Styling and Theming

- **Material-UI Theme**: Custom theme with consistent colors and typography
- **Responsive Design**: Mobile-first approach with Material-UI's grid system
- **Component Overrides**: Custom styling for cards, buttons, and form elements
- **Dark/Light Mode**: Ready for theme switching (not yet implemented)

## State Management

- **React Hooks**: useState, useEffect for local component state
- **Custom Hooks**: Reusable stateful logic for projects, features, attendees
- **Context**: Notification system for app-wide toast messages
- **Local Storage**: Potential for caching user preferences

## Development Guidelines

- **Component Structure**: Each component in its own directory with index file
- **Error Boundaries**: Graceful error handling throughout the app
- **Loading States**: Visual feedback for async operations
- **Form Validation**: Client-side validation with server-side backup
- **Accessibility**: ARIA labels and semantic HTML where possible

## Future Enhancements

- **Pairwise Comparison Interface**: Interactive comparison grid
- **Fibonacci Scoring Interface**: Visual scoring with Fibonacci sequence
- **Results Dashboard**: Charts and graphs for prioritization results
- **Real-time Collaboration**: WebSocket integration for live updates
- **Export Options**: Multiple export formats (PDF, Excel, etc.)
- **User Authentication**: Login system and user management
- **Project Templates**: Pre-configured project types
- **Mobile Optimization**: Enhanced mobile experience

## Environment Variables

Create a `.env` file for environment-specific configuration:

```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_TITLE=PairWise
```

## Building for Production

```bash
# Build optimized production bundle
npm run build

# The dist/ directory contains the built application
# Serve with any static file server
```

## Contributing

1. Follow the existing code structure and naming conventions
2. Add appropriate error handling and loading states
3. Include form validation for user inputs
4. Test on both desktop and mobile viewports
5. Add comments for complex logic or business rules

## Known Issues

- CSV import currently uses basic parsing (consider upgrading to a robust CSV library)
- File upload is limited to 5MB (configurable)
- No offline support yet
- Limited browser compatibility testing

## Support

For issues related to the frontend, check:

1. Browser console for JavaScript errors
2. Network tab for API communication issues
3. React DevTools for component state debugging
4. Backend server logs for API errors
