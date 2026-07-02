'use client'

import Link from 'next/link'
import { Layers, GraduationCap, Languages, MessagesSquare, TrendingUp } from 'lucide-react'
import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'

const tiles = [
	{
		title: 'Decks',
		description: 'Create decks and fill them with the words you are learning.',
		href: '/pollyglot/decks',
		icon: Layers,
		available: true,
	},
	{
		title: 'Study',
		description: 'Review the cards that are due today.',
		href: '/pollyglot/study',
		icon: GraduationCap,
		available: true,
	},
	{
		title: 'Translate',
		description: 'Translate a phrase and save it as a card.',
		href: '/pollyglot/translate',
		icon: Languages,
		available: false,
	},
	{
		title: 'Conversation',
		description: 'Practice with a tutor that asks before it answers.',
		href: '/pollyglot/conversation',
		icon: MessagesSquare,
		available: false,
	},
	{
		title: 'Progress',
		description: 'Streaks, reviews per day, and every word you have met.',
		href: '/pollyglot/stats',
		icon: TrendingUp,
		available: false,
	},
]

export default function PollyglotHome() {
	return (
		<ProtectedRoute>
			<MainLayout>
				<div className='mx-auto max-w-4xl'>
					<div className='mb-8'>
						<h1 className='text-2xl font-bold tracking-tight'>Pollyglot</h1>
						<p className='text-muted-foreground'>
							Flashcards, translation, and conversation practice in one place.
						</p>
					</div>
					<div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
						{tiles.map((tile) =>
							tile.available ? (
								<Link
									key={tile.title}
									href={tile.href}
									className='group rounded-xl border bg-card p-6 transition-colors hover:border-emerald-500/50 hover:bg-accent'>
									<tile.icon className='mb-4 h-5 w-5 text-emerald-600 dark:text-emerald-400' />
									<h2 className='mb-1 text-sm font-semibold'>{tile.title}</h2>
									<p className='text-sm text-muted-foreground'>{tile.description}</p>
								</Link>
							) : (
								<div
									key={tile.title}
									aria-disabled='true'
									className='rounded-xl border border-dashed bg-card/50 p-6 opacity-70'>
									<tile.icon className='mb-4 h-5 w-5 text-muted-foreground' />
									<div className='mb-1 flex items-center gap-2'>
										<h2 className='text-sm font-semibold'>{tile.title}</h2>
										<span className='rounded-full bg-muted px-2 py-0.5 text-[10px] font-medium uppercase tracking-wide text-muted-foreground'>
											Soon
										</span>
									</div>
									<p className='text-sm text-muted-foreground'>{tile.description}</p>
								</div>
							)
						)}
					</div>
				</div>
			</MainLayout>
		</ProtectedRoute>
	)
}
