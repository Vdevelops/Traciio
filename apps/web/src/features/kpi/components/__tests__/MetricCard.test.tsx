import { render, screen } from '@testing-library/react';
import MetricCard from '../MetricCard';

describe('MetricCard', () => {
  it('renders no-data state when value is null', () => {
    render(<MetricCard label="Conversion Rate" value={null} suffix="%" />);
    expect(screen.getByText(/Belum ada data/i)).toBeInTheDocument();
  });

  it('renders numeric value with suffix', () => {
    render(<MetricCard label="Conversion Rate" value={25} suffix="%" />);
    expect(screen.getByText(/25/)).toBeInTheDocument();
    expect(screen.getByLabelText(/metric-value-conversion-rate/i)).toBeInTheDocument();
  });
});
