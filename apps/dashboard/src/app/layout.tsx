import { Outfit } from 'next/font/google';
import './globals.css';
import "flatpickr/dist/flatpickr.css";
import { SidebarProvider } from '@/context/SidebarContext';
import { ThemeProvider } from '@/context/ThemeContext';
import { ProcessProvider } from '@/context/ProcessContext';

import { getAllProcess } from '@/lib/documentProcessing';

const outfit = Outfit({
  subsets: ["latin"],
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {

  const processPromise = getAllProcess();

  return (
    <html lang="en">
      <body className={`${outfit.className} dark:bg-gray-900`}>
        <ThemeProvider>
          <SidebarProvider>
            <ProcessProvider processPromise={processPromise}>
              {children}
            </ProcessProvider>
          </SidebarProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
