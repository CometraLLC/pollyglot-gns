'use client'

import { QueryClientProvider } from '@tanstack/react-query'
import { queryClient } from '@/src/lib/query-client'
import { ThemeProvider } from '@/src/presentation/components/theme-provider'
import { NextIntlClientProvider } from 'next-intl'
import { useLocaleStore } from '@/src/infrastructure/stores/locale-store'
import { ReactNode, useEffect, useState } from 'react'

export function Providers({ children }: { children: ReactNode }) {
    const { locale, hasHydrated } = useLocaleStore();
    const [messages, setMessages] = useState<any>(null);

    useEffect(() => {
        if (!hasHydrated) {
            useLocaleStore.persist.rehydrate();
        }
    }, [hasHydrated]);

    useEffect(() => {
        // Load messages for current locale
        import(`../locales/${hasHydrated ? locale : 'en'}.json`)
            .then((msgs) => setMessages(msgs.default))
            .catch(() => setMessages({}));
    }, [locale, hasHydrated]);

    // Detect theme for loading screen (before ThemeProvider renders)
    const [isDark, setIsDark] = useState(true); // Default to dark

    useEffect(() => {
        // Check localStorage for saved theme
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme) {
            setIsDark(savedTheme === 'dark');
        } else {
            // Fallback to system preference
            setIsDark(window.matchMedia('(prefers-color-scheme: dark)').matches);
        }
    }, []);

    if (!messages) {
        return (
            <div className={`fixed inset-0 z-50 flex items-center justify-center ${isDark ? 'bg-[#0a0a0a]' : 'bg-white'}`}>

                {/* Loading Content */}
                <div className="relative flex flex-col items-center gap-6">
                    {/* Animated Logo */}
                    <div className="relative">
                        <div className="absolute inset-0 animate-ping opacity-10">
                            {/* eslint-disable-next-line @next/next/no-img-element */}
                            <img src="/pollyglot.svg" alt="" className="h-16 w-16" />
                        </div>
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img src="/pollyglot.svg" alt="Pollyglot" className="h-16 w-16 animate-pulse" />
                    </div>

                    {/* Brand Name */}
                    <div className="flex flex-col items-center gap-2">
                        <span className={`text-2xl font-bold tracking-tight ${isDark ? 'text-white' : 'text-gray-900'}`}>Pollyglot</span>
                        <span className={`text-sm ${isDark ? 'text-gray-400' : 'text-gray-500'}`}>Loading your experience...</span>
                    </div>

                    {/* Loading Bar */}
                    <div className={`w-48 h-1 ${isDark ? 'bg-gray-800' : 'bg-gray-200'} rounded-full overflow-hidden`}>
                        <div className={`h-full ${isDark ? 'bg-white' : 'bg-gray-900'} rounded-full animate-loading-bar`} />
                    </div>
                </div>
            </div>
        );
    }

    return (
        <NextIntlClientProvider locale={hasHydrated ? locale : 'en'} messages={messages}>
            <ThemeProvider
                attribute="class"
                defaultTheme="dark"
                enableSystem={false}
                disableTransitionOnChange
            >
                <QueryClientProvider client={queryClient}>
                    {children}
                </QueryClientProvider>
            </ThemeProvider>
        </NextIntlClientProvider>
    )
}
