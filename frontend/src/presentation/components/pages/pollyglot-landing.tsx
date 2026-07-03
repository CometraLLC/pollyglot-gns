'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowRight, Layers, Languages, MessagesSquare, TrendingUp } from 'lucide-react'
import { ThemeToggle } from '@/src/presentation/components/theme-toggle'
import { Button } from '@/src/presentation/components/ui/button'

// ─── Demo deck for the hero flashcard ────────────────────
const demoCards = [
    { language: 'Japanese', front: 'こんにちは', back: 'hello' },
    { language: 'Spanish', front: 'gato', back: 'cat' },
    { language: 'French', front: 'merci', back: 'thank you' },
]

const ratings = ['Forgot', 'Difficult', 'Okay', 'Almost', 'Got it!'] as const

function HeroFlashcard() {
    const [index, setIndex] = useState(0)
    const [flipped, setFlipped] = useState(false)
    const card = demoCards[index % demoCards.length]

    const rate = () => {
        setFlipped(false)
        // Wait for the flip-back before swapping content so the answer
        // of the next card is never visible mid-turn.
        setTimeout(() => setIndex((i) => i + 1), 250)
    }

    return (
        <div className="w-full max-w-sm">
            <p className="mb-3 text-center text-xs font-medium uppercase tracking-widest text-muted-foreground">
                Try a card
            </p>
            <div className="[perspective:1200px]">
                <button
                    type="button"
                    aria-pressed={flipped}
                    aria-label={flipped ? 'Hide answer' : 'Show answer'}
                    onClick={() => setFlipped((f) => !f)}
                    className="neu-card relative block h-60 w-full cursor-pointer rounded-2xl text-card-foreground transition-transform duration-500 [transform-style:preserve-3d] motion-reduce:transition-none focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
                    style={{ transform: flipped ? 'rotateY(180deg)' : 'rotateY(0deg)' }}
                >
                    {/* Front */}
                    <span className="absolute inset-0 flex flex-col items-center justify-center gap-3 p-6 [backface-visibility:hidden]">
                        <span className="text-xs font-medium uppercase tracking-widest text-emerald-600 dark:text-emerald-400">
                            {card.language}
                        </span>
                        <span className="text-3xl font-semibold">{card.front}</span>
                        <span className="text-xs text-muted-foreground">Tap to reveal</span>
                    </span>
                    {/* Back */}
                    <span
                        className="absolute inset-0 flex flex-col items-center justify-center gap-3 rounded-2xl bg-emerald-600/5 p-6 [backface-visibility:hidden]"
                        style={{ transform: 'rotateY(180deg)' }}
                    >
                        <span className="text-xs font-medium uppercase tracking-widest text-muted-foreground">
                            English
                        </span>
                        <span className="text-3xl font-semibold">{card.back}</span>
                        <span className="text-xs text-muted-foreground">How well did you know it?</span>
                    </span>
                </button>
            </div>
            {/* Rating row — active only once the answer is showing */}
            <div
                className={`mt-4 flex flex-wrap justify-center gap-2 transition-opacity duration-300 ${
                    flipped ? 'opacity-100' : 'pointer-events-none opacity-30'
                }`}
            >
                {ratings.map((label) => (
                    <button
                        key={label}
                        type="button"
                        aria-label={`Rate as ${label}`}
                        tabIndex={flipped ? 0 : -1}
                        onClick={rate}
                        className="neu-btn rounded-full px-3 py-1 text-xs font-medium text-muted-foreground hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
                    >
                        {label}
                    </button>
                ))}
            </div>
        </div>
    )
}

// ─── Landing ─────────────────────────────────────────────
const features = [
    {
        icon: Layers,
        title: 'Spaced repetition',
        body: 'Rate each card from Forgot to Got it! — scheduling brings it back just before you would lose it.',
    },
    {
        icon: Languages,
        title: 'Translate and keep it',
        body: 'Translate any phrase you meet, then save it straight into a deck so it becomes a card, not a memory.',
    },
    {
        icon: MessagesSquare,
        title: 'Conversation practice',
        body: 'A tutor that asks before it answers, so the words you learn come out of your own mouth first.',
    },
    {
        icon: TrendingUp,
        title: 'Progress you can see',
        body: 'Streaks, reviews per day, and every word you have met — counted, not guessed.',
    },
]

const greetings = 'hello · hola · こんにちは · bonjour · 안녕하세요 · ciao · مرحبا · olá · hallo · привет'

export function PollyglotLanding() {
    return (
        <div className="min-h-screen bg-background text-foreground">
            {/* Nav */}
            <nav className="fixed left-0 right-0 top-0 z-50 border-b border-border/50 bg-background/80 backdrop-blur-xl">
                <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
                    <Link href="/" className="flex items-center gap-2.5">
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img src="/pollyglot.svg" alt="" className="h-8 w-8" />
                        <span className="text-lg font-semibold tracking-tight">Pollyglot</span>
                    </Link>
                    <div className="flex items-center gap-2">
                        <ThemeToggle />
                        <Link href="/auth/login">
                            <Button variant="ghost" size="sm">
                                Sign in
                            </Button>
                        </Link>
                        <Link href="/auth/register">
                            <Button size="sm" className="bg-emerald-600 text-white hover:bg-emerald-700">
                                Start learning
                            </Button>
                        </Link>
                    </div>
                </div>
            </nav>

            {/* Hero */}
            <section className="relative overflow-hidden pt-32 pb-16 md:pt-40 md:pb-24">
                <div className="pointer-events-none absolute inset-0 overflow-hidden">
                    <div className="absolute left-[70%] top-[20%] h-[420px] w-[420px] -translate-x-1/2 rounded-full bg-emerald-500/10 blur-[120px]" />
                </div>
                <div className="relative mx-auto grid max-w-6xl items-center gap-12 px-6 md:grid-cols-2">
                    <div>
                        <p className="mb-6 text-xs font-medium uppercase tracking-widest text-emerald-600 dark:text-emerald-400">
                            Flashcards · Translation · Conversation
                        </p>
                        <h1 className="mb-6 text-4xl font-bold leading-[1.1] tracking-tight md:text-5xl">
                            Learn languages the way memory works.
                        </h1>
                        <p className="mb-8 max-w-md text-muted-foreground">
                            Pollyglot schedules each card for the moment you are about to forget it,
                            translates anything you meet in the wild, and gives you a partner to
                            practice with.
                        </p>
                        <div className="flex flex-wrap gap-3">
                            <Link href="/auth/register">
                                <Button className="bg-emerald-600 text-white hover:bg-emerald-700">
                                    Start learning
                                    <ArrowRight className="ml-2 h-4 w-4" />
                                </Button>
                            </Link>
                            <Link href="/auth/login">
                                <Button variant="outline">Sign in</Button>
                            </Link>
                        </div>
                    </div>
                    <div className="flex justify-center">
                        <HeroFlashcard />
                    </div>
                </div>
                {/* Greeting strip */}
                <p className="mx-auto mt-16 max-w-6xl px-6 text-center text-sm text-muted-foreground/60" lang="und">
                    {greetings}
                </p>
            </section>

            {/* Features */}
            <section className="border-t bg-muted/30 py-20">
                <div className="mx-auto max-w-6xl px-6">
                    <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
                        {features.map((feature) => (
                            <div key={feature.title} className="neu-card p-6">
                                <feature.icon className="mb-4 h-5 w-5 text-emerald-600 dark:text-emerald-400" />
                                <h2 className="mb-2 text-sm font-semibold">{feature.title}</h2>
                                <p className="text-sm leading-relaxed text-muted-foreground">{feature.body}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Footer */}
            <footer className="border-t py-8">
                <div className="mx-auto flex max-w-6xl items-center justify-between px-6">
                    <div className="flex items-center gap-2">
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img src="/pollyglot.svg" alt="" className="h-5 w-5" />
                        <span className="text-sm text-muted-foreground">
                            Pollyglot © {new Date().getFullYear()} — built by Cometra
                        </span>
                    </div>
                    <Link
                        href="/docs"
                        className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                        Documentation
                    </Link>
                </div>
            </footer>
        </div>
    )
}
