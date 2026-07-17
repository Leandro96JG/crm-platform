import type { Metadata } from 'next';
import ProgressBarProvider from './components/common/progress-bar/ProgressBarProvider';
import { Providers } from './providers';
import { roboto } from './ui/fonts';

import './ui/global.css';

export const metadata: Metadata = {
  title: {
    template: '%s | Viva',
    default: 'Viva',
  },
  description: 'Gestión de taller de stickers',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="es" className="dark" style={{ colorScheme: 'dark' }}>
      <head>
        <script
          dangerouslySetInnerHTML={{
            __html: `(function(){try{var t=localStorage.getItem('viva-theme');if(t!=='light'&&t!=='dark'){t='dark';}var r=document.documentElement;r.classList.toggle('dark',t==='dark');r.style.colorScheme=t;}catch(e){}})();`,
          }}
        />
      </head>
      <body className={roboto.className}>
        <Providers>
          <ProgressBarProvider>{children}</ProgressBarProvider>
        </Providers>
      </body>
    </html>
  );
}
