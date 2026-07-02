'use client'

import Link from 'next/link'
import { GraduationCap } from 'lucide-react'
import { useDecks } from '@/src/application/hooks/use-decks'

export function StudyPage() {
	const { data: decks, isPending, isError } = useDecks()

	return (
		<div className='mx-auto max-w-4xl'>
			<div className='mb-8'>
				<h1 className='text-2xl font-bold tracking-tight'>Study</h1>
				<p className='text-muted-foreground'>Pick a deck to review what is due.</p>
			</div>

			{isPending && <p className='text-muted-foreground'>Loading decks…</p>}

			{isError && (
				<p className='text-red-600 dark:text-red-400'>
					Could not load your decks. Check that the API is running, then reload.
				</p>
			)}

			{decks && decks.length === 0 && (
				<div className='rounded-xl border border-dashed p-12 text-center'>
					<GraduationCap className='mx-auto mb-4 h-8 w-8 text-muted-foreground' />
					<p className='mb-1 font-medium'>No decks to study</p>
					<p className='mb-6 text-sm text-muted-foreground'>
						You need a deck with cards before you can review.
					</p>
					<Link
						href='/pollyglot/decks'
						className='text-sm font-medium text-emerald-600 hover:underline dark:text-emerald-400'
					>
						Create a deck
					</Link>
				</div>
			)}

			{decks && decks.length > 0 && (
				<div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
					{decks.map((deck) => (
						<Link
							key={deck.id}
							href={`/pollyglot/study/${deck.id}`}
							className='group rounded-xl border bg-card p-6 transition-colors hover:border-emerald-500/50 hover:bg-accent'
						>
							<GraduationCap className='mb-4 h-5 w-5 text-emerald-600 dark:text-emerald-400' />
							<h2 className='mb-1 text-sm font-semibold'>{deck.name}</h2>
							<p className='text-sm text-muted-foreground'>
								{deck.source_lang} → {deck.target_lang}
							</p>
							<p className='mt-2 text-xs text-muted-foreground'>
								{deck.card_count} {deck.card_count === 1 ? 'card' : 'cards'}
							</p>
						</Link>
					))}
				</div>
			)}
		</div>
	)
}
