import { render } from '@testing-library/react';
import KPIPage from '@/app/[locale]/kpi/page';
import * as authUtils from '@/features/auth/utils/getCurrentUser';

vi.mock('@/features/auth/utils/getCurrentUser');

describe('KPI route RBAC', () => {
  it('renders rep view for sales_rep role', async () => {
    (authUtils.getCurrentUser as any).mockResolvedValue({ id: 'u1', role: 'sales_rep' });
    const { container } = render(await KPIPage());
    expect(container).toBeTruthy();
  });

  it('renders manager view for sales_manager role', async () => {
    (authUtils.getCurrentUser as any).mockResolvedValue({ id: 'm1', role: 'sales_manager' });
    const { container } = render(await KPIPage());
    expect(container).toBeTruthy();
  });

  it('redirects anonymous to login', async () => {
    (authUtils.getCurrentUser as any).mockResolvedValue(null);
    await expect(KPIPage()).rejects.toThrow();
  });
});
