'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { MessagesSquare, Plus } from 'lucide-react'
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
import { useConversations, useCreateConversation } from '@/src/application/hooks/use-conversations'

export function ConversationPage() {
	const router = useRouter()
	const { data: conversations, isPending, isError } = useConversations()
	const createConversation = useCreateConversation()
	const [dialogOpen, setDialogOpen] = useState(false)
	const [language, setLanguage] = useState('')

	const start = (e: React.FormEvent) => {
		e.preventDefault()
		if (!language.trim()) return
		createConversation.mutate(
			{ language: language.trim() },
			{
				onSuccess: (conversation) => {
					setDialogOpen(false)
					router.push(`/pollyglot/conversation/${conversation.id}`)
				},
			}
		)
	}

	return (
		<div className='mx-auto max-w-4xl'>
			<div className='mb-8 flex items-center justify-between'>
				<div>
					<h1 className='text-2xl font-bold tracking-tight'>Conversation</h1>
					<p className='text-muted-foreground'>
						Practice with a tutor that asks before it answers.
					</p>
				</div>
				<Button
					className='bg-emerald-600 text-white hover:bg-emerald-700'
					onClick={() => setDialogOpen(true)}
				>
					<Plus className='mr-2 h-4 w-4' />
					New conversation
				</Button>
			</div>

			{isPending && <p className='text-muted-foreground'>Loading conversations…</p>}

			{isError && (
				<p className='text-red-600 dark:text-red-400'>
					Could not load your conversations. Check that the API is running, then reload.
				</p>
			)}

			{conversations && conversations.length === 0 && (
				<div className='neu-inset rounded-xl p-12 text-center'>
					<MessagesSquare className='mx-auto mb-4 h-8 w-8 text-muted-foreground' />
					<p className='mb-1 font-medium'>No conversations yet</p>
					<p className='text-sm text-muted-foreground'>
						Start one and let the tutor draw the words out of you.
					</p>
				</div>
			)}

			{conversations && conversations.length > 0 && (
				<div className='grid gap-4 sm:grid-cols-2'>
					{conversations.map((conversation) => (
						<Link
							key={conversation.id}
							href={`/pollyglot/conversation/${conversation.id}`}
							className='group neu-card neu-interactive p-6'
						>
							<MessagesSquare className='mb-4 h-5 w-5 text-emerald-600 dark:text-emerald-400' />
							<h2 className='mb-1 text-sm font-semibold'>{conversation.title}</h2>
							<p className='text-sm text-muted-foreground'>{conversation.language}</p>
						</Link>
					))}
				</div>
			)}

			<Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
				<DialogContent className='sm:max-w-[425px]'>
					<DialogHeader>
						<DialogTitle>New conversation</DialogTitle>
						<DialogDescription>
							Tell the tutor which language you want to practice.
						</DialogDescription>
					</DialogHeader>
					<form onSubmit={start} className='space-y-4'>
						<div className='space-y-2'>
							<Label htmlFor='conversation-language'>Language</Label>
							<Input
								id='conversation-language'
								value={language}
								onChange={(e) => setLanguage(e.target.value)}
								placeholder='Japanese'
							/>
						</div>
						<DialogFooter>
							<Button type='submit' disabled={createConversation.isPending}>
								Start practicing
							</Button>
						</DialogFooter>
					</form>
				</DialogContent>
			</Dialog>
		</div>
	)
}
