import { render, screen } from '@testing-library/react';
import CompositeScoreCard from '../CompositeScoreCard';

describe('CompositeScoreCard', () => {
  it('shows empty state when score is null', () => {
    render(
      <CompositeScoreCard
        score={null}
        grade={null}
        trend={null}
      />
    );
    expect(screen.getByText(/Tidak cukup data untuk evaluasi/i)).toBeInTheDocument();
  });

  it('shows score and grade when provided', () => {
    render(
      <CompositeScoreCard
        score={83.5}
        grade={'Good'}
        trend={{ previousCompositeScore: 78, delta: 5.5, direction: 'up' }}
      />
    );
    expect(screen.getByText(/83.5/)).toBeInTheDocument();
    expect(screen.getByText(/Good/)).toBeInTheDocument();
    expect(screen.getByText(/naik 5.5 poin dari periode lalu/i)).toBeInTheDocument();
  });
});
