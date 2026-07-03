'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, PartyPopper } from 'lucide-react'
import { Button } from '@/src/presentation/components/ui/button'
import { useDeck, useReviewCard, useStudyQueue } from '@/src/application/hooks/use-decks'
import type { Card, StudyRating } from '@/src/domain/services/decks.service'

const ratings: Array<{ label: string; value: StudyRating }> = [
	{ label: 'Forgot', value: 0 },
	{ label: 'Difficult', value: 1 },
	{ label: 'Okay', value: 2 },
	{ label: 'Almost', value: 3 },
	{ label: 'Got it!', value: 4 },
]

export function StudySessionPage({ deckId }: { deckId: string }) {
	const { data: deck } = useDeck(deckId)
	const { data: queue, isPending, isError } = useStudyQueue(deckId)
	const reviewCard = useReviewCard(deckId)

	// Snapshot the queue when it first arrives: a background refetch must
	// never reorder or shrink a session the learner is in the middle of.
	const [session, setSession] = useState<Card[] | null>(null)
	const [index, setIndex] = useState(0)
	const [flipped, setFlipped] = useState(false)
	const [flips, setFlips] = useState(0)
	const [reviewError, setReviewError] = useState(false)

	useEffect(() => {
		if (queue && session === null) {
			setSession(queue)
		}
	}, [queue, session])

	if (isError) {
		return (
			<div className='mx-auto max-w-xl text-center'>
				<p className='text-red-600 dark:text-red-400'>
					Could not load the study queue. Check that the API is running, then reload.
				</p>
			</div>
		)
	}

	if (isPending || session === null) {
		return <p className='text-center text-muted-foreground'>Loading your session…</p>
	}

	if (session.length === 0) {
		return (
			<div className='mx-auto max-w-xl neu-inset rounded-xl p-12 text-center'>
				<PartyPopper className='mx-auto mb-4 h-8 w-8 text-emerald-600 dark:text-emerald-400' />
				<p className='mb-1 font-medium'>All caught up!</p>
				<p className='mb-6 text-sm text-muted-foreground'>
					Nothing is due in {deck?.name ?? 'this deck'} right now.
				</p>
				<Link href='/pollyglot/decks' className='text-sm font-medium text-emerald-600 hover:underline dark:text-emerald-400'>
					Back to decks
				</Link>
			</div>
		)
	}

	if (index >= session.length) {
		return (
			<div className='mx-auto max-w-xl neu-inset rounded-xl p-12 text-center'>
				<PartyPopper className='mx-auto mb-4 h-8 w-8 text-emerald-600 dark:text-emerald-400' />
				<p className='mb-1 font-medium'>Session complete</p>
				<p className='mb-6 text-sm text-muted-foreground'>
					You reviewed {session.length} {session.length === 1 ? 'card' : 'cards'}.
				</p>
				<Link href='/pollyglot/decks' className='text-sm font-medium text-emerald-600 hover:underline dark:text-emerald-400'>
					Back to decks
				</Link>
			</div>
		)
	}

	const card = session[index]

	const toggleFlip = () => {
		setFlipped((f) => !f)
		setFlips((n) => n + 1)
	}

	const rate = (value: StudyRating) => {
		setReviewError(false)
		reviewCard.mutate(
			{ cardId: card.id, rating: value },
			{
				onSuccess: () => {
					setFlipped(false)
					setIndex((i) => i + 1)
				},
				onError: () => setReviewError(true),
			}
		)
	}

	return (
		<div className='mx-auto max-w-xl'>
			<div className='mb-6 flex items-center justify-between'>
				<Link
					href='/pollyglot/study'
					className='inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-foreground'
				>
					<ArrowLeft className='h-4 w-4' />
					{deck?.name ?? 'Study'}
				</Link>
				<p className='text-sm text-muted-foreground'>
					{index + 1} of {session.length}
				</p>
			</div>

			<div className='[perspective:1200px]'>
				<button
					type='button'
					aria-pressed={flipped}
					aria-label={flipped ? 'Hide answer' : 'Show answer'}
					onClick={toggleFlip}
					className='neu-card relative block h-64 w-full cursor-pointer rounded-2xl text-card-foreground transition-transform duration-500 [transform-style:preserve-3d] motion-reduce:transition-none focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500'
					style={{ transform: flipped ? 'rotateY(180deg)' : 'rotateY(0deg)' }}
				>
					<span className='absolute inset-0 flex flex-col items-center justify-center gap-3 p-6 [backface-visibility:hidden]'>
						<span className='text-xs font-medium uppercase tracking-widest text-emerald-600 dark:text-emerald-400'>
							{deck?.source_lang ?? ''}
						</span>
						<span className='text-3xl font-semibold'>{card.front}</span>
						<span className='text-xs text-muted-foreground'>Tap to reveal</span>
					</span>
					<span
						className='absolute inset-0 flex flex-col items-center justify-center gap-3 rounded-2xl bg-emerald-600/5 p-6 [backface-visibility:hidden]'
						style={{ transform: 'rotateY(180deg)' }}
					>
						<span className='text-xs font-medium uppercase tracking-widest text-muted-foreground'>
							{deck?.target_lang ?? ''}
						</span>
						<span className='text-3xl font-semibold'>{flipped ? card.back : ''}</span>
						<span className='text-xs text-muted-foreground'>How well did you know it?</span>
					</span>
				</button>
			</div>

			<div
				className={`mt-6 flex flex-wrap justify-center gap-2 transition-opacity duration-300 ${
					flipped ? 'opacity-100' : 'pointer-events-none opacity-30'
				}`}
			>
				{ratings.map(({ label, value }) => (
					<button
						key={label}
						type='button'
						aria-label={`Rate as ${label}`}
						tabIndex={flipped ? 0 : -1}
						disabled={reviewCard.isPending}
						onClick={() => rate(value)}
						className='neu-btn rounded-full px-4 py-1.5 text-sm font-medium text-muted-foreground hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 disabled:opacity-50'
					>
						{label}
					</button>
				))}
			</div>

			{reviewError && (
				<p className='mt-4 text-center text-sm text-red-600 dark:text-red-400'>
					Could not save that review. Check your connection and try again.
				</p>
			)}

			<p
				className='mt-8 text-center text-sm text-muted-foreground'
				aria-label={`Cards flipped: ${flips}`}
			>
				Cards Flipped: {flips}
			</p>
		</div>
	)
}
