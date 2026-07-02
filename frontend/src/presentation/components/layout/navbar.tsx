'use client'

import { Menu } from 'lucide-react'
import { Button } from '@/src/presentation/components/ui/button'
import { ThemeToggle } from '@/src/presentation/components/theme-toggle'

interface NavbarProps {
    onMenuToggle: () => void
}

export function Navbar({ onMenuToggle }: NavbarProps) {
    return (
        <nav className="sticky top-0 z-30 border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
            <div className="flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
                {/* Menu toggle for mobile/tablet */}
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={onMenuToggle}
                    className="lg:hidden"
                    aria-label="Toggle sidebar"
                >
                    <Menu className="h-5 w-5" />
                </Button>

                {/* Logo (visible on mobile when sidebar is hidden) */}
                <div className="flex items-center space-x-2 lg:hidden">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img src="/pollyglot.svg" alt="" className="h-8 w-8" />
                    <span className="text-xl font-bold">Pollyglot</span>
                </div>

                {/* Spacer for desktop */}
                <div className="hidden lg:block flex-1" />

                {/* Right side actions */}
                <div className="flex items-center gap-2">
                    <ThemeToggle />
                </div>
            </div>
        </nav>
    )
}
