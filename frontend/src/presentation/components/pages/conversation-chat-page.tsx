'use client'

import { useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, Send, Volume2 } from 'lucide-react'
import { Button } from '@/src/presentation/components/ui/button'
import { Input } from '@/src/presentation/components/ui/input'
import { useConversations, useMessages, useSendMessage } from '@/src/application/hooks/use-conversations'
import { speakWithFallback } from '@/src/lib/speech'
import { cn } from '@/src/lib/utils'
import type { Message } from '@/src/domain/services/conversation.service'

function MessageBubble({ message, language }: { message: Message; language: string }) {
	const isTutor = message.role === 'tutor'
	return (
		<div
			aria-label={isTutor ? 'Tutor said' : 'You said'}
			className={cn('flex items-end gap-1', isTutor ? 'justify-start' : 'justify-end')}
		>
			<div
				className={cn(
					'max-w-[80%] rounded-2xl px-4 py-2.5 text-sm leading-relaxed',
					isTutor
						? 'neu-card-sm rounded-2xl rounded-bl-sm'
						: 'rounded-br-sm bg-emerald-600 text-white'
				)}
			>
				{message.content}
			</div>
			{isTutor && (
				<button
					type='button'
					aria-label='Play message'
					onClick={() => void speakWithFallback(message.content, language)}
					className='neu-btn rounded-full p-1.5 text-muted-foreground hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500'
				>
					<Volume2 className='h-3.5 w-3.5' aria-hidden />
				</button>
			)}
		</div>
	)
}

export function ConversationChatPage({ conversationId }: { conversationId: string }) {
	const { data: conversations } = useConversations()
	const { data: messages, isPending, isError } = useMessages(conversationId)
	const sendMessage = useSendMessage(conversationId)
	const [draft, setDraft] = useState('')
	const bottomRef = useRef<HTMLDivElement>(null)

	const conversation = conversations?.find((c) => c.id === conversationId)

	useEffect(() => {
		bottomRef.current?.scrollIntoView({ block: 'end' })
	}, [messages?.length])

	const submit = (e: React.FormEvent) => {
		e.preventDefault()
		const content = draft.trim()
		if (!content) return
		sendMessage.mutate(content, {
			onSuccess: () => setDraft(''),
		})
	}

	return (
		<div className='mx-auto flex h-[calc(100dvh-8rem)] max-w-2xl flex-col'>
			<div className='mb-4 flex items-center gap-3'>
				<Link
					href='/pollyglot/conversation'
					className='inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-foreground'
				>
					<ArrowLeft className='h-4 w-4' />
					Conversations
				</Link>
				{conversation && (
					<span className='text-sm font-semibold'>{conversation.title}</span>
				)}
			</div>

			<div className='neu-inset flex-1 space-y-3 overflow-y-auto rounded-xl p-4'>
				{isPending && <p className='text-sm text-muted-foreground'>Loading conversation…</p>}
				{isError && (
					<p className='text-sm text-red-600 dark:text-red-400'>
						Could not load this conversation. Check that the API is running, then reload.
					</p>
				)}
				{messages?.map((message) => (
					<MessageBubble key={message.id} message={message} language={conversation?.language ?? ''} />
				))}
				{sendMessage.isPending && (
					<p className='text-xs text-muted-foreground'>The tutor is thinking…</p>
				)}
				<div ref={bottomRef} />
			</div>

			{sendMessage.isError && (
				<p className='mt-2 text-sm text-red-600 dark:text-red-400'>
					Could not send that. Check your connection and try again.
				</p>
			)}

			<form onSubmit={submit} className='mt-4 flex gap-2'>
				<label htmlFor='chat-message' className='sr-only'>
					Message
				</label>
				<Input
					id='chat-message'
					value={draft}
					onChange={(e) => setDraft(e.target.value)}
					placeholder='Reply to your tutor…'
					autoComplete='off'
				/>
				<Button
					type='submit'
					disabled={sendMessage.isPending}
					className='bg-emerald-600 text-white hover:bg-emerald-700'
					aria-label='Send'
				>
					<Send className='h-4 w-4' />
				</Button>
			</form>
		</div>
	)
}
