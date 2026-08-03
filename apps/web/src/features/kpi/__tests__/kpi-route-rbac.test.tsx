import { render } from '@testing-library/react';
import KPIPage from '@/app/[locale]/kpi/page';
import * as authUtils from '@/features/auth/utils/getCurrentUser';

vi.mock('@/features/auth/utils/getCurrentUser');

const pageProps = {
  params: Promise.resolve({ locale: 'en' }),
};

const mockedGetCurrentUser = vi.mocked(authUtils.getCurrentUser);

describe('KPI route RBAC', () => {
  it('renders rep view for sales_rep role', async () => {
    mockedGetCurrentUser.mockResolvedValue({ id: 'u1', role: 'sales_rep' } as never);
    const { container } = render(await KPIPage(pageProps as Parameters<typeof KPIPage>[0]));
    expect(container).toBeTruthy();
  });

  it('renders manager view for sales_manager role', async () => {
    mockedGetCurrentUser.mockResolvedValue({ id: 'm1', role: 'sales_manager' } as never);
    const { container } = render(await KPIPage(pageProps as Parameters<typeof KPIPage>[0]));
    expect(container).toBeTruthy();
  });

  it('redirects anonymous to login', async () => {
    mockedGetCurrentUser.mockResolvedValue(null as never);
    await expect(KPIPage(pageProps as Parameters<typeof KPIPage>[0])).rejects.toThrow();
  });
});
