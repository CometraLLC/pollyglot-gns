'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, Pencil, Plus, Trash2 } from 'lucide-react'
import { Button } from '@/src/presentation/components/ui/button'
import { Input } from '@/src/presentation/components/ui/input'
import { Label } from '@/src/presentation/components/ui/label'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/src/presentation/components/ui/dialog'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from '@/src/presentation/components/ui/alert-dialog'
import {
	useCards,
	useCreateCard,
	useDeck,
	useDeleteCard,
	useUpdateCard,
} from '@/src/application/hooks/use-decks'
import type { Card, CardInput } from '@/src/domain/services/decks.service'

interface CardFormDialogProps {
	open: boolean
	onOpenChange: (open: boolean) => void
	title: string
	initial?: CardInput
	onSubmit: (input: CardInput) => void
	pending: boolean
}

function CardFormDialog({ open, onOpenChange, title, initial, onSubmit, pending }: CardFormDialogProps) {
	const [front, setFront] = useState(initial?.front ?? '')
	const [back, setBack] = useState(initial?.back ?? '')

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault()
		if (!front.trim() || !back.trim()) return
		onSubmit({ front: front.trim(), back: back.trim() })
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className='sm:max-w-[425px]'>
				<DialogHeader>
					<DialogTitle>{title}</DialogTitle>
					<DialogDescription>
						The front is the word you are learning; the back is what it means.
					</DialogDescription>
				</DialogHeader>
				<form onSubmit={handleSubmit} className='space-y-4'>
					<div className='space-y-2'>
						<Label htmlFor='card-front'>Front</Label>
						<Input
							id='card-front'
							value={front}
							onChange={(e) => setFront(e.target.value)}
							placeholder='こんにちは'
						/>
					</div>
					<div className='space-y-2'>
						<Label htmlFor='card-back'>Back</Label>
						<Input
							id='card-back'
							value={back}
							onChange={(e) => setBack(e.target.value)}
							placeholder='hello'
						/>
					</div>
					<DialogFooter>
						<Button type='submit' disabled={pending}>
							Save card
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	)
}

function CardRow({ card, deckId }: { card: Card; deckId: string }) {
	const [editOpen, setEditOpen] = useState(false)
	const [deleteOpen, setDeleteOpen] = useState(false)
	const updateCard = useUpdateCard(deckId)
	const deleteCard = useDeleteCard(deckId)

	const due = new Date(card.due_at) <= new Date()

	return (
		<div className='flex items-center justify-between gap-4 rounded-lg border bg-card px-4 py-3'>
			<div className='min-w-0 flex-1'>
				<p className='truncate font-medium'>{card.front}</p>
				<p className='truncate text-sm text-muted-foreground'>{card.back}</p>
			</div>
			<div className='flex items-center gap-2'>
				{due ? (
					<span className='rounded-full bg-emerald-500/15 px-2 py-0.5 text-xs font-medium text-emerald-600 dark:text-emerald-400'>
						Due
					</span>
				) : (
					<span className='text-xs text-muted-foreground'>
						Due {new Date(card.due_at).toLocaleDateString()}
					</span>
				)}
				<Button
					variant='ghost'
					size='icon'
					aria-label={`Edit card ${card.front}`}
					onClick={() => setEditOpen(true)}
				>
					<Pencil className='h-4 w-4' />
				</Button>
				<Button
					variant='ghost'
					size='icon'
					aria-label={`Delete card ${card.front}`}
					className='text-red-600 hover:text-red-600'
					onClick={() => setDeleteOpen(true)}
				>
					<Trash2 className='h-4 w-4' />
				</Button>
			</div>

			{editOpen && (
				<CardFormDialog
					open={editOpen}
					onOpenChange={setEditOpen}
					title='Edit card'
					initial={{ front: card.front, back: card.back }}
					pending={updateCard.isPending}
					onSubmit={(input) =>
						updateCard.mutate(
							{ cardId: card.id, data: input },
							{ onSuccess: () => setEditOpen(false) }
						)
					}
				/>
			)}

			<AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>Delete this card?</AlertDialogTitle>
						<AlertDialogDescription>
							“{card.front}” and its review history will be removed.
						</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel>Cancel</AlertDialogCancel>
						<AlertDialogAction
							className='bg-red-600 text-white hover:bg-red-700'
							onClick={() => deleteCard.mutate(card.id, { onSuccess: () => setDeleteOpen(false) })}
						>
							Delete
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>
		</div>
	)
}

export function DeckDetailPage({ deckId }: { deckId: string }) {
	const { data: deck } = useDeck(deckId)
	const { data: cards, isPending, isError } = useCards(deckId)
	const [addOpen, setAddOpen] = useState(false)
	const createCard = useCreateCard(deckId)

	return (
		<div className='mx-auto max-w-4xl'>
			<Link
				href='/pollyglot/decks'
				className='mb-4 inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-foreground'
			>
				<ArrowLeft className='h-4 w-4' />
				Back to decks
			</Link>

			<div className='mb-8 flex items-center justify-between'>
				<div>
					<h1 className='text-2xl font-bold tracking-tight'>{deck?.name ?? '…'}</h1>
					{deck && (
						<p className='text-muted-foreground'>
							{deck.source_lang} → {deck.target_lang}
						</p>
					)}
				</div>
				<Button
					className='bg-emerald-600 text-white hover:bg-emerald-700'
					onClick={() => setAddOpen(true)}
				>
					<Plus className='mr-2 h-4 w-4' />
					Add card
				</Button>
			</div>

			{isPending && <p className='text-muted-foreground'>Loading cards…</p>}

			{isError && (
				<p className='text-red-600 dark:text-red-400'>
					Could not load cards. Check that the API is running, then reload.
				</p>
			)}

			{cards && cards.length === 0 && (
				<div className='rounded-xl border border-dashed p-12 text-center'>
					<p className='mb-1 font-medium'>No cards yet</p>
					<p className='text-sm text-muted-foreground'>
						Add the first word you want to remember.
					</p>
				</div>
			)}

			{cards && cards.length > 0 && (
				<div className='space-y-2'>
					{cards.map((card) => (
						<CardRow key={card.id} card={card} deckId={deckId} />
					))}
				</div>
			)}

			{addOpen && (
				<CardFormDialog
					open={addOpen}
					onOpenChange={setAddOpen}
					title='Add card'
					pending={createCard.isPending}
					onSubmit={(input) => createCard.mutate(input, { onSuccess: () => setAddOpen(false) })}
				/>
			)}
		</div>
	)
}
