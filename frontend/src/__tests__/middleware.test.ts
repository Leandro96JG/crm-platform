/**
 * @jest-environment node
 */
import { NextRequest } from 'next/server';
import { middleware } from '../middleware';

function buildRequest(
  pathname: string,
  options: { jwt?: string; userAgent?: string } = {}
): NextRequest {
  const url = `http://localhost${pathname}`;
  const req = new NextRequest(url, {
    headers: {
      'user-agent': options.userAgent ?? 'Mozilla/5.0 (Windows NT 10.0)',
    },
  });
  if (options.jwt) {
    req.cookies.set('jwt', options.jwt);
  }
  return req;
}

describe('middleware', () => {
  describe('unauthenticated user', () => {
    it('should redirect to /login when accessing a protected route', () => {
      const req = buildRequest('/home');
      const res = middleware(req);
      expect(res?.headers.get('location')).toContain('/login');
    });

    it('should allow access to /login', () => {
      const req = buildRequest('/login');
      const res = middleware(req);
      expect(res?.headers.get('location')).toBeNull();
    });
  });

  describe('authenticated user', () => {
    it('should redirect from /login to /home', () => {
      const req = buildRequest('/login', { jwt: 'token' });
      const res = middleware(req);
      expect(res?.headers.get('location')).toContain('/home');
    });

    it('should allow access to /home', () => {
      const req = buildRequest('/home', { jwt: 'token' });
      const res = middleware(req);
      expect(res?.headers.get('location')).toBeNull();
    });
  });
});
