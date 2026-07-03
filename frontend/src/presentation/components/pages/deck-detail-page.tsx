'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, Download, Pencil, Plus, Share2, Trash2, Upload } from 'lucide-react'
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
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { deckKeys } from '@/src/application/hooks/use-decks'
import { decksService } from '@/src/domain/services/decks.service'
import type { ImportResult } from '@/src/domain/services/decks.service'
import type { Card, CardInput, CardType } from '@/src/domain/services/decks.service'

interface CardFormDialogProps {
	open: boolean
	onOpenChange: (open: boolean) => void
	title: string
	initial?: CardInput
	onSubmit: (input: CardInput) => void
	pending: boolean
	// type/reverse options only apply when creating; edits change text only
	withTypeOptions?: boolean
}

function CardFormDialog({
	open,
	onOpenChange,
	title,
	initial,
	onSubmit,
	pending,
	withTypeOptions = false,
}: CardFormDialogProps) {
	const [front, setFront] = useState(initial?.front ?? '')
	const [back, setBack] = useState(initial?.back ?? '')
	const [cardType, setCardType] = useState<CardType>('basic')
	const [reverse, setReverse] = useState(false)

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault()
		if (!front.trim() || !back.trim()) return
		const input: CardInput = { front: front.trim(), back: back.trim() }
		if (withTypeOptions && cardType === 'cloze') input.card_type = 'cloze'
		if (withTypeOptions && cardType === 'basic' && reverse) input.reverse = true
		onSubmit(input)
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
					{withTypeOptions && (
						<div className='space-y-2'>
							<Label htmlFor='card-type'>Card type</Label>
							<select
								id='card-type'
								value={cardType}
								onChange={(e) => setCardType(e.target.value as CardType)}
								className='flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500'
							>
								<option value='basic'>Basic — front and back</option>
								<option value='cloze'>Cloze — fill in the blank</option>
							</select>
						</div>
					)}
					<div className='space-y-2'>
						<Label htmlFor='card-front'>Front</Label>
						<Input
							id='card-front'
							value={front}
							onChange={(e) => setFront(e.target.value)}
							placeholder={withTypeOptions && cardType === 'cloze' ? '水を{{c1::飲みます}}' : 'こんにちは'}
						/>
						{withTypeOptions && cardType === 'cloze' && (
							<p className='text-xs text-muted-foreground'>
								{'Wrap the hidden text like {{c1::word}} — it becomes the blank.'}
							</p>
						)}
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
					{withTypeOptions && cardType === 'basic' && (
						<div className='flex items-center gap-2'>
							<input
								id='card-reverse'
								type='checkbox'
								checked={reverse}
								onChange={(e) => setReverse(e.target.checked)}
								className='h-4 w-4 accent-emerald-600'
							/>
							<Label htmlFor='card-reverse'>Also create reversed card (back → front)</Label>
						</div>
					)}
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
		<div className='flex items-center justify-between gap-4 neu-card-sm px-4 py-3'>
			<div className='min-w-0 flex-1'>
				<p className='flex items-center gap-2 truncate font-medium'>
					{card.front}
					{card.card_type === 'cloze' && (
						<span className='rounded-full bg-muted px-2 py-0.5 text-[10px] font-medium uppercase tracking-wide text-muted-foreground'>
							cloze
						</span>
					)}
				</p>
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

function ImportDialog({ deckId, open, onOpenChange }: { deckId: string; open: boolean; onOpenChange: (open: boolean) => void }) {
	const queryClient = useQueryClient()
	const [file, setFile] = useState<File | null>(null)
	const [result, setResult] = useState<ImportResult | null>(null)
	const importDeck = useMutation({
		mutationFn: (upload: File) => decksService.importDeck(deckId, upload),
		onSuccess: (summary) => {
			setResult(summary)
			queryClient.invalidateQueries({ queryKey: deckKeys.cards(deckId) })
			queryClient.invalidateQueries({ queryKey: deckKeys.detail(deckId) })
			queryClient.invalidateQueries({ queryKey: deckKeys.lists() })
		},
	})

	const submit = (e: React.FormEvent) => {
		e.preventDefault()
		if (!file) return
		setResult(null)
		importDeck.mutate(file)
	}

	return (
		<Dialog open={open} onOpenChange={(next) => { onOpenChange(next); if (!next) { setFile(null); setResult(null) } }}>
			<DialogContent className='sm:max-w-[425px]'>
				<DialogHeader>
					<DialogTitle>Import cards</DialogTitle>
					<DialogDescription>
						CSV or TSV with front and back columns (a third card_type column is optional). Anki TSV exports work as-is.
					</DialogDescription>
				</DialogHeader>
				<form onSubmit={submit} className='space-y-4'>
					<div className='space-y-2'>
						<Label htmlFor='import-file'>CSV or TSV file</Label>
						<Input
							id='import-file'
							type='file'
							accept='.csv,.tsv,text/csv,text/tab-separated-values'
							onChange={(e) => setFile(e.target.files?.[0] ?? null)}
						/>
					</div>
					<DialogFooter>
						<Button type='submit' disabled={importDeck.isPending}>
							Import cards
						</Button>
					</DialogFooter>
				</form>
				{importDeck.isError && (
					<p className='text-sm text-red-600 dark:text-red-400'>
						Import failed. Check the file format and try again.
					</p>
				)}
				{result && (
					<div className='space-y-1 text-sm'>
						<p className='font-medium text-emerald-600 dark:text-emerald-400'>
							Imported {result.imported} {result.imported === 1 ? 'card' : 'cards'}.
						</p>
						{result.skipped.length > 0 && (
							<ul className='list-inside list-disc text-muted-foreground'>
								{result.skipped.map((row) => (
									<li key={row.line}>
										Line {row.line}: {row.error}
									</li>
								))}
							</ul>
						)}
					</div>
				)}
			</DialogContent>
		</Dialog>
	)
}

function ShareDialog({ deckId, shareCode, open, onOpenChange }: {
	deckId: string
	shareCode: string | null | undefined
	open: boolean
	onOpenChange: (open: boolean) => void
}) {
	const queryClient = useQueryClient()
	const invalidate = () => {
		queryClient.invalidateQueries({ queryKey: deckKeys.detail(deckId) })
		queryClient.invalidateQueries({ queryKey: deckKeys.lists() })
	}
	const share = useMutation({
		mutationFn: () => decksService.shareDeck(deckId),
		onSuccess: invalidate,
	})
	const unshare = useMutation({
		mutationFn: () => decksService.unshareDeck(deckId),
		onSuccess: invalidate,
	})

	const code = share.data?.share_code ?? shareCode
	const link = code ? `/pollyglot/shared/${code}` : null

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className='sm:max-w-[425px]'>
				<DialogHeader>
					<DialogTitle>Share deck</DialogTitle>
					<DialogDescription>
						Anyone signed in with the link can preview this deck and add their own copy.
					</DialogDescription>
				</DialogHeader>
				{link ? (
					<div className='space-y-4'>
						<div className='neu-inset rounded-lg px-4 py-3 text-sm break-all'>{link}</div>
						<div className='flex gap-2'>
							<Button
								variant='outline'
								onClick={() => void navigator.clipboard?.writeText(`${window.location.origin}${link}`)}
							>
								Copy link
							</Button>
							<Button
								variant='outline'
								className='text-red-600 hover:text-red-600'
								disabled={unshare.isPending}
								onClick={() => unshare.mutate()}
							>
								Disable sharing
							</Button>
						</div>
					</div>
				) : (
					<Button
						className='bg-emerald-600 text-white hover:bg-emerald-700'
						disabled={share.isPending}
						onClick={() => share.mutate()}
					>
						Enable sharing
					</Button>
				)}
				{(share.isError || unshare.isError) && (
					<p className='text-sm text-red-600 dark:text-red-400'>
						Could not update sharing. Try again.
					</p>
				)}
			</DialogContent>
		</Dialog>
	)
}

function downloadBlob(blob: Blob, filename: string) {
	const url = URL.createObjectURL(blob)
	const link = document.createElement('a')
	link.href = url
	link.download = filename
	link.click()
	URL.revokeObjectURL(url)
}

export function DeckDetailPage({ deckId }: { deckId: string }) {
	const { data: deck } = useDeck(deckId)
	const { data: cards, isPending, isError } = useCards(deckId)
	const [addOpen, setAddOpen] = useState(false)
	const [importOpen, setImportOpen] = useState(false)
	const [shareOpen, setShareOpen] = useState(false)
	const createCard = useCreateCard(deckId)
	const exportDeck = useMutation({
		mutationFn: () => decksService.exportDeck(deckId, 'csv'),
		onSuccess: (blob) => downloadBlob(blob, `${deck?.name ?? 'deck'}.csv`),
	})

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
				<div className='flex flex-wrap gap-2'>
					<Button variant='outline' onClick={() => setShareOpen(true)}>
						<Share2 className='mr-2 h-4 w-4' />
						Share
					</Button>
					<Button variant='outline' onClick={() => exportDeck.mutate()} disabled={exportDeck.isPending}>
						<Download className='mr-2 h-4 w-4' />
						Export CSV
					</Button>
					<Button variant='outline' onClick={() => setImportOpen(true)}>
						<Upload className='mr-2 h-4 w-4' />
						Import
					</Button>
					<Button
						className='bg-emerald-600 text-white hover:bg-emerald-700'
						onClick={() => setAddOpen(true)}
					>
						<Plus className='mr-2 h-4 w-4' />
						Add card
					</Button>
				</div>
			</div>

			{isPending && <p className='text-muted-foreground'>Loading cards…</p>}

			{isError && (
				<p className='text-red-600 dark:text-red-400'>
					Could not load cards. Check that the API is running, then reload.
				</p>
			)}

			{cards && cards.length === 0 && (
				<div className='neu-inset rounded-xl p-12 text-center'>
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

			<ImportDialog deckId={deckId} open={importOpen} onOpenChange={setImportOpen} />
			<ShareDialog deckId={deckId} shareCode={deck?.share_code} open={shareOpen} onOpenChange={setShareOpen} />

			{addOpen && (
				<CardFormDialog
					open={addOpen}
					onOpenChange={setAddOpen}
					title='Add card'
					withTypeOptions
					pending={createCard.isPending}
					onSubmit={(input) => createCard.mutate(input, { onSuccess: () => setAddOpen(false) })}
				/>
			)}
		</div>
	)
}
