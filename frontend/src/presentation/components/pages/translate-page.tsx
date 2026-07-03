'use client'

import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { AxiosError } from 'axios'
import { ArrowLeftRight, Languages } from 'lucide-react'
import { Button } from '@/src/presentation/components/ui/button'
import { Input } from '@/src/presentation/components/ui/input'
import { Label } from '@/src/presentation/components/ui/label'
import { Textarea } from '@/src/presentation/components/ui/textarea'
import { useCreateCard, useDecks } from '@/src/application/hooks/use-decks'
import { translateService } from '@/src/domain/services/translate.service'
import type { Translation } from '@/src/domain/services/translate.service'

function errorMessage(error: unknown): string {
	if (error instanceof AxiosError) {
		const apiError = (error.response?.data as { error?: string } | undefined)?.error
		if (apiError) return apiError
	}
	return 'Translation failed. Check that the API is running, then try again.'
}

function SaveToDeck({ translation }: { translation: Translation }) {
	const { data: decks } = useDecks()
	const [deckId, setDeckId] = useState('')
	const [savedTo, setSavedTo] = useState<string | null>(null)
	const createCard = useCreateCard(deckId)

	if (!decks || decks.length === 0) return null

	const save = () => {
		if (!deckId) return
		const deck = decks.find((d) => d.id === deckId)
		createCard.mutate(
			{ front: translation.text, back: translation.translation },
			{ onSuccess: () => setSavedTo(deck?.name ?? 'deck') }
		)
	}

	return (
		<div className='mt-4 flex flex-wrap items-end gap-3 border-t pt-4'>
			<div className='space-y-2'>
				<Label htmlFor='save-deck'>Save to deck</Label>
				<select
					id='save-deck'
					value={deckId}
					onChange={(e) => {
						setDeckId(e.target.value)
						setSavedTo(null)
					}}
					className='flex h-9 w-56 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500'
				>
					<option value=''>Choose a deck…</option>
					{decks.map((deck) => (
						<option key={deck.id} value={deck.id}>
							{deck.name}
						</option>
					))}
				</select>
			</div>
			<Button
				variant='outline'
				disabled={!deckId || createCard.isPending}
				onClick={save}
			>
				Save as card
			</Button>
			{savedTo && (
				<p className='text-sm text-emerald-600 dark:text-emerald-400'>Saved to {savedTo}.</p>
			)}
			{createCard.isError && (
				<p className='text-sm text-red-600 dark:text-red-400'>
					Could not save the card. Try again.
				</p>
			)}
		</div>
	)
}

export function TranslatePage() {
	const [text, setText] = useState('')
	const [from, setFrom] = useState('Japanese')
	const [to, setTo] = useState('English')

	const translate = useMutation({
		mutationFn: (input: Parameters<typeof translateService.translate>[0]) =>
			translateService.translate(input),
	})

	const submit = (e: React.FormEvent) => {
		e.preventDefault()
		if (!text.trim()) return
		translate.mutate({ text: text.trim(), from: from.trim(), to: to.trim() })
	}

	const swap = () => {
		setFrom(to)
		setTo(from)
	}

	return (
		<div className='mx-auto max-w-2xl'>
			<div className='mb-8'>
				<h1 className='text-2xl font-bold tracking-tight'>Translate</h1>
				<p className='text-muted-foreground'>
					Translate a phrase, then keep it as a card so it becomes vocabulary.
				</p>
			</div>

			<form onSubmit={submit} className='space-y-4'>
				<div className='flex items-end gap-2'>
					<div className='flex-1 space-y-2'>
						<Label htmlFor='translate-from'>From</Label>
						<Input
							id='translate-from'
							value={from}
							onChange={(e) => setFrom(e.target.value)}
							placeholder='Japanese'
						/>
					</div>
					<Button
						type='button'
						variant='ghost'
						size='icon'
						aria-label='Swap languages'
						onClick={swap}
					>
						<ArrowLeftRight className='h-4 w-4' />
					</Button>
					<div className='flex-1 space-y-2'>
						<Label htmlFor='translate-to'>To</Label>
						<Input
							id='translate-to'
							value={to}
							onChange={(e) => setTo(e.target.value)}
							placeholder='English'
						/>
					</div>
				</div>

				<div className='space-y-2'>
					<Label htmlFor='translate-text'>Text to translate</Label>
					<Textarea
						id='translate-text'
						value={text}
						onChange={(e) => setText(e.target.value)}
						placeholder='こんにちは'
						rows={3}
					/>
				</div>

				<Button
					type='submit'
					disabled={translate.isPending}
					className='bg-emerald-600 text-white hover:bg-emerald-700'
				>
					<Languages className='mr-2 h-4 w-4' />
					Translate
				</Button>
			</form>

			{translate.isError && (
				<p className='mt-6 text-sm text-red-600 dark:text-red-400'>
					{errorMessage(translate.error)}
				</p>
			)}

			{translate.data && (
				<div className='mt-6 neu-card p-6'>
					<p className='mb-1 text-xs font-medium uppercase tracking-widest text-muted-foreground'>
						{translate.data.to}
					</p>
					<p className='text-2xl font-semibold'>{translate.data.translation}</p>
					<SaveToDeck translation={translate.data} />
				</div>
			)}
		</div>
	)
}
