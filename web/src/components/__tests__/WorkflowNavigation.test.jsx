import { fireEvent, render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import WorkflowNavigation from '../WorkflowNavigation/WorkflowNavigation';

const mockProgress = {
  projectId: 1,
  currentPhase: 'setup',
  setupCompleted: true,
  attendeesAdded: false,
  featuresAdded: false,
  pairwiseValueCompleted: false,
  pairwiseComplexityCompleted: false,
  fibonacciValueCompleted: false,
  fibonacciComplexityCompleted: false,
  resultsCalculated: false,
};

const MockWorkflowNavigation = ({ progress = mockProgress, onPhaseChange = vi.fn() }) => {
  return (
    <BrowserRouter>
      <WorkflowNavigation progress={progress} onPhaseChange={onPhaseChange} />
    </BrowserRouter>
  );
};

describe('WorkflowNavigation', () => {
  it('renders all workflow phases', () => {
    render(<MockWorkflowNavigation />);

    expect(screen.getByText(/Setup/i)).toBeInTheDocument();
    expect(screen.getByText(/Attendees/i)).toBeInTheDocument();
    expect(screen.getByText(/Features/i)).toBeInTheDocument();
    expect(screen.getByText(/Pairwise Value/i)).toBeInTheDocument();
    expect(screen.getByText(/Pairwise Complexity/i)).toBeInTheDocument();
    expect(screen.getByText(/Fibonacci Value/i)).toBeInTheDocument();
    expect(screen.getByText(/Fibonacci Complexity/i)).toBeInTheDocument();
    expect(screen.getByText(/Results/i)).toBeInTheDocument();
  });

  it('shows completed phases as checked', () => {
    render(<MockWorkflowNavigation />);

    const setupPhase = screen.getByText(/Setup/i).closest('.workflow-phase');
    expect(setupPhase).toHaveClass('completed');
  });

  it('shows current phase as active', () => {
    render(<MockWorkflowNavigation />);

    const currentPhase = screen.getByText(/Setup/i).closest('.workflow-phase');
    expect(currentPhase).toHaveClass('active');
  });

  it('disables future phases', () => {
    render(<MockWorkflowNavigation />);

    const futurePhase = screen.getByText(/Features/i).closest('.workflow-phase');
    expect(futurePhase).toHaveClass('disabled');
  });

  it('calls onPhaseChange when clicking available phase', () => {
    const mockOnPhaseChange = vi.fn();
    render(<MockWorkflowNavigation onPhaseChange={mockOnPhaseChange} />);

    const setupPhase = screen.getByText(/Setup/i);
    fireEvent.click(setupPhase);

    expect(mockOnPhaseChange).toHaveBeenCalledWith('setup');
  });

  it('shows progress percentage', () => {
    const progressWithMultipleCompleted = {
      ...mockProgress,
      setupCompleted: true,
      attendeesAdded: true,
      featuresAdded: true,
    };

    render(<MockWorkflowNavigation progress={progressWithMultipleCompleted} />);

    // Should show progress for 3 out of 8 phases completed (37.5%)
    expect(screen.getByText(/37%/i)).toBeInTheDocument();
  });
});