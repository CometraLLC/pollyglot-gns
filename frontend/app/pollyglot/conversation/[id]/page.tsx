import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { ConversationChatPage } from '@/src/presentation/components/pages/conversation-chat-page'

export default async function Page({ params }: { params: Promise<{ id: string }> }) {
	const { id } = await params
	return (
		<ProtectedRoute>
			<MainLayout>
				<ConversationChatPage conversationId={id} />
			</MainLayout>
		</ProtectedRoute>
	)
}
