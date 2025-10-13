import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import PairwiseGrid from '../PairwiseGrid/PairwiseGrid';

const mockFeatures = [
  { id: 1, title: 'Feature A', description: 'First feature' },
  { id: 2, title: 'Feature B', description: 'Second feature' },
  { id: 3, title: 'Feature C', description: 'Third feature' },
];

const mockComparisons = [
  {
    id: 1,
    featureA: mockFeatures[0],
    featureB: mockFeatures[1],
    winner: null,
    consensusReached: false,
    votes: [],
  },
  {
    id: 2,
    featureA: mockFeatures[0],
    featureB: mockFeatures[2],
    winner: null,
    consensusReached: false,
    votes: [],
  },
];

const mockAttendees = [
  { id: 1, name: 'John Doe', isFacilitator: true },
  { id: 2, name: 'Jane Smith', isFacilitator: false },
];

const MockPairwiseGrid = ({ 
  features = mockFeatures,
  comparisons = mockComparisons,
  attendees = mockAttendees,
  currentAttendeeId = 1,
  onVote = vi.fn(),
  criterionType = 'value'
}) => {
  return (
    <BrowserRouter>
      <PairwiseGrid
        features={features}
        comparisons={comparisons}
        attendees={attendees}
        currentAttendeeId={currentAttendeeId}
        onVote={onVote}
        criterionType={criterionType}
      />
    </BrowserRouter>
  );
};

describe('PairwiseGrid', () => {
  it('renders feature comparison grid', () => {
    render(<MockPairwiseGrid />);

    expect(screen.getByText('Feature A')).toBeInTheDocument();
    expect(screen.getByText('Feature B')).toBeInTheDocument();
    expect(screen.getByText('Feature C')).toBeInTheDocument();
  });

  it('shows comparison buttons for feature pairs', () => {
    render(<MockPairwiseGrid />);

    // Should show buttons for Feature A vs B comparison
    const comparisonButtons = screen.getAllByRole('button');
    expect(comparisonButtons.length).toBeGreaterThan(0);
  });

  it('calls onVote when comparison is selected', async () => {
    const mockOnVote = vi.fn();
    render(<MockPairwiseGrid onVote={mockOnVote} />);

    const voteButton = screen.getAllByRole('button')[0];
    fireEvent.click(voteButton);

    await waitFor(() => {
      expect(mockOnVote).toHaveBeenCalled();
    });
  });

  it('shows progress indicator', () => {
    render(<MockPairwiseGrid />);

    // Should show some form of progress (e.g., "0 of 2 comparisons completed")
    expect(screen.getByText(/comparison/i)).toBeInTheDocument();
  });

  it('displays different criterion types', () => {
    render(<MockPairwiseGrid criterionType="complexity" />);

    expect(screen.getByText(/complexity/i)).toBeInTheDocument();
  });

  it('shows consensus status for comparisons', () => {
    const completedComparisons = [
      {
        ...mockComparisons[0],
        winner: mockFeatures[0],
        consensusReached: true,
      },
    ];

    render(<MockPairwiseGrid comparisons={completedComparisons} />);

    expect(screen.getByText(/consensus/i)).toBeInTheDocument();
  });
});