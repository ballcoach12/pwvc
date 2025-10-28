import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import FibonacciScorer from '../FibonacciScorer/FibonacciScorer';

const mockFeatures = [
  { id: 1, title: 'Feature A', description: 'First feature' },
  { id: 2, title: 'Feature B', description: 'Second feature' },
];

const mockScores = [
  { featureId: 1, attendeeId: 1, score: 5, scoreType: 'value' },
  { featureId: 2, attendeeId: 1, score: 3, scoreType: 'value' },
];

const MockFibonacciScorer = ({ 
  features = mockFeatures,
  scores = mockScores,
  currentAttendeeId = 1,
  scoreType = 'value',
  onScoreChange = vi.fn(),
  fibonacciValues = [1, 2, 3, 5, 8, 13, 21]
}) => {
  return (
    <FibonacciScorer
      features={features}
      scores={scores}
      currentAttendeeId={currentAttendeeId}
      scoreType={scoreType}
      onScoreChange={onScoreChange}
      fibonacciValues={fibonacciValues}
    />
  );
};

describe('FibonacciScorer', () => {
  it('renders all features for scoring', () => {
    render(<MockFibonacciScorer />);

    expect(screen.getByText('Feature A')).toBeInTheDocument();
    expect(screen.getByText('Feature B')).toBeInTheDocument();
  });

  it('shows Fibonacci value options', () => {
    render(<MockFibonacciScorer />);

    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
    expect(screen.getByText('8')).toBeInTheDocument();
    expect(screen.getByText('13')).toBeInTheDocument();
    expect(screen.getByText('21')).toBeInTheDocument();
  });

  it('calls onScoreChange when score is selected', () => {
    const mockOnScoreChange = vi.fn();
    render(<MockFibonacciScorer onScoreChange={mockOnScoreChange} />);

    const scoreButton = screen.getByText('8');
    fireEvent.click(scoreButton);

    expect(mockOnScoreChange).toHaveBeenCalledWith(
      expect.objectContaining({
        featureId: expect.any(Number),
        score: 8,
        scoreType: 'value'
      })
    );
  });

  it('shows current scores for features', () => {
    render(<MockFibonacciScorer />);

    // Should show that Feature A has score 5 and Feature B has score 3
    const scoreElements = screen.getAllByText(/score/i);
    expect(scoreElements.length).toBeGreaterThan(0);
  });

  it('displays different score types', () => {
    render(<MockFibonacciScorer scoreType="complexity" />);

    expect(screen.getByText(/complexity/i)).toBeInTheDocument();
  });

  it('shows progress indicator', () => {
    render(<MockFibonacciScorer />);

    // Should show some form of progress (e.g., "2 of 2 features scored")
    expect(screen.getByText(/feature/i)).toBeInTheDocument();
  });

  it('highlights selected score for each feature', () => {
    render(<MockFibonacciScorer />);

    // The score of 5 for Feature A should be highlighted/selected
    const selectedScore = screen.getByText('5');
    expect(selectedScore.closest('button')).toHaveClass('selected');
  });

  it('shows feature descriptions for context', () => {
    render(<MockFibonacciScorer />);

    expect(screen.getByText('First feature')).toBeInTheDocument();
    expect(screen.getByText('Second feature')).toBeInTheDocument();
  });
});