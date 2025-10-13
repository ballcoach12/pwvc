import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { ErrorProvider } from '../../contexts/ErrorContext';
import ErrorDisplay from '../ErrorDisplay/ErrorDisplay';

// Mock error context value
const mockError = {
  id: '1',
  type: 'api',
  severity: 'error',
  message: 'Test error message',
  details: 'Test error details',
  timestamp: new Date(),
};

const MockErrorProvider = ({ children, error = null }) => {
  const value = {
    error,
    addError: () => {},
    clearError: () => {},
    clearAllErrors: () => {},
  };

  return (
    <ErrorProvider value={value}>
      {children}
    </ErrorProvider>
  );
};

describe('ErrorDisplay', () => {
  it('renders error message when error exists', () => {
    render(
      <MockErrorProvider error={mockError}>
        <ErrorDisplay />
      </MockErrorProvider>
    );

    expect(screen.getByText('Test error message')).toBeInTheDocument();
    expect(screen.getByText('Test error details')).toBeInTheDocument();
  });

  it('does not render when no error exists', () => {
    render(
      <MockErrorProvider error={null}>
        <ErrorDisplay />
      </MockErrorProvider>
    );

    expect(screen.queryByText('Test error message')).not.toBeInTheDocument();
  });

  it('renders different severity levels', () => {
    const warningError = {
      ...mockError,
      severity: 'warning',
      message: 'Warning message',
    };

    render(
      <MockErrorProvider error={warningError}>
        <ErrorDisplay />
      </MockErrorProvider>
    );

    expect(screen.getByText('Warning message')).toBeInTheDocument();
  });
});