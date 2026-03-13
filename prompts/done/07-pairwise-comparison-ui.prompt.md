# Prompt 7: Pairwise Comparison UI

Build the pairwise comparison grid interface showing feature matchups, voting buttons, real-time consensus indicators, and progress tracking through all required comparisons.

## Requirements
- Create pairwise comparison grid component
- Implement voting interface for each comparison
- Add real-time WebSocket integration for live updates
- Show consensus status and progress tracking
- Handle both Value and Complexity comparison sessions
- Display attendee voting status for each comparison

## Components to Create

### PairwiseGrid Component
- Grid layout showing all feature comparisons
- Vote buttons for each attendee
- Real-time consensus indicators
- Progress bar for session completion

### ComparisonCard Component
- Display two features side-by-side
- Voting buttons (Feature A, Tie, Feature B)
- Show individual attendee votes
- Consensus status indicator
- Feature details on hover/expand

### SessionProgress Component
- Progress bar showing completion percentage
- Count of completed vs. total comparisons
- List of attendees and their voting progress
- Session timer and status

### AttendeeVotingPanel Component
- Current attendee's voting interface
- Clear indication of unvoted comparisons
- Easy navigation between comparisons
- Submit/change vote functionality

## WebSocket Integration
- Connect to session WebSocket on page load
- Real-time vote updates from other attendees
- Live consensus status changes
- Session progress synchronization
- Attendee presence indicators

## User Experience Features
- Keyboard shortcuts for quick voting
- Visual indicators for voted/unvoted comparisons
- Smooth animations for status changes
- Mobile-responsive design
- Clear instructions and help text

## State Management
- Local voting state with optimistic updates
- WebSocket event handling
- Error recovery for failed votes
- Session synchronization on reconnect

## Navigation Flow
1. Select criterion type (Value or Complexity)
2. Start pairwise comparison session
3. Vote on all feature pairs
4. Wait for group consensus on all comparisons
5. Proceed to Fibonacci scoring or view results