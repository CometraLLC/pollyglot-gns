'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Layers, Pencil, Plus, Trash2 } from 'lucide-react'
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
	useCreateDeck,
	useDecks,
	useDeleteDeck,
	useUpdateDeck,
} from '@/src/application/hooks/use-decks'
import type { Deck, DeckInput } from '@/src/domain/services/decks.service'

interface DeckFormDialogProps {
	open: boolean
	onOpenChange: (open: boolean) => void
	title: string
	submitLabel: string
	initial?: DeckInput
	onSubmit: (input: DeckInput) => void
	pending: boolean
}

function DeckFormDialog({
	open,
	onOpenChange,
	title,
	submitLabel,
	initial,
	onSubmit,
	pending,
}: DeckFormDialogProps) {
	const [name, setName] = useState(initial?.name ?? '')
	const [sourceLang, setSourceLang] = useState(initial?.source_lang ?? '')
	const [targetLang, setTargetLang] = useState(initial?.target_lang ?? '')

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault()
		if (!name.trim() || !sourceLang.trim() || !targetLang.trim()) return
		onSubmit({ name: name.trim(), source_lang: sourceLang.trim(), target_lang: targetLang.trim() })
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className='sm:max-w-[425px]'>
				<DialogHeader>
					<DialogTitle>{title}</DialogTitle>
					<DialogDescription>
						A deck holds the cards for one language you are learning.
					</DialogDescription>
				</DialogHeader>
				<form onSubmit={handleSubmit} className='space-y-4'>
					<div className='space-y-2'>
						<Label htmlFor='deck-name'>Name</Label>
						<Input
							id='deck-name'
							value={name}
							onChange={(e) => setName(e.target.value)}
							placeholder='Japanese Basics'
						/>
					</div>
					<div className='space-y-2'>
						<Label htmlFor='deck-source'>Learning language</Label>
						<Input
							id='deck-source'
							value={sourceLang}
							onChange={(e) => setSourceLang(e.target.value)}
							placeholder='Japanese'
						/>
					</div>
					<div className='space-y-2'>
						<Label htmlFor='deck-target'>Your language</Label>
						<Input
							id='deck-target'
							value={targetLang}
							onChange={(e) => setTargetLang(e.target.value)}
							placeholder='English'
						/>
					</div>
					<DialogFooter>
						<Button type='submit' disabled={pending}>
							{submitLabel}
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	)
}

function DeckCard({ deck }: { deck: Deck }) {
	const [editOpen, setEditOpen] = useState(false)
	const [deleteOpen, setDeleteOpen] = useState(false)
	const updateDeck = useUpdateDeck()
	const deleteDeck = useDeleteDeck()

	return (
		<div className='group neu-card neu-interactive p-6'>
			<div className='mb-4 flex items-start justify-between'>
				<Layers className='h-5 w-5 text-emerald-600 dark:text-emerald-400' />
				<div className='flex gap-1 opacity-100 lg:opacity-0 lg:transition-opacity lg:group-hover:opacity-100 lg:group-focus-within:opacity-100'>
					<Button
						variant='ghost'
						size='icon'
						aria-label={`Edit ${deck.name}`}
						onClick={() => setEditOpen(true)}
					>
						<Pencil className='h-4 w-4' />
					</Button>
					<Button
						variant='ghost'
						size='icon'
						aria-label={`Delete ${deck.name}`}
						className='text-red-600 hover:text-red-600'
						onClick={() => setDeleteOpen(true)}
					>
						<Trash2 className='h-4 w-4' />
					</Button>
				</div>
			</div>
			<Link href={`/pollyglot/decks/${deck.id}`} className='block'>
				<h2 className='mb-1 text-sm font-semibold hover:underline'>{deck.name}</h2>
			</Link>
			<p className='text-sm text-muted-foreground'>
				{deck.source_lang} → {deck.target_lang}
			</p>
			<p className='mt-2 flex items-center gap-2 text-xs text-muted-foreground'>
				{deck.card_count} {deck.card_count === 1 ? 'card' : 'cards'}
				{deck.due_count > 0 && (
					<span className='rounded-full bg-emerald-500/15 px-2 py-0.5 text-xs font-medium text-emerald-600 dark:text-emerald-400'>
						{deck.due_count} due
					</span>
				)}
			</p>

			{editOpen && (
				<DeckFormDialog
					open={editOpen}
					onOpenChange={setEditOpen}
					title='Edit deck'
					submitLabel='Save'
					initial={{ name: deck.name, source_lang: deck.source_lang, target_lang: deck.target_lang }}
					pending={updateDeck.isPending}
					onSubmit={(input) =>
						updateDeck.mutate(
							{ id: deck.id, data: input },
							{ onSuccess: () => setEditOpen(false) }
						)
					}
				/>
			)}

			<AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>Delete “{deck.name}”?</AlertDialogTitle>
						<AlertDialogDescription>
							This removes the deck and its {deck.card_count}{' '}
							{deck.card_count === 1 ? 'card' : 'cards'} from your study rotation.
						</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel>Cancel</AlertDialogCancel>
						<AlertDialogAction
							className='bg-red-600 text-white hover:bg-red-700'
							onClick={() => deleteDeck.mutate(deck.id, { onSuccess: () => setDeleteOpen(false) })}
						>
							Delete
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>
		</div>
	)
}

export function DecksPage() {
	const { data: decks, isPending, isError } = useDecks()
	const [createOpen, setCreateOpen] = useState(false)
	const createDeck = useCreateDeck()

	return (
		<div className='mx-auto max-w-4xl'>
			<div className='mb-8 flex items-center justify-between'>
				<div>
					<h1 className='text-2xl font-bold tracking-tight'>Decks</h1>
					<p className='text-muted-foreground'>Each deck is a language you are learning.</p>
				</div>
				<Button
					className='bg-emerald-600 text-white hover:bg-emerald-700'
					onClick={() => setCreateOpen(true)}
				>
					<Plus className='mr-2 h-4 w-4' />
					New deck
				</Button>
			</div>

			{isPending && <p className='text-muted-foreground'>Loading decks…</p>}

			{isError && (
				<p className='text-red-600 dark:text-red-400'>
					Could not load your decks. Check that the API is running, then reload.
				</p>
			)}

			{decks && decks.length === 0 && (
				<div className='neu-inset rounded-xl p-12 text-center'>
					<Layers className='mx-auto mb-4 h-8 w-8 text-muted-foreground' />
					<p className='mb-1 font-medium'>No decks yet</p>
					<p className='text-sm text-muted-foreground'>
						Create your first deck to start collecting words.
					</p>
				</div>
			)}

			{decks && decks.length > 0 && (
				<div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
					{decks.map((deck) => (
						<DeckCard key={deck.id} deck={deck} />
					))}
				</div>
			)}

			{createOpen && (
				<DeckFormDialog
					open={createOpen}
					onOpenChange={setCreateOpen}
					title='New deck'
					submitLabel='Create deck'
					pending={createDeck.isPending}
					onSubmit={(input) =>
						createDeck.mutate(input, { onSuccess: () => setCreateOpen(false) })
					}
				/>
			)}
		</div>
	)
}
