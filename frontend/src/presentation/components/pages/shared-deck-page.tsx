'use client'

import { useMutation, useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { Layers, Copy } from 'lucide-react'
import { Button } from '@/src/presentation/components/ui/button'
import { decksService } from '@/src/domain/services/decks.service'

export function SharedDeckPage({ code }: { code: string }) {
	const router = useRouter()
	const { data: preview, isPending, isError } = useQuery({
		queryKey: ['shared', code],
		queryFn: () => decksService.getSharedDeck(code),
		retry: false,
	})
	const clone = useMutation({
		mutationFn: () => decksService.cloneSharedDeck(code),
		onSuccess: (deck) => router.push(`/pollyglot/decks/${deck.id}`),
	})

	if (isPending) {
		return <p className='text-center text-muted-foreground'>Loading shared deck…</p>
	}

	if (isError || !preview) {
		return (
			<div className='mx-auto max-w-xl'>
				<div className='neu-inset rounded-xl p-12 text-center'>
					<p className='mb-1 font-medium'>This deck is not found or no longer shared</p>
					<p className='text-sm text-muted-foreground'>
						Ask the owner for a fresh link.
					</p>
				</div>
			</div>
		)
	}

	return (
		<div className='mx-auto max-w-xl'>
			<div className='neu-card p-8'>
				<Layers className='mb-4 h-6 w-6 text-emerald-600 dark:text-emerald-400' />
				<h1 className='mb-1 text-2xl font-bold tracking-tight'>{preview.name}</h1>
				<p className='text-muted-foreground'>
					{preview.source_lang} → {preview.target_lang}
				</p>
				<p className='mt-1 text-sm text-muted-foreground'>
					{preview.card_count} {preview.card_count === 1 ? 'card' : 'cards'}
				</p>

				{preview.sample_cards.length > 0 && (
					<div className='mt-6 space-y-2'>
						<p className='text-xs font-medium uppercase tracking-widest text-muted-foreground'>
							Sample cards
						</p>
						{preview.sample_cards.map((card) => (
							<div key={card.front} className='neu-inset flex items-center justify-between rounded-lg px-4 py-2 text-sm'>
								<span className='font-medium'>{card.front}</span>
								<span className='text-muted-foreground'>{card.back}</span>
							</div>
						))}
					</div>
				)}

				<Button
					className='mt-8 w-full bg-emerald-600 text-white hover:bg-emerald-700'
					disabled={clone.isPending}
					onClick={() => clone.mutate()}
				>
					<Copy className='mr-2 h-4 w-4' />
					Add to my decks
				</Button>
				{clone.isError && (
					<p className='mt-2 text-sm text-red-600 dark:text-red-400'>
						Could not clone this deck. Try again.
					</p>
				)}
			</div>
		</div>
	)
}
